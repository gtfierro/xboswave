# Quick script to autogenerate the hvac.proto file with definitions for unique equipment in the specified buildings
import pymortar
import subprocess
import os
import json
from collections import OrderedDict

def main():
    
    pymortar_client = pymortar.Client()
    buildings = ['orinda-public-library']
    proto_json = gen_proto_json() if os.path.exists('hvac.proto') and os.stat("hvac.proto").st_size != 0 else {}
    f = open("hvac.proto","w+")
    clustered_equipment = {"Damper": {}, "Fan": {}, "Thermostat": {}}
    
    f.write("syntax = \"proto3\";\r\n\n")
    f.write('package xbospb;\n\nimport "brick.proto";\nimport "nullabletypes.proto";\n\n')

    for building in buildings:
        # Go through all buildings, get equipment, and their respective points
        equipment = get_equipment_list(building, pymortar_client)

        #Add proto definitions for each piece of equipment
        for eq in equipment:
            points = dict.fromkeys(get_points_for_equipment(eq, building, pymortar_client), 0)

            # Dictates if the current equipment being processed is part of a bigger cluster or
            # is going to be its own message in hvac.proto
            clustered = False

            for clustered_eq in clustered_equipment:
                # Union all the clustered equipment
                if clustered_eq.lower() in eq.lower():
                    points = set(clustered_equipment[clustered_eq].keys()).union(points.keys())
                    clustered_equipment[clustered_eq] = dict.fromkeys(points, 0)
                    clustered = True
                    break

            if not clustered:
                # Write unclustered equipment to hvac.proto
                match_proto_numbering_and_write(f, eq, points, proto_json)

    # Write all clustered equipment to hvac.proto
    for eq in clustered_equipment:
        points = clustered_equipment[eq]
        if points:
            match_proto_numbering_and_write(f, eq, points, proto_json)
    
    f.close()

def match_proto_numbering_and_write(file, equipment, points, proto_json):
    """Matches the numbering of the existing points in hvac.proto and writes the points to the hvac.proto file"""

    if equipment in proto_json:
        # If the equipment already exists in hvac.proto
        new_points = []
        for p in points:
            if p in proto_json[equipment]["Fields"]:
                points[p] = proto_json[equipment]["Fields"][p]["Number"]
            else:
                new_points.append(p)

        if new_points:    
            for i in range(1, len(points.values()) + 1):
                if new_points and i not in points.values():
                    points[new_points[0]] = i
                    new_points.pop(0)

        points = OrderedDict(sorted(points.items(), key=lambda kv: kv[1]))

        write_proto_message(file, points, equipment)
    else:
        # If the equipment is new in hvac.proto
        write_proto_message(file, dict(zip(points.keys(), list(range(1, len(points) + 1)))), equipment)

def gen_proto_json():
    """Generates json dictionary out of the hvac.proto file using protodef_to_json"""

    # Run the protodef_to_json on hvac.proto to generate a json representation
    cur_dir = os.getcwd()
    os.chdir("../tools/protodef_to_json")
    subprocess.call(["go", "build"])
    proto_str = subprocess.Popen(["./protodef_to_json", "../../proto/hvac.proto"], stdout=subprocess.PIPE).communicate()[0].decode("utf-8")
    proto_str = proto_str.replace("'", "")
    os.chdir(cur_dir)
    pjson = json.loads(proto_str)

    # Processing the json object to make it more easy to use
    keys, values = [], []
    for record in pjson:
        keys.append(record["Class"])
        k, v = [], []
        for field in record["Fields"]:
            k.append(field["Class"])
            v.append(field)

        record["Fields"] = dict(zip(k, v))
        values.append(record)

    return dict(zip(keys, values))

def get_equipment_list(building, pymortar_client):
    """Gets unique list of equipment in the building"""

    #Get unique list of equipment
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

    # Excludes CentralPlant for now
    equip_query = res.query("SELECT DISTINCT SUBSTR(equipclass, INSTR(equipclass, '#') + 1) AS eq FROM equipment WHERE equipname NOT LIKE '%Central%'")

    return [row[0] for row in equip_query]

def get_points_for_equipment(equipment, building, pymortar_client):
    """Gets points for each equipment in the list of equipment"""

    v = pymortar.View(
        name=equipment,
        definition="""
        SELECT ?{eq_lower} ?point ?class FROM {building} WHERE {{
            ?{eq_lower} rdf:type brick:{eq_class} .
            ?{eq_lower} bf:hasPoint ?point .
            ?point rdf:type ?class
        }};""".format(eq_lower=equipment.lower(), eq_class=equipment, building=building) 
    )

    res = pymortar_client.fetch(pymortar.FetchRequest(
        sites=[building],
        views=[v]
    ))

    return res.view(equipment)['class'].unique()

def write_proto_message(file, points, equipment):
    """Writes a single message struct in the specified protobuf file"""

    file.write("message %s {\r\n" % equipment)
    file.write("\toption (brick_equip_class).namespace = 'https://brickschema.org/schema/1.0.3/Brick#';\r\n")
    file.write("\toption (brick_equip_class).value = '%s';\r\n\n" % equipment)

    for p in points:
        point_type = "Double"
        file.write("""\t{point_type} {point_lower} = {proto_index} [(brick_point_class).namespace='https://brickschema.org/schema/1.0.3/Brick#',(brick_point_class).value='{point}'];\r\n""".format(point_type=point_type, point_lower=p.lower(), proto_index=points[p], point=p))

    file.write("}\r\n\n")


if __name__ == '__main__':
    main()