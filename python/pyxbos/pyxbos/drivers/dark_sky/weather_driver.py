from pyxbos.driver import *
from pyxbos import weather_station_pb2
import pandas as pd
from pvlib import solarposition, irradiance
import os,sys
import json
import requests
import yaml
import argparse
import logging
import numpy as np

class WeatherDriver(Driver):
    def setup(self, cfg):
        self.baseurl = cfg['darksky']['url']
        self.apikey = cfg['darksky']['apikey']
        self.coords = cfg['darksky']['coordinates']
        self.lat = float(self.coords.split(',')[0])
        self.lng = float(self.coords.split(',')[1])
        self.url = self.baseurl + self.apikey + '/' + self.coords
        self.service_name = cfg['service_name']

    def solar_model_ZhHu(self, forecast_df, sin_alt, zh_solar_const):
        '''
        Estimate Global Horizontal Irradiance (GHI) from Zhang-Huang solar forecast model
        Params: sin_alt, sine of solar altitude
                forecast_df should include the following columns:
                cloudCover: [0,1];
                temperature: degC;
                relative humidity: %;
                windSpeed: m/s;
        Returns: estimated GHI
        '''

        # df = forecast_df.copy()

        c0 = 0.5598
        c1 = 0.4982
        c2 = -0.6762
        c3 = 0.02842
        c4 = -0.00317
        c5 = 0.014
        d = -17.853
        k = 0.843

        forecast_df['temperature_c'] = (forecast_df['temperature'] - 32) * 5.0/9.0
        shift_temp = pd.Series(data=np.roll(forecast_df.temperature_c, 3), index=forecast_df.index)
        forecast_df['deltaT'] = forecast_df['temperature_c'] - shift_temp
        forecast_df['estimatedGhi'] = (zh_solar_const * sin_alt * (c0 + c1 * forecast_df['cloudCover']
                                                                    + c2 * forecast_df['cloudCover'] ** 2 + c3 *
                                                                    forecast_df['deltaT'] + c4 * forecast_df[
                                                                        'humidity'] * 100
                                                                    + c5 * forecast_df['windSpeed']) + d) / k
        forecast_df.loc[forecast_df.loc[forecast_df.estimatedGhi <= 0].index, 'estimatedGhi'] = 0

        return forecast_df

    def Perez_split(self, forecast_df):
        '''
        Estimate beam radiation and diffuse radiation from GHI and solar altitude
        Params: forecast_df includes the estimated GHI
        Returns: beam_rad, beam radiation, W/m2
                 diff_rad, diffuse radiation, W/m2
        '''

        # datetime should include timezone information, otherwise UTC time by default
        alt_ang = solarposition.get_solarposition(forecast_df.index, self.lat, self.lng)['elevation']
        sin_alt = np.sin(np.radians(alt_ang))
        zh_solar_const = 1355  # W/m2, solar constant used by Zhang-Huang model
        solar_const = 1367  # general solar constant

        df = self.solar_model_ZhHu(forecast_df=forecast_df, sin_alt=sin_alt, zh_solar_const=zh_solar_const)

        clear_index_kt = df['estimatedGhi'] / (solar_const * sin_alt)
        clear_index_ktc = 0.4268 + 0.1934 * sin_alt

        diff = (clear_index_kt < clear_index_ktc) * 1  # *1 converts boolean to integer
        clear_index_kds = diff * ((3.996 - 3.862 * sin_alt + 1.54 * (sin_alt) ** 2) * (clear_index_kt) ** 3) + \
                          (1 - diff) * (clear_index_kt - (1.107 + 0.03569 * sin_alt + 1.681 * (sin_alt) ** 2) * (
                    1.0 - clear_index_kt) ** 3)

        # Calculate direct normal radiation, W/m2
        df['beamRadiation'] = zh_solar_const * sin_alt * clear_index_kds * (1.0 - clear_index_kt) / (1.0 - clear_index_kds)
        # Calculate diffuse horizontal radiation, W/m2
        df['diffuseRadiation'] = zh_solar_const * sin_alt * (clear_index_kt - clear_index_kds) / (1.0 - clear_index_kds)
        return df

    def plane_of_array(self, df):
        """
        :param df: data frame includes GHI, beamRadiation, diffRadiation
        :return: df with plane of array solar radiation on pv and windows
        """
        pv_tilt = 8
        pv_azimuth = 37
        albedo = 0.2
        win_tilt = 90
        win_azimuth = 0
        datetime = df.index
        alt_ang = solarposition.get_solarposition(datetime, self.lat, self.lng)['elevation']
        azi_ang = solarposition.get_solarposition(datetime, self.lat, self.lng)['azimuth']
        df['poaSrOnPV'] = irradiance.get_total_irradiance(pv_tilt, pv_azimuth, alt_ang, azi_ang, df['beamRadiation'],
                                                       df['estimatedGhi'], df['diffuseRadiation'], albedo)['poa_global']
        df['poaSrOnWindows'] = irradiance.get_total_irradiance(win_tilt, win_azimuth, alt_ang, azi_ang, df['beamRadiation'],
                                                        df['estimatedGhi'], df['diffuseRadiation'], albedo)['poa_global']
        return df

    def read(self, requestid=None):
        try:
            response = requests.get(self.url)
            json_data = json.loads(response.text)

            if 'hourly' not in json_data:
                return
            else:
                hourly_df = pd.DataFrame.from_dict(json_data['hourly']['data'])
                hourly_df.time = pd.to_datetime(hourly_df.time, unit='s', utc=True)
                hourly_df = hourly_df.set_index('time')

                hourly_df = self.Perez_split(forecast_df=hourly_df)
                hourly_df = self.plane_of_array(df=hourly_df)
                hourly_df = hourly_df.reset_index()
                hourly_df['time'] = hourly_df['time'].astype(int)
                hourly_dict = hourly_df.drop(columns=['precipType']).to_dict('records')
                json_data['hourly']['data'] = hourly_dict

            hourly = json_data['hourly']
            predictions = []
            hourly_output = {}

            for hour in hourly.get('data',[]):
                for key, value in hour.items():
                    hourly_output[key] = value

                if 'humidity' in hourly_output:
                    hourly_output['humidity'] *= 100

                timestamp = int(hourly_output.get('time',None))
                predictions.append(weather_station_pb2.WeatherStationPrediction.Prediction(
                    prediction_time=timestamp,
                    prediction=weather_station_pb2.WeatherStation(
                        time  =   types.Int64(value=timestamp),
                        icon  =  hourly_output.get('icon',None),
                        nearestStormDistance  =   types.Double(value=hourly_output.get('nearestStormDistance',None)),
                        nearestStormBearing  =   types.Double(value=hourly_output.get('nearestStormBearing',None)),
                        precipIntensity  =   types.Double(value=hourly_output.get('precipIntensity',None)),
                        precipIntensityError  =   types.Double(value=hourly_output.get('precipIntensityError',None)),
                        precipProbability  =   types.Double(value=hourly_output.get('precipProbability',None)),
                        precipType  =  hourly_output.get('precipType',None),
                        temperature  =   types.Double(value=hourly_output.get('temperature',None)),
                        apparentTemperature  =   types.Double(value=hourly_output.get('apparentTemperature',None)),
                        dewPoint  =   types.Double(value=hourly_output.get('dewPoint',None)),
                        humidity  =   types.Double(value=hourly_output.get('humidity',None)),
                        pressure  =   types.Double(value=hourly_output.get('pressure',None)),
                        windSpeed  =   types.Double(value=hourly_output.get('windSpeed',None)),
                        windGust  =   types.Double(value=hourly_output.get('windGust',None)),
                        windBearing  =   types.Double(value=hourly_output.get('windBearing',None)),
                        cloudCover  =   types.Double(value=hourly_output.get('cloudCover',None)),
                        uvIndex  =   types.Double(value=hourly_output.get('uvIndex',None)),
                        visibility  =   types.Double(value=hourly_output.get('visibility',None)),
                        ozone  =   types.Double(value=hourly_output.get('ozone',None)),
                        estimatedGhi  =   types.Double(value=hourly_output.get('estimatedGhi',None)),
                        beamRadiation  =   types.Double(value=hourly_output.get('beamRadiation',None)),
                        diffuseRadiation  =   types.Double(value=hourly_output.get('diffuseRadiation',None)),
                        poaSrOnPV =   types.Double(value=hourly_output.get('poaSrOnPV',None)),
                        poaSrOnWindows  =   types.Double(value=hourly_output.get('poaSrOnWindows',None)),
                    )
                ))

            time_now = int(time.time() * 1e9)

            hourly_msg = xbos_pb2.XBOS(
                XBOSIoTDeviceState=iot_pb2.XBOSIoTDeviceState(
                    time=time_now,
                    weather_station_prediction=weather_station_pb2.WeatherStationPrediction(
                        predictions=predictions
                    )
                )
            )
            self.report(self.service_name+ '/prediction', hourly_msg)

            if 'currently' not in json_data: return
            currently_output = {}
            for key, value in json_data['currently'].items():
                currently_output[key] = value

            if 'humidity' in currently_output:
                currently_output['humidity'] *= 100 # change from decimal to percent

            currently_msg = xbos_pb2.XBOS(
                XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                    time = time_now,
                    weather_station = weather_station_pb2.WeatherStation(
                        time  =   types.Int64(value=currently_output.get('time',None)),
                        icon  =  currently_output.get('icon',None),
                        nearestStormDistance  =   types.Double(value=currently_output.get('nearestStormDistance',None)),
                        nearestStormBearing  =   types.Double(value=currently_output.get('nearestStormBearing',None)),
                        precipIntensity  =   types.Double(value=currently_output.get('precipIntensity',None)),
                        precipIntensityError  =   types.Double(value=currently_output.get('precipIntensityError',None)),
                        precipProbability  =   types.Double(value=currently_output.get('precipProbability',None)),
                        precipType  =  currently_output.get('precipType',None),
                        temperature  =   types.Double(value=currently_output.get('temperature',None)),
                        apparentTemperature  =   types.Double(value=currently_output.get('apparentTemperature',None)),
                        dewPoint  =   types.Double(value=currently_output.get('dewPoint',None)),
                        humidity  =   types.Double(value=currently_output.get('humidity',None)),
                        pressure  =   types.Double(value=currently_output.get('pressure',None)),
                        windSpeed  =   types.Double(value=currently_output.get('windSpeed',None)),
                        windGust  =   types.Double(value=currently_output.get('windGust',None)),
                        windBearing  =   types.Double(value=currently_output.get('windBearing',None)),
                        cloudCover  =   types.Double(value=currently_output.get('cloudCover',None)),
                        uvIndex  =   types.Double(value=currently_output.get('uvIndex',None)),
                        visibility  =   types.Double(value=currently_output.get('visibility',None)),
                        ozone  =   types.Double(value=currently_output.get('ozone',None)),
                    )
                )
            )
            self.report(self.service_name, currently_msg)
        except:
            print("error occured! continuing")

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("config_file", help="config file with api key as well as namespace")
    args = parser.parse_args()
    config_file = args.config_file

    with open(config_file) as f:
        driverConfig = yaml.safe_load(f)

    xbosConfig = driverConfig['xbos']
    waved = xbosConfig.get('waved', 'localhost:777')
    wavemq = xbosConfig.get('wavemq', 'locahost:4516')
    namespace = xbosConfig.get('namespace')
    base_resource = xbosConfig.get('base_resource')
    service_name = xbosConfig.get('service_name')
    entity = xbosConfig.get('entity')
    rate = xbosConfig.get('rate')
    driver_id = xbosConfig.get('id', 'darksky-driver')

    darkskyConfig = driverConfig['dark_sky']
    url = darkskyConfig.get('url', 'https://api.darksky.net/forecast/')
    api = darkskyConfig.get('api')
    coordinates = darkskyConfig.get('latlong')

    xbos_cfg = {
        'waved': waved,
        'wavemq': wavemq,
        'namespace': namespace,
        'base_resource': base_resource,
        'entity': entity,
        'id': driver_id,
        'rate': rate, 
        'service_name': service_name,
        
        'darksky': {
            'apikey': api,
            'url': url,
            'coordinates': coordinates
        }
    }

    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    weather_driver = WeatherDriver(xbos_cfg)
    weather_driver.begin()
