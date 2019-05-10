from pyxbos.driver import *
from pyxbos import weather_station_pb2
import os,sys
import json
import requests
import yaml
import argparse
from pyxbos import *
import os,sys
import json
import requests
import yaml
import argparse
import logging


class WeatherPredictionDriver(Driver):
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
        output = {}

        for hour in hourly.get('data',[]):
            for key, value in hour.items():
                output[key] = value
            if 'humidity' in output:
                output['humidity'] *= 100 # change from decimal to percent
            #print(hour)
            timestamp = int(hour.get('time') * 1e9) # nanoseconds
            predictions.append(weather_station_pb2.WeatherStationPrediction.Prediction(
                prediction_time=timestamp,
                prediction=weather_station_pb2.WeatherStation(
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
            ))
            #print(predictions)

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                weather_prediction = weather_station_pb2.WeatherStationPrediction(
                    predictions=predictions
                )
            )
        )
        self.report(self.coords+'/prediction', msg)



if __name__ == '__main__':

    parser = argparse.ArgumentParser()
    parser.add_argument("config_file", help="config file with api key as well as namespace")
    parser.add_argument("ent_file", help="entity file")
    args = parser.parse_args()
    config_file = args.config_file
    ent_file = args.ent_file

    with open(config_file) as f:
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
        'base_resource': 'weather_prediction',
        'entity': ent_file,
        'id': 'pyxbos-driver-prediction-1',
        #'rate': 1800, # half hour
        'rate': 20, # 15 min
    }
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    prediction_driver = WeatherPredictionDriver(cfg)
    prediction_driver.begin()
