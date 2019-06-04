import os
import re

import argparse
import configparser
import io
import json
import urllib
import requests
from bs4 import BeautifulSoup as bs
import pandas as pd

from auth import *
import sensordb

try:
    import ordereddict
except ImportError:
    import collections as ordereddict

def to_pathname(value):
    s = re.sub(r'[\W/]+', '_', value)
    s = re.sub(r'_*$', '', s)
    return s.lower()

def remove_nbsp(s):
    if not s:
        return s
    s = re.sub("&nbsp;", '', s)
    return s

def createFolder(directory):
    try:
        if not os.path.exists(directory):
            os.makedirs(directory)
    except OSError:
        raise OSError("Error encountered when creating directory")

def crawler(start, end, getDevices, getConfig):
    starttime = start.split("-")
    endtime = end.split("-")

    datadir = "data"
    createFolder(datadir)

    # find all the AcquiSuite boxes
    devices = {}
    response = requests.get(BMOROOT + STATUSPAGE, auth=AUTH)
    soup = bs(response.content, features="html.parser")

    for tr in soup.findAll('tr'):
        tds = tr('td')
        if len(tds) != 6:
            continue

        name = tds[0].a.string
        devices[name] = {
            'ip' : remove_nbsp(tds[3].string),
            'href' : tds[0].a['href'],
            }

    # look at all the meters hanging off each of them
    for location in devices:
        print("Location: ", location, " URL: ", BMOROOT + devices[location]['href'])
        response = requests.get(BMOROOT + devices[location]['href'], auth=AUTH)
        soup = bs(response.content, features="html.parser")
        subdevices = []
        for tr in soup.findAll('tr'):
            tds = tr('td')
            if len(tds) != 5 or tds[3].a != None:
                continue
            subdevices.append({
                'address' : re.sub("<.*?>", "", str(tds[0])),
                'status': remove_nbsp(tds[1].string),
                'name' : remove_nbsp(tds[2].string),
                'type' : remove_nbsp(tds[3].string),
                'firmware': remove_nbsp(tds[4].string)
            })
        devices[location]['subdevices'] = subdevices

    if getDevices:
        devices_json = open("devices.json", "w")
        with open("devices.json", "w") as out:
            out.write(json.dumps(devices))
        devices_json.close()
        print("Exported Device Data to JSON")

    conf = {}
    for location, devs in devices.items():
        params = urllib.parse.parse_qs(urllib.parse.urlsplit(devs['href']).query)
        if "AS" not in params or "DB" not in params:
            continue
        if location in AUTH:
            continue
        thisconf = {}
        for dev in devs['subdevices']:
            if sensordb.get_map(dev['type'], location) != None:
                dlurl = BMOROOT + 'mbdev_export.php/' + params['AS'][0] + '_' +  \
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

    # generate config file
    cf = configparser.RawConfigParser('', ordereddict.OrderedDict)
    cf.optionxform = str

    cf.add_section('server')
    cf.set('server', 'SuggestThreadPool', '20')
    cf.set('server', 'Port', '9051')
    cf.add_section('/')
    cf.set('/', 'Metadata/Location/Campus', 'UCB')
    cf.set('/', 'Metadata/SourceName', 'buildingmanageronline archive')
    cf.set('/', 'uuid', '91dde108-d02b-11e0-8542-0026bb56ec92')

    requests_total = 0
    requests_failed = 0
    requests_nodata = 0

    for building in conf:
        building_path = '/' + to_pathname(building)
        cf.add_section(building_path)
        cf.set(building_path, 'type', 'Collection')
        for metername in conf[building].keys():
            metertype, url = conf[building][metername]

            building_name = building
            if "New" in building_name:
                building_name = building_name[:building_name.index("New")]
            if "NEW" in building_name:
                building_name = building_name[:building_name.index("NEW")]

            meter_path = building_path + '/' + to_pathname(metername)
            cf.add_section(meter_path)
            cf.set(meter_path, 'Metadata/Extra/MeterName', metername)
            cf.set(meter_path, 'Metadata/Instrument/Model', '"' + metertype + '"')
            cf.set(meter_path, 'Metadata/Location/Building', building_name)
            cf.set(meter_path, 'Url', url)
            # add any extra config options specific to this meter type
            sensor_map = sensordb.get_map(metertype, building_name)
            if 'extra' in sensor_map:
                for k, val in sensor_map['extra'].items():
                    cf.set(meter_path, k, val)

            req_url = url.format(start, end) \
                + "&mnuStartMonth=" + starttime[1] \
                + "&mnuStartDay=" + starttime[2] \
                + "&mnuStartYear=" + starttime[0] \
                + "&mnuStartTime=0%3A0" \
                + "&mnuEndMonth=" + endtime[1] \
                + "&mnuEndDay=" + endtime[2] \
                + "&mnuEndYear=" + endtime[0] \
                + "&mnuEndTime=23%3A59"

            response = requests.get(url=req_url, auth=AUTH)
            requests_total += 1
            if response:
                data = response.content.decode('utf-8')

                if "No data found within range" in data:
                    print("[Warning] No data found for "+meter_path+" within given date range")
                    requests_nodata += 1
                    continue

                filename = meter_path.replace(r"[\[\]]", "")
                filename = datadir + "/" + filename.strip("/").replace("/", "-")
                filename += "-" + start + "to" + end
                filename += ".csv"

                try:
                    csv_df = pd.read_csv(io.StringIO(data))
                    csv_df.to_csv(filename, index=False)
                    print("Downloaded Data for", meter_path)
                except pd.errors.ParserError as error:
                    print("[Error] Failed to download file "+meter_path+". Error: "+error)
                    requests_failed += 1
            else:
                requests_failed += 1
                print("Received", response.status_code, "response for request:", req_url)

    if getConfig:
        config_file = open("config.ini", "w")
        cf.write(config_file)
        config_file.close()

    print("\n-------------")
    print("# of Requests - No Data: ", requests_nodata)
    print("# of Requests - Failed: ", requests_failed)
    print("# of Requests - Total: ", requests_total)
    no_data_percent = float(requests_nodata)/float(requests_total)*100
    failed_percent = float(requests_failed)/float(requests_total)*100
    print(f"Percent of Requests - No Data: {no_data_percent:.3f}")
    print(f"Percent of Requests - Failed: {failed_percent:.3f}")

if __name__ == '__main__':
    # Defining Arguments
    parser = argparse.ArgumentParser(description="Specify Start and End \
                        Dates of Building Data Query")
    parser.add_argument("-s", "--start", required=True, type=str, metavar="YYYY-MM-DD",
                        help="[Required] Start Date of Query (YYYY-MM-DD)")
    parser.add_argument("-e", "--end", required=True, type=str, metavar="YYYY-MM-DD",
                        help='[Required] End Date of Query (YYYY-MM-DD)')
    parser.add_argument("-d", "--devices", required=False, type=bool,
                        default=False, help="[Optional] Download JSON list of buildings + meters")
    parser.add_argument("-c", "--config", required=False, type=bool,
                        default=False, help="[Optional] Download config file of buildings + meters")
    ARGS = parser.parse_args()

    crawler(ARGS.start, ARGS.end, ARGS.devices, ARGS.config)
