from pyxbos.driver import *
from pyxbos import hvac_pb2
import pymortar
#import BAC0
import re
import os,sys
import json
import time

class HVACDriver(Driver):
    def setup(self, cfg):
        # A dictionary mapping point names to their classes
        self.pymortar_client = pymortar.Client({
            'username': cfg['mortar_api_username'],
            'password': cfg['mortar_api_password']
        })
        #self.bacnet = BAC0.connect()
        self.point_class_map = get_point_class_mapping(self.pymortar_client, cfg['building'])
        #self.devices = [BAC0.device(dev[2], dev[3], cfg["bacnet"]) for dev in cfg["devices"]]

    def read(self, requestid=None):
        #for dev in self.devices:
            # A dictionary mapping point names to their bacnet values
            #point_value_map = get_point_value_mapping(dev.points)
           
            #AHU-1.DPR-O Damper point MISSING from Bacnet points

        point_value_map = get_point_value_mapping([])

        #print(point_value_map)

        output = {}
        for equipclass in self.point_class_map:
            if equipclass not in ['Boiler', 'Chilled_Water_Pump', "Hot_Water_Pump"]:
                if equipclass not in output:
                    output[equipclass] = {}

                for point in self.point_class_map[equipclass]:
                    pieces = point.split(".")
                    new_point = point
                    if len(pieces) > 2:
                        new_point = pieces[0] + "." + pieces[2]

                    if "DPRPOS" in point:
                        new_point = pieces[0] + "." + "DMPRPOS"

                    if new_point in point_value_map:
                        pt_class = self.point_class_map[equipclass][point]['point_class']
                        if pieces[0] not in output[equipclass]:
                            output[equipclass][pieces[0]] = {}
                        output[equipclass][pieces[0]][pt_class.lower()] = point_value_map[new_point]

        for equipclass in output:
            data = {}
            all_eq = output[equipclass]
            for eq in all_eq:
                props = {}
                for prop in all_eq[eq]:
                    if prop not in ['vvt', 'vav']:
                        print("EQ", eq, prop)
                        value = all_eq[eq].get(prop, None)
                        if isinstance(value, float):
                            props[prop] = types.Double(value=value)
                        elif isinstance(value, int):
                            props[prop] = types.Int64(value=value)
                        elif isinstance(value, bool):
                            props[prop] = types.Bool(value=value)
                        elif isinstance(value, str):
                            props[prop] = value

                data[equipclass.lower()] = getattr(hvac_pb2, equipclass)(**props)
                data['time'] = int(time.time()*1e9)

                msg = xbos_pb2.XBOS(XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(**data))
                print("EQUIPMENT", data)
                self.report(eq, msg) # Put name of equipment here as the first argument

                # ahu = output['AHU']['AHU-2'] #Points from AHU-2 only
                # vav = output['VAV']['VAV2-1'] #Points from VAV2-1 only
                # fan = output['Fan']['AHU-1'] #Points from AHU-1 only
                # damper = output['Damper']['VAV2-1'] #Points from VAV2-1 only
                # economizer = output['Economizer']['AHU-3'] #Points from AHU-3 only

                # msg = xbos_pb2.XBOS(
                #     XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                #         time = int(time.time()*1e9),
                #         ahu = hvac_pb2.AHU(
                #             outside_air_temperature_sensor = types.Double(value=ahu.get('outside_air_temperature_sensor', None)), 
                #             filter_status = types.Bool(value=ahu.get('filter_status', None)),
                #             discharge_air_static_pressure_setpoint = types.Double(value=ahu.get('discharge_air_static_pressure_setpoint', None)),
                #             building_static_pressure_setpoint = types.Double(value=ahu.get('building_static_pressure_setpoint', None)),
                #             heating_valve_command = types.Double(value=ahu.get('heating_valve_command', None)),
                #             occupancy_command = types.Bool(value=ahu.get('occupancy_command', None)),
                #             cooling_demand = types.Bool(value=ahu.get('cooling_demand', None)),
                #             cooling_valve_command = types.Double(value=ahu.get('cooling_valve_command', None)),
                #             discharge_air_temperature_sensor = types.Double(value=ahu.get('discharge_air_temperature_sensor', None)),
                #             return_air_temperature_sensor = types.Double(value=ahu.get('return_air_temperature_sensor', None)),
                #             building_static_pressure_sensor = types.Double(value=ahu.get('building_static_pressure_sensor', None)),
                #             discharge_air_temperature_setpoint = types.Double(value=ahu.get('discharge_air_temperature_setpoint', None)),
                #             mixed_air_temperature_sensor = types.Double(value=ahu.get('mixed_air_temperature_setpoint', None)),
                #             discharge_air_static_pressure_sensor = types.Double(value=ahu.get('discharge_air_static_pressure_sensor', None)),
                #             zone_temperature_sensor = types.Double(value=ahu.get('zone_temperature_sensor', None)),
                #             shutdown_command = types.Bool(value=ahu.get('shutdown_command', None)),
                #             supply_air_damper_min_position_setpoint = types.Double(value=ahu.get('supply_air_damper_min_position_setpoint', None)),
                #             mixed_air_temperature_low_limit_setpoint = types.Double(value=ahu.get('mixed_air_temperature_low_limit_setpoint', None)),
                #             mixed_air_temperature_setpoint = types.Double(value=ahu.get('mixed_air_temperature_setpoint', None)),
                #             occupied_mode_status = types.Bool(value=ahu.get('occupied_mode_status', None)),
                #             zone_temperature_setpoint = types.Double(value=ahu.get('mixed_air_temperature_setpoint', None))
                #         ),
                #         vav = hvac_pb2.VAV(
                #             supply_air_flow_sensor = types.Double(value=vav.get('supply_air_flow_sensor', None)),
                #             zone_temperature_sensor = types.Double(value=vav.get('zone_temperature_sensor', None)),
                #             zone_temperature_setpoint = types.Double(value=vav.get('zone_temperature_setpoint', None)),
                #             occupied_heating_min_supply_air_flow_setpoint = types.Double(value=vav.get('occupied_heating_min_supply_air_flow_setpoint', None)),
                #             thermostat_adjust_setpoint = types.Double(value=vav.get('thermostat_adjust_setpoint', None)),
                #             supply_air_velocity_pressure_sensor = types.Double(value=vav.get('supply_air_velocity_pressure_sensor', None)),
                #             discharge_air_temperature_sensor = types.Double(value=vav.get('discharge_air_temperature_sensor', None)),
                #             cooling_max_supply_air_flow_setpoint = types.Double(value=vav.get('cooling_max_supply_air_flow_setpoint', None)),
                #             box_mode = vav.get('box_mode', None),
                #             cooling_demand = types.Double(value=vav.get('cooling_demand', None)),
                #             heating_demand = types.Double(value=vav.get('heating_demand', None)),
                #             supply_air_flow_setpoint = types.Double(value=vav.get('supply_air_flow_setpoint', None))
                #         ),
                #         fan = hvac_pb2.Fan(
                #             vfd_alarm = types.Bool(value=fan.get('vfd_alarm', None)),
                #             fan_status = types.Bool(value=fan.get('fan_status', None)),
                #             fan_reset_command = types.Bool(value=fan.get('fan_reset_command ', None)),
                #             fan_overload_alarm = types.Bool(value=fan.get('fan_overload_alarm', None)),
                #             on_off_command = types.Bool(value=fan.get('on_off_command', None)),
                #             fan_speed_setpoint = types.Double(value=fan.get('fan_speed_setpoint', None))
                #         ),
                #         economizer = hvac_pb2.Economizer(
                #             economizer_differential_air_temperature_setpoint = types.Double(value=economizer.get('economizer_differential_air_temperature_setpoint', None))
                #         ),
                #         damper = hvac_pb2.Damper(
                #             damper_position_command = types.Double(value=damper.get('damper_position_command', None)),
                #             damper_position_sensor = types.Double(value=damper.get('damper_position_sensor', None))
                #         )
                #     )
                #)

def get_point_value_mapping(device_points):
    # points = [str(point).strip() for point in device_points]
    # point_class_map = {}

    # for point in points:
    #     point_name, point_value = point.split(" : ")

    #     point_name = point_name.split(".", 1)[1]

    #     if point_value == "True":
    #         point_class_map[point_name] = True
    #     elif point_value == "False":
    #         point_class_map[point_name] = False
    #     elif any(char.isdigit() for char in point_value):
    #         point_class_map[point_name] = float(re.sub(r'[^0-9.]','', point_value))
    #     elif point_value == "Occupied":
    #         point_class_map[point_name] = True
    #     else:
    #         point_class_map[point_name] = point_value
    
    with open('point_to_value.json', 'r') as f:
        point_class_map = json.load(f)
        return point_class_map

    #return point_class_map

def get_point_class_mapping(pymortar_client, building):

    v = pymortar.View(
            name="equipment",
            definition="""
            SELECT ?equipname ?equipclass ?point ?pointclass FROM %s WHERE {
                ?equipclass rdfs:subClassOf+ brick:Equipment .
                ?equipname rdf:type ?equipclass .
                ?equipname bf:hasPoint ?point .
                ?point rdf:type ?pointclass
            };""" % building
        )

    res = pymortar_client.fetch(pymortar.FetchRequest(
            sites=[building],
            views=[v]
        ))

    eq_view = res.view('equipment')

    point_class_map = {}
    groups = eq_view.groupby('equipclass').groups
    clustered_equipment = {"Damper": None, "Fan": None}
    clustered = None

    for equipclass, indexes in groups.items():
        for clustered_eq in clustered_equipment:
            # Union all the clustered equipment
            if clustered_eq.lower() in equipclass.lower():
                clustered = clustered_eq
                break
        
        if clustered:
            if clustered not in point_class_map:
                point_class_map[clustered] = {}
        else:
            if equipclass not in point_class_map:
                point_class_map[equipclass] = {}
                
        points = eq_view.iloc[indexes][['point','pointclass']].values
        for p in points:
            point_name, point_class = p[0], p[1]
            if clustered:
                point_class_map[clustered][point_name] = { "point_name": point_name, "point_class": point_class}
            else:
                point_class_map[equipclass][point_name] = { "point_name": point_name, "point_class": point_class}

        clustered = None
    
    return point_class_map

if __name__ == '__main__':
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

    cfg = {
        'wavemq': 'localhost:4516',
        'waved': 'localhost:410',
        'namespace': 'GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg==',
        'base_resource': 'sriharsha',
        'entity': 'sriharsha.ent',
        'id': 'pyxbos-driver-hvac-1',
        'rate': 10,
        'building': 'orinda-public-library',
        'mortar_api_username': os.environ['MORTAR_API_USERNAME'],
        'mortar_api_password': os.environ['MORTAR_API_PASSWORD']
    }

    current_driver = HVACDriver(cfg)
    current_driver.begin()