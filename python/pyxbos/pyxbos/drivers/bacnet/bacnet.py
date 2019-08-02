from pyxbos.driver import *
from pyxbos import hvac_pb2
import pymortar
import BAC0
import re
import os,sys
import json
import time

class BACnetDriver(Driver):
    def setup(self, cfg):
        self.pymortar_client = pymortar.Client({
            'username': cfg['mortar_api_username'],
            'password': cfg['mortar_api_password']
        })
        self.bacnet = BAC0.connect()
        # A dictionary mapping BACnet point names to their Brick classes grouped by Equipment class
        self.point_class_dict = get_point_class_dict(self.pymortar_client, cfg['building'])
        # A list of all the BACnet devices in the network
        self.devices = [BAC0.device(dev[2], dev[3], self.bacnet) for dev in self.bacnet.devices]

    def read(self, requestid=None):
        for dev in self.devices:
            # A dictionary mapping BACnet point names mapped their respective values
            point_value_map = get_point_value_dict(dev.points)

            #AHU-1.DPR-O Damper point MISSING from BACnet points for orinda-public-library

            # Output will look like { "AHU": { "AHU2-1": { "return_air_temperature_sensor": 67.47, ..... }, ....}, .....}
            output = {}
            for equipclass in self.point_class_dict:
                if equipclass not in output:
                    output[equipclass] = {}

                for point in self.point_class_dict[equipclass]:
                    pieces = point.split(".")
                    new_point = point
                    if len(pieces) > 2:
                        # This is to fix point names like "AHU-2.SF-.SF-O" to "AHU-2.SF-O"
                        new_point = pieces[0] + "." + pieces[-1]

                    if "DPRPOS" in point:
                        # Pymortar may return point names with DPRPOS while BACnet's naming convention is DMPRPOS
                        # This fixes those points to follow BACnet's naming convention
                        new_point = pieces[0] + "." + "DMPRPOS"

                    if new_point in point_value_map:
                        # Gets the Brick class of the point name using the point to class name dict
                        point_class = self.point_class_dict[equipclass][point]['point_class']
                        
                        if pieces[0] not in output[equipclass]:
                            output[equipclass][pieces[0]] = {}
                        
                        output[equipclass][pieces[0]][point_class.lower()] = point_value_map[new_point]
            
            # Goes through every equipment class in output, and sends out a xbos_pb2 message
            for equipclass in output:
                # Maps all properties of each equipment to their respective BACnet values
                # Looks like { "time": 1564278587100909056, ahu": hvac_pb2.AHU(return_air_temperature_sensor = 67.47, ....), .... }
                device_state_data = {}
                # all_eq looks like { "AHU2-1": {...}, "AHU2-3": {....}, .... }
                all_eq = output[equipclass]
                for eq in all_eq:
                    # Dictionary of properties for this piece of equipment
                    # Looks like { "return_air_temperature_sensor": 67.47, .... }
                    props = {}
                    for prop in all_eq[eq]:
                        # For each property like "return_air_temperature_sensor"
                        if prop not in ['vvt', 'vav']:
                            # Gets the value of that property or sets it to None if it doesn't exist
                            value = all_eq[eq].get(prop, None)

                            if not isinstance(value, str):
                                # Changes the value into a Double type (defined in the proto file) if not a string
                                # If the point_value_dict was created correctly, then all values should be converted into doubles already
                                props[prop] = types.Double(value=value)

                    # e.g. Calls hvac_pb2.AHU(return_air_temperature_sensor = 67.47)
                    device_state_data[equipclass.lower()] = getattr(hvac_pb2, equipclass)(**props)
                    device_state_data['time'] = int(time.time()*1e9)

                    msg = xbos_pb2.XBOS(XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(**device_state_data))
                    self.report(eq, msg)

def get_point_value_dict(device_points):
    """ Creates a dictionary of BACnet point names mapped to their respective values

    e.g. { "AHU-1.SF-A": false, "VVT-4.MAXDMP": 75, "Garage.ZN2-Q": 1.77, .... }

    :param device_points: List of all the point names in a particular BACnet device
    :return: Dictionary of BACnet point names mapped to their respective values
    """
    points = [str(point).strip() for point in device_points]
    point_value_dict = {}

    for point in points:
        point_name, point_value = point.split(" : ")

        point_name = point_name.split(".", 1)[1]

        if point_value.lower() in ["true", "occupied", "cooling"]:
            # If the value is "occupied" or "cooling" then, converting the value to a boolean True to
            # make it easier to track
            point_value_dict[point_name] = True
        elif point_value.lower() in ["false", "unknown", "unoccupied"]:
            # If the value is "unknown" or "unoccupied" then, converting the value to a boolean False to
            # make it easier to track
            point_value_dict[point_name] = False
        elif any(char.isdigit() for char in point_value):
            # If the BACnet value contains any numbers, get rid of all 
            # non-digit characters and convert to float
            point_value_dict[point_name] = float(re.sub(r'[^0-9.]','', point_value))
        else:
            point_value_dict[point_name] = point_value

    return point_value_dict

def get_point_class_dict(pymortar_client, building):
    """ Creates a dictionary of BACnet point names mapped to their respective Brick classes.
    Points are also grouped by Equipment class. Point names and classes are obtained from Pymortar. 

    e.g. { "AHU": { 
            "AHU-3.RA-T": {
                "point_name": "AHU-3.RA-T",
                "point_class": "Return_Air_Temperature_Sensor"
            }, .....,
          "VAV": { 
            .....
           }
        }

    :param pymortar_client: Pymortar client
    :param building: Name of building that is being analyzed
    :return: Dictionary of BACnet point names mapped to their respective Brick classes grouped by Equipment class
    """
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

    point_class_dict = {}
    groups = eq_view.groupby('equipclass').groups
    # Equipment classes that are not being considered and will be excluded from the dictionary
    # Add to this list of excluded equipment for other buildings
    excluded_equipment = ['Boiler', 'Chilled_Water_Pump', 'Hot_Water_Pump']
    # Clusters all points with subclasses of Fan and Damper (i.e. Exhaust_Fan)
    clustered_equipment = dict.fromkeys(["Damper", "Fan"], None)
    clustered = None

    for equipclass, indexes in groups.items():
        if equipclass not in excluded_equipment:
            for clustered_eq in clustered_equipment:
                # Union all the clustered equipment
                if clustered_eq.lower() in equipclass.lower():
                    clustered = clustered_eq
                    break

            if clustered:
                if clustered not in point_class_dict:
                    point_class_dict[clustered] = {}
            else:
                if equipclass not in point_class_dict:
                    point_class_dict[equipclass] = {}

            points = eq_view.iloc[indexes][['point','pointclass']].values
            for p in points:
                point_name, point_class = p[0], p[1]
                if clustered:
                    point_class_dict[clustered][point_name] = { "point_name": point_name, "point_class": point_class}
                else:
                    point_class_dict[equipclass][point_name] = { "point_name": point_name, "point_class": point_class}

            clustered = None

    return point_class_dict

if __name__ == '__main__':
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

    building = 'orinda-public-library'

    cfg = {
        'wavemq': 'localhost:4516',
        'waved': 'localhost:410',
        'namespace': 'GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg==',
        'base_resource': 'sriharsha',
        'entity': 'sriharsha.ent',
        'id': 'pyxbos-driver-hvac-1',
        'rate': 10,
        'building': building,
        'mortar_api_username': os.environ['MORTAR_API_USERNAME'],
        'mortar_api_password': os.environ['MORTAR_API_PASSWORD']
    }

    current_driver = BACnetDriver(cfg)
    current_driver.begin()