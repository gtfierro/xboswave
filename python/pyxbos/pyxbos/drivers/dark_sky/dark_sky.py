from pyxbos.driver import *
from pyxbos import weather_station_pb2
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
        print("In prediction driver")
        response = requests.get(self.url)
        json_data = json.loads(response.text)
        if 'hourly' not in json_data: return

        hourly = json_data['hourly']
        #print(json_data)
        predictions = []

        for hour in hourly.get('data',[]):
            timestamp = int(hour.get('time') * 1e9) # nanoseconds
            temperature = hour.get('apparentTemperature', None)
            precipIntensity = hour.get('precipIntensity', None)
            precipProbability = hour.get('precipProbability', None)
            humidity = hour.get('humidity', None)
            if humidity is not None:
                humidity *= 100 # change from decimal to percent

            predictions.append(weather_station_pb2.WeatherStationPrediction.Prediction(
                prediction_time=timestamp,
                prediction=weather_station_pb2.WeatherStation(
                    temperature=types.Double(value=temperature),
                    precipIntensity=types.Double(value=precipIntensity),
                    humidity=types.Double(value=humidity),
                )
            ))

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                weather_prediction = weather_station_pb2.WeatherStationPrediction(
                    predictions=predictions
                )
            )
        )
        self.report(self.coords+'/prediction', msg)

class DarkSkyDriver(Driver):
    def setup(self, cfg):
        self.baseurl = cfg['darksky']['url']
        self.apikey = cfg['darksky']['apikey']
        self.coords = cfg['darksky']['coordinates']
        self.url = self.baseurl + self.apikey + '/' + self.coords

    def read(self, requestid=None):
        print("In current driver")
        response = requests.get(self.url)
        json_data = json.loads(response.text)
        if 'currently' not in json_data: return

        logging.info("currently {0}".format(json_data['currently']))
        nearestStormDistance =  json_data['currently'].get('nearestStormDistance',None)
        nearestStormBearing =   json_data['currently'].get('nearestStormBearing',None)
        precipIntensity =       json_data['currently'].get('precipIntensity',None)
        apparentTemperature =   json_data['currently'].get('apparentTemperature',None)
        humidity =              json_data['currently'].get('humidity',None)
        #print(json_data['currently'])
        if humidity is not None:
            humidity *= 100 # change from decimal to percent

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                weather_current = weather_station_pb2.WeatherStation(
                    nearestStormDistance  =   types.Double(value=nearestStormDistance),
                    nearestStormBearing   =   types.Double(value=nearestStormBearing),
                    precipIntensity        =   types.Double(value=precipIntensity),
                    temperature             =   types.Double(value=apparentTemperature),
                    humidity                =   types.Double(value=humidity),
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
        'base_resource': 'dark_sky',
        'entity': 'dark_sky.ent',
        'id': 'pyxbos-driver-darksky-1',
        #'rate': 1800, # half hour
        'rate': 20, # 15 min
    }
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    current_driver = DarkSkyDriver(cfg)
    prediction_driver = DarkSkyPredictionDriver(cfg)
    current_driver.begin()
    prediction_driver.begin()
