# Protodef to JSON conversion

Creates a JSON representation of messages, fields, field numbers and Brick classes + namespaces from XBOS `proto` definition files.
This is typically easier to work with in other languages that don't have an easily accessible (or completE) protobuf definition parser (many thanks to https://github.com/yoheimuta/go-protoparser)

## Usage

```
$ go build
$ ./protodef_to_json path/to/protofile.proto
```

## Example

```
$ ./protodef_to_json ../../proto/hvac.proto
```

produces

```json
[
  {
    "Name": "VAV",
    "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
    "Class": "'VAV'",
    "Fields": [
      {
        "Name": "discharge_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Temperature_Sensor'",
        "Number": 1
      },
      {
        "Name": "thermostat_adjust_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Thermostat_Adjust_Setpoint'",
        "Number": 2
      },
      {
        "Name": "zone_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Temperature_Sensor'",
        "Number": 3
      },
      {
        "Name": "cooling_max_supply_air_flow_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Cooling_Max_Supply_Air_Flow_Setpoint'",
        "Number": 4
      },
      {
        "Name": "supply_air_flow_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Flow_Sensor'",
        "Number": 5
      },
      {
        "Name": "zone_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Temperature_Setpoint'",
        "Number": 6
      },
      {
        "Name": "box_mode",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Box_Mode'",
        "Number": 7
      },
      {
        "Name": "cooling_demand",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Cooling_Demand'",
        "Number": 8
      },
      {
        "Name": "occupied_heating_min_supply_air_flow_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupied_Heating_Min_Supply_Air_Flow_Setpoint'",
        "Number": 9
      },
      {
        "Name": "supply_air_flow_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Flow_Setpoint'",
        "Number": 10
      },
      {
        "Name": "supply_air_velocity_pressure_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Velocity_Pressure_Sensor'",
        "Number": 11
      },
      {
        "Name": "heating_demand",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Heating_Demand'",
        "Number": 12
      }
    ]
  },
  {
    "Name": "VVT",
    "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
    "Class": "'VVT'",
    "Fields": [
      {
        "Name": "occupied_dead_band",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupied_Dead_Band'",
        "Number": 1
      },
      {
        "Name": "zone_demand",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Demand'",
        "Number": 2
      },
      {
        "Name": "zone_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Temperature_Sensor'",
        "Number": 3
      },
      {
        "Name": "supply_air_damper_min_position_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Damper_Min_Position_Setpoint'",
        "Number": 4
      },
      {
        "Name": "thermostat_adjust_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Thermostat_Adjust_Setpoint'",
        "Number": 5
      },
      {
        "Name": "occupied_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupied_Temperature_Setpoint'",
        "Number": 6
      },
      {
        "Name": "occupied_mode_status",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupied_Mode_Status'",
        "Number": 7
      },
      {
        "Name": "supply_air_damper_max_position_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Damper_Max_Position_Setpoint'",
        "Number": 8
      },
      {
        "Name": "unoccupied_dead_band",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Unoccupied_Dead_Band'",
        "Number": 9
      },
      {
        "Name": "discharge_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Temperature_Sensor'",
        "Number": 10
      },
      {
        "Name": "occupancy_command",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupancy_Command'",
        "Number": 11
      },
      {
        "Name": "unoccupied_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Unoccupied_Temperature_Setpoint'",
        "Number": 12
      }
    ]
  },
  {
    "Name": "AHU",
    "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
    "Class": "'AHU'",
    "Fields": [
      {
        "Name": "zone_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Temperature_Setpoint'",
        "Number": 1
      },
      {
        "Name": "outside_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Outside_Air_Temperature_Sensor'",
        "Number": 2
      },
      {
        "Name": "mixed_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Mixed_Air_Temperature_Sensor'",
        "Number": 3
      },
      {
        "Name": "occupied_mode_status",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupied_Mode_Status'",
        "Number": 4
      },
      {
        "Name": "filter_status",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Filter_Status'",
        "Number": 5
      },
      {
        "Name": "shutdown_command",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Shutdown_Command'",
        "Number": 6
      },
      {
        "Name": "mixed_air_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Mixed_Air_Temperature_Setpoint'",
        "Number": 7
      },
      {
        "Name": "supply_air_damper_min_position_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Supply_Air_Damper_Min_Position_Setpoint'",
        "Number": 8
      },
      {
        "Name": "mixed_air_temperature_low_limit_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Mixed_Air_Temperature_Low_Limit_Setpoint'",
        "Number": 9
      },
      {
        "Name": "zone_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Zone_Temperature_Sensor'",
        "Number": 10
      },
      {
        "Name": "return_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Return_Air_Temperature_Sensor'",
        "Number": 11
      },
      {
        "Name": "heating_valve_command",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Heating_Valve_Command'",
        "Number": 12
      },
      {
        "Name": "discharge_air_temperature_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Temperature_Sensor'",
        "Number": 13
      },
      {
        "Name": "discharge_air_temperature_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Temperature_Setpoint'",
        "Number": 14
      },
      {
        "Name": "building_static_pressure_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Building_Static_Pressure_Sensor'",
        "Number": 15
      },
      {
        "Name": "cooling_valve_command",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Cooling_Valve_Command'",
        "Number": 16
      },
      {
        "Name": "occupancy_command",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Occupancy_Command'",
        "Number": 17
      },
      {
        "Name": "discharge_air_static_pressure_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Static_Pressure_Setpoint'",
        "Number": 18
      },
      {
        "Name": "discharge_air_static_pressure_sensor",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Discharge_Air_Static_Pressure_Sensor'",
        "Number": 19
      },
      {
        "Name": "cooling_demand",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Cooling_Demand'",
        "Number": 20
      },
      {
        "Name": "building_static_pressure_setpoint",
        "Namespace": "'https://brickschema.org/schema/1.0.3/Brick#'",
        "Class": "'Building_Static_Pressure_Setpoint'",
        "Number": 21
      }
    ]
  }
]
```
