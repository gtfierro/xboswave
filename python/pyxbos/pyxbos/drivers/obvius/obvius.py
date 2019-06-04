from pyxbos.driver import *
import logging
import os, re
import requests
import time
import urllib

from bs4 import BeautifulSoup as bs
from datetime import datetime, timedelta
import sensordb

try:
    import ordereddict
except ImportError:
    import collections as ordereddict

class ObviusDriver(Driver):
    def setup(self, cfg):
        self.bmoroot = cfg['obvius']['bmoroot']
        self.statuspage = cfg['obvius']['statuspage']
        self.auth = (cfg['obvius']['username'], cfg['obvius']['password'])

    def read(self, requestid=None):
        end = datetime.today()
        start = end - timedelta(days=1)
        crawl(start.strftime("%Y-%m-%d"), end.strftime("%Y-%m-%d"))

    def crawl(self, start, end, getConfig):
        starttime = start.split("-")
        endtime = end.split("-")

        # find all the AcquiSuite boxes
        devices = {}
        response = requests.get(self.bmoroot + self.statuspage, auth=self.auth)
        soup = bs(response.content, features="html.parser")

        for tr in soup.findAll('tr'):
            tds = tr('td')
            if len(tds) != 6:
                continue

            name = tds[0].a.string
            devices[name] = {
                'ip' : self.remove_nbsp(tds[3].string),
                'href' : tds[0].a['href'],
                }

         # look at all the meters hanging off each of them
         for location in devices:
             response = requests.get(self.bmoroot + devices[location]['href'], auth=self.auth)
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
             devices[location]['subdevices'] = subdevices

        conf = {}
        for location, devs in devices.items():
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
                conf[location] = thisconf

        for building in conf:
            building_path = '/' + self.to_pathname(building)
            for metername in conf[building].keys():
                meter_path = building_path + '/' + self.to_pathname(metername)
                req_url = url.format(start, end) \
                    + "&mnuStartMonth=" + starttime[1] + "&mnuStartDay=" + starttime[2] \
                    + "&mnuStartYear=" + starttime[0] + "&mnuStartTime=0%3A0" \
                    + "&mnuEndMonth=" + endtime[1] + "&mnuEndDay=" + endtime[2] \
                    + "&mnuEndYear=" + endtime[0] + "&mnuEndTime=23%3A59"

                response = requests.get(url=req_url, auth=self.auth)
                if response:
                    data = response.content.decode('utf-8')

                    if "No data found within range" in data:
                        print("[Warning] No data found for "+meter_path+" within given date range")
                        continue

                    # TODO (john-b-yang): Clean + Package Data within corresponding proto message
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
        }
        'wavemq': 'localhost:4516',
        'namespace': 'GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==',
        'base_resource': 'obvius',
        'entity': 'obvius.ent',
        'id': 'pyxbos-driver-obvius',
        'rate': 86400, # 1 Day
    }

    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    obvius_driver = ObviusDriver(cfg)
    obvius_driver.begin()
