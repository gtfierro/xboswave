from pyxbos.driver import *
import logging
import os, re, io
import requests
import time
import urllib
import yaml

from bs4 import BeautifulSoup as bs
from datetime import datetime, timedelta
import sensordb
import pandas as pd

try:
    import ordereddict
except ImportError:
    import collections as ordereddict

class ObviusDriver(Driver):
    def setup(self, cfg):
        self.bmoroot = cfg['obvius']['bmoroot']
        self.statuspage = cfg['obvius']['statuspage']
        self.auth = (cfg['obvius']['username'], cfg['obvius']['password'])
        self.devices = {}
        self.conf = {}
        self.rate = cfg['rate']
        self.discoverMeters()

    def discoverMeters(self):
        # find all the AcquiSuite boxes
        response = requests.get(self.bmoroot + self.statuspage, auth=self.auth)
        soup = bs(response.content, features="html.parser")

        for tr in soup.findAll('tr'):
            tds = tr('td')
            if len(tds) != 6:
                continue

            name = tds[0].a.string
            self.devices[name] = {
                'ip' : self.remove_nbsp(tds[3].string),
                'href' : tds[0].a['href'],
                }

        # look at all the meters hanging off each of them
        for location in self.devices:
             response = requests.get(self.bmoroot + self.devices[location]['href'], auth=self.auth)
             soup = bs(response.content, features="html.parser")
             subdevices = []
             for tr in soup.findAll('tr'):
                 tds = tr('td')
                 if len(tds) != 5 or tds[3].a != None:
                     continue
                 subdevices.append({
                     'address' : re.sub("<.*?>", "", str(tds[0])),
                     'status': self.remove_nbsp(tds[1].string),
                     'name' : self.remove_nbsp(tds[2].string),
                     'type' : self.remove_nbsp(tds[3].string),
                     'firmware': self.remove_nbsp(tds[4].string)
                 })
             self.devices[location]['subdevices'] = subdevices

        for location, devs in self.devices.items():
            params = urllib.parse.parse_qs(urllib.parse.urlsplit(devs['href']).query)
            if "AS" not in params or "DB" not in params:
                continue
            if location in self.auth:
                continue
            thisconf = {}
            for dev in devs['subdevices']:
                if sensordb.get_map(dev['type'], location) != None:
                    dlurl = self.bmoroot + 'mbdev_export.php/' + params['AS'][0] + '_' +  \
                        dev['address'] + '.csv' + "?DB=" + params['DB'][0] + '&AS=' + \
                        params['AS'][0] + '&MB=' + dev['address'] + '&DOWNLOAD=YES' + \
                        "&COLNAMES=ON&EXPORTTIMEZONE=UTC&DELIMITER=COMMA" + \
                        '&DATE_RANGE_STARTTIME={}&DATE_RANGE_ENDTIME={}'
                    dlurl = dlurl.replace(" ", "")
                    thisconf[dev['name']] = (
                        dev['type'],
                        dlurl)

            if thisconf:
                self.conf[location] = thisconf

    def read(self, requestid=None):
        endtime = datetime.now()
        starttime = endtime - timedelta(seconds=self.rate)

        for building in self.conf:
            building_path = '/' + self.to_pathname(building)
            for metername in self.conf[building].keys():
                metertype, url = self.conf[building][metername]
                meter_path = building_path + '/' + self.to_pathname(metername)
                req_url = url.format(starttime, endtime) \
                    + "&mnuStartMonth=" + str(starttime.month) + "&mnuStartDay=" + str(starttime.day) \
                    + "&mnuStartYear=" + str(starttime.year) + "&mnuStartTime=" + str(starttime.hour) + "%3A" + str(starttime.minute) \
                    + "&mnuEndMonth=" + str(endtime.month) + "&mnuEndDay=" + str(endtime.day) \
                    + "&mnuEndYear=" + str(endtime.year) + "&mnuEndTime=" + str(endtime.hour) + "%3A" + str(endtime.minute) \

                response = requests.get(url=req_url, auth=self.auth)
                if response:
                    data = response.content.decode('utf-8')
                    if "No data found within range" in data:
                        print("No data found for "+meter_path+" within given date range")
                        continue
                    csv_df = pd.read_csv(io.StringIO(data))

                    if "water" in meter_path:
                        for index, row in csv_df.iterrows():
                            rowDict = dict(row)
                            meter_data = iot_pb2.Meter(
                                water_total = types.Double(value=rowDict.get('Water (Gallons)')),
                                water_rate = types.Double(value=rowDict.get('Water Ave Rate (Gpm)')),
                                water_instantaneous = types.Double(value=rowDict.get('Water Instantaneous (Gpm)')),
                            )
                            msg = xbos_pb2.XBOS(
                                XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                                    time = int(time.time()*1e9),
                                    meter = meter_data
                                )
                            )
                            self.report(building_path + '/water_meter', msg)

                    if "condensate" in meter_path:
                        for index, row in csv_df.iterrows():
                            rowDict = dict(row)
                            meter_data = iot_pb2.Meter(
                                condense_total = types.Double(value=rowDict.get('Steam Condensate Meter (Gallons)')),
                                condense_rate = types.Double(value=rowDict.get('Steam Condensate Meter Ave Rate (Gpm)')),
                                condense_instantaneous = types.Double(value=rowDict.get('Steam Condensate Meter Instantaneous (Gpm))')),
                            )
                            msg = xbos_pb2.XBOS(
                                XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                                    time = int(time.time()*1e9),
                                    meter = meter_data
                                )
                            )
                            self.report(building_path + '/condensation_meter', msg)

                    if "electric" in meter_path:
                        for index, row in csv_df.iterrows():
                            rowDict = dict(row)
                            meter_data = iot_pb2.Meter(
                                power = types.Double(value=rowDict.get('Real Power (kW)')),
                                apparent_power = types.Double(value=rowDict.get('Apparent Power (kVA)')),
                                reactive_power = types.Double(value=rowDict.get('Reactive Power (kVAR)')),
                                voltage = types.Double(value=rowDict.get('Voltage')),
                                energy = types.Double(value=rowDict.get('Current (Amps)')),
                                demand = types.Double(value=rowDict.get('Present Demand (kW)'))
                            )
                            msg = xbos_pb2.XBOS(
                                XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                                    time = int(time.time()*1e9),
                                    meter = meter_data
                                )
                            )
                            self.report(building_path + '/electric_meter', msg)

                    else:
                        print("Unrecognized type of meter: "+meter_path)
                        continue
                else:
                    print("Received", response.status_code, "response for request:", req_url)

    def to_pathname(self, value):
        s = re.sub(r'[\W/]+', '_', value)
        s = re.sub(r'_*$', '', s)
        return s.lower()

    def remove_nbsp(self, s):
        if not s:
            return s
        s = re.sub("&nbsp;", '', s)
        return s

if __name__ == '__main__':
    with open('obvius.yaml') as f:
        driverConfig = yaml.safe_load(f)

    cfg = {
        'obvius': {
            'bmoroot': driverConfig['bmoroot'],
            'statuspage': driverConfig['statuspage'],
            'username': driverConfig['username'],
            'password': driverConfig['password'],
        },
        'wavemq': 'localhost:9516',
        'namespace': 'GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==',
        'base_resource': 'obvius',
        'entity': 'obvius.ent',
        'id': 'pyxbos-driver-obvius',
        'rate': 900, # 15 Minutes
    }

    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    obvius_driver = ObviusDriver(cfg)
    obvius_driver.begin()
