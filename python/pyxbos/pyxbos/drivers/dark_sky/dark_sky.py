from pyxbos import *
import os,sys
import json
import requests
import yaml
import argparse


class DarkSkyPredictionDriver(Driver):
    def setup(self, cfg):
        self.baseurl = cfg['darksky']['url']
        self.apikey = cfg['darksky']['apikey']
        self.coords = cfg['darksky']['coordinates']
        self.url = self.baseurl + self.apikey + '/' + self.coords

    def read(self, requestid=None):
        response = requests.get(self.url)
        json_data = json.loads(response.text)
        logging.info("json {0}".format(json_data))
        #TODO: finish this
        
class DarkSkyDriver(Driver):
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
        nearestStormDistance =  json_data['currently'].get('nearestStormDistance',None)
        nearestStormBearing =   json_data['currently'].get('nearestStormBearing',None)
        precipIntensity =       json_data['currently'].get('precipIntensity',None)
        apparentTemperature =   json_data['currently'].get('apparentTemperature',None)
        humidity =              json_data['currently'].get('humidity',None)

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                weather_station = iot_pb2.WeatherStation(
                    nearest_storm_distance  =   types.Double(value=nearestStormDistance),
                    nearest_storm_bearing   =   types.Int32(value=nearestStormBearing),
                    precip_intensity        =   types.Double(value=precipIntensity),
                    temperature             =   types.Double(value=apparentTemperature),
                    humidity                =   types.Double(value=humidity),
                )
            )
        )
        self.report('blr1', msg)


cfg = {
    'darksky': {
        'apikey': '<api key here>',
        'url': 'https://api.darksky.net/forecast/',
        'coordinates': '40.5301,-124.0000' # Should be near BLR
    },
    'wavemq': 'localhost:4516',
    'namespace': 'GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==',
    'base_resource': 'test/darksky',
    'entity': 'gabedarksky.ent',
    'id': 'pyxbos-driver-darksky-1',
    #'rate': 1800, # half hour
    'rate': 900, # 15 min
}
logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
e = DarkSkyDriver(cfg)
e.begin()
