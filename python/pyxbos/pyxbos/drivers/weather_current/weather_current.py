from pyxbos import *
import os,sys
import json
import requests
import yaml
import argparse


class WeatherCurrentDriver(Driver):
    def setup(self, cfg):
        self.baseurl = cfg['darksky']['url']
        self.apikey = cfg['darksky']['apikey']
        self.coords = cfg['darksky']['coordinates']
        self.url = self.baseurl + self.apikey + '/' + self.coords

    def read(self, requestid=None):
        response = requests.get(self.url)
        json_data = json.loads(response.text)
        if 'currently' not in json_data: return

        logging.info("currently {0}".format(json_data['currently']))
        print(json_data['currently'])
        output = {}
        for key, value in json_data['currently'].items():
            output[key] = value

        #print(json_data['currently'])
        if 'humidity' in output:
            output['humidity'] *= 100 # change from decimal to percent

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                weather_current = weather_current_pb2.Weather_Current_State(
                    time  =   types.Int64(value=output.get('time',None)),
                    icon  =  output.get('icon',None),
                    nearestStormDistance  =   types.Double(value=output.get('nearestStormDistance',None)),
                    nearestStormBearing  =   types.Double(value=output.get('nearestStormBearing',None)),
                    precipIntensity  =   types.Double(value=output.get('precipIntensity',None)),
                    precipIntensityError  =   types.Double(value=output.get('precipIntensityError',None)),
                    precipProbability  =   types.Double(value=output.get('precipProbability',None)),
                    precipType  =  output.get('precipType',None),
                    temperature  =   types.Double(value=output.get('temperature',None)),
                    apparentTemperature  =   types.Double(value=output.get('apparentTemperature',None)),
                    dewPoint  =   types.Double(value=output.get('dewPoint',None)),
                    humidity  =   types.Double(value=output.get('humidity',None)),
                    pressure  =   types.Double(value=output.get('pressure',None)),
                    windSpeed  =   types.Double(value=output.get('windSpeed',None)),
                    windGust  =   types.Double(value=output.get('windGust',None)),
                    windBearing  =   types.Double(value=output.get('windBearing',None)),
                    cloudCover  =   types.Double(value=output.get('cloudCover',None)),
                    uvIndex  =   types.Double(value=output.get('uvIndex',None)),
                    visibility  =   types.Double(value=output.get('visibility',None)),
                    ozone  =   types.Double(value=output.get('ozone',None)),
                )
            )
        )
        self.report(self.coords, msg)


if __name__ == '__main__':
    with open('dark_sky.yaml') as f:
        # use safe_load instead load for security reasons
        driverConfig = yaml.safe_load(f)

    namespace = driverConfig['wavemq']['namespace']
    api = driverConfig['dark_sky']['api']
    cfg = {
        'darksky': {
            'apikey': api,
            'url': 'https://api.darksky.net/forecast/',
            'coordinates': '40.5301,-124.0000' # Should be near BLR
        },
        'wavemq': 'localhost:4516',
        'namespace': namespace,
        'base_resource': 'weather_current',
        'entity': 'weather_current.ent',
        'id': 'pyxbos-driver-current-1',
        #'rate': 1800, # half hour
        'rate': 20, # 15 min
    }
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    current_driver = WeatherCurrentDriver(cfg)
    current_driver.begin()
