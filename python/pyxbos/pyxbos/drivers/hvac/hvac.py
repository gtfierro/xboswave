from pyxbos.driver import *
from pyxbos import hvac_pb2
import pymortar
import BAC0
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
        self.bacnet = BAC0.connect()
        self.point_class_map = get_point_class_mapping(self.pymortar_client, cfg['building'])
        self.devices = [BAC0.device(dev[2], dev[3], self.bacnet) for dev in self.bacnet.devices]

    def read(self, requestid=None):
        for dev in self.devices:
            # A dictionary mapping point names to their bacnet values
            point_value_map = get_point_value_mapping(dev.points)

            #AHU-1.DPR-O Damper point MISSING from Bacnet points

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
                            if not isinstance(value, str):
                                props[prop] = types.Double(value=value)
                            else:
                                if value == "Cooling":
                                    props[prop] = 1
                                else:
                                    props[prop] = 0

                    data[equipclass.lower()] = getattr(hvac_pb2, equipclass)(**props)
                    data['time'] = int(time.time()*1e9)

                    msg = xbos_pb2.XBOS(XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(**data))
                    self.report(eq, msg) # Put name of equipment here as the first argument

def get_point_value_mapping(device_points):
    points = [str(point).strip() for point in device_points]
    point_class_map = {}

    for point in points:
        point_name, point_value = point.split(" : ")

        point_name = point_name.split(".", 1)[1]

        if point_value == "True":
            point_class_map[point_name] = True
        elif point_value == "False":
            point_class_map[point_name] = False
        elif any(char.isdigit() for char in point_value):
            point_class_map[point_name] = float(re.sub(r'[^0-9.]','', point_value))
        elif point_value == "Occupied":
            point_class_map[point_name] = True
        else:
            point_class_map[point_name] = point_value

    return point_class_map

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