# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: iot.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import nullabletypes_pb2 as nullabletypes__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='iot.proto',
  package='xbospb',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\tiot.proto\x12\x06xbospb\x1a\x13nullabletypes.proto\"\'\n\x03URI\x12\x11\n\tnamespace\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t\"c\n\x06Triple\x12\x1c\n\x07subject\x18\x01 \x01(\x0b\x32\x0b.xbospb.URI\x12\x1e\n\tpredicate\x18\x02 \x01(\x0b\x32\x0b.xbospb.URI\x12\x1b\n\x06object\x18\x03 \x01(\x0b\x32\x0b.xbospb.URI\"\xc4\x01\n\x12XBOSIoTDeviceState\x12\x0c\n\x04time\x18\x01 \x01(\x04\x12\x11\n\trequestid\x18\x02 \x01(\x03\x12\r\n\x05\x65rror\x18\x03 \x01(\t\x12&\n\nthermostat\x18\x04 \x01(\x0b\x32\x12.xbospb.Thermostat\x12\x1c\n\x05meter\x18\x05 \x01(\x0b\x32\r.xbospb.Meter\x12\x1c\n\x05light\x18\x06 \x01(\x0b\x32\r.xbospb.Light\x12\x1a\n\x04\x65vse\x18\x07 \x01(\x0b\x32\x0c.xbospb.EVSE\"\xb9\x01\n\x16XBOSIoTDeviceActuation\x12\x0c\n\x04time\x18\x01 \x01(\x04\x12\x11\n\trequestid\x18\x02 \x01(\x03\x12&\n\nthermostat\x18\x03 \x01(\x0b\x32\x12.xbospb.Thermostat\x12\x1c\n\x05meter\x18\x04 \x01(\x0b\x32\r.xbospb.Meter\x12\x1c\n\x05light\x18\x05 \x01(\x0b\x32\r.xbospb.Light\x12\x1a\n\x04\x65vse\x18\x06 \x01(\x0b\x32\x0c.xbospb.EVSE\"?\n\x0eXBOSIoTContext\x12\x0c\n\x04time\x18\x01 \x01(\x04\x12\x1f\n\x07\x63ontext\x18\x02 \x03(\x0b\x32\x0e.xbospb.Triple\"\xda\x02\n\nThermostat\x12#\n\x0btemperature\x18\x01 \x01(\x0b\x32\x0e.xbospb.Double\x12)\n\x11relative_humidity\x18\x02 \x01(\x0b\x32\x0e.xbospb.Double\x12\x1e\n\x08override\x18\x03 \x01(\x0b\x32\x0c.xbospb.Bool\x12\x1f\n\tfan_state\x18\x04 \x01(\x0b\x32\x0c.xbospb.Bool\x12!\n\x08\x66\x61n_mode\x18\x05 \x01(\x0e\x32\x0f.xbospb.FanMode\x12\x1e\n\x04mode\x18\x06 \x01(\x0e\x32\x10.xbospb.HVACMode\x12 \n\x05state\x18\x07 \x01(\x0e\x32\x11.xbospb.HVACState\x12*\n\x13\x65nabled_heat_stages\x18\x08 \x01(\x0b\x32\r.xbospb.Int32\x12*\n\x13\x65nabled_cool_stages\x18\t \x01(\x0b\x32\r.xbospb.Int32\"o\n\x05Meter\x12\x1d\n\x05power\x18\x01 \x01(\x0b\x32\x0e.xbospb.Double\x12\x1f\n\x07voltage\x18\x02 \x01(\x0b\x32\x0e.xbospb.Double\x12&\n\x0e\x61pparent_power\x18\x03 \x01(\x0b\x32\x0e.xbospb.Double\"G\n\x05Light\x12\x1b\n\x05state\x18\x01 \x01(\x0b\x32\x0c.xbospb.Bool\x12!\n\nbrightness\x18\x02 \x01(\x0b\x32\r.xbospb.Int64\"\xb7\x01\n\x04\x45VSE\x12%\n\rcurrent_limit\x18\x01 \x01(\x0b\x32\x0e.xbospb.Double\x12\x1f\n\x07\x63urrent\x18\x02 \x01(\x0b\x32\x0e.xbospb.Double\x12\x1f\n\x07voltage\x18\x03 \x01(\x0b\x32\x0e.xbospb.Double\x12)\n\x12\x63harging_time_left\x18\x04 \x01(\x0b\x32\r.xbospb.Int32\x12\x1b\n\x05state\x18\x05 \x01(\x0b\x32\x0c.xbospb.Bool*-\n\x07\x46\x61nMode\x12\x0b\n\x07\x46\x61nAuto\x10\x00\x12\t\n\x05\x46\x61nOn\x10\x01\x12\n\n\x06\x46\x61nOff\x10\x02*Y\n\x08HVACMode\x12\x0f\n\x0bHVACModeOff\x10\x00\x12\x14\n\x10HVACModeHeatOnly\x10\x01\x12\x14\n\x10HVACModeCoolOnly\x10\x02\x12\x10\n\x0cHVACModeAuto\x10\x03*\x81\x01\n\tHVACState\x12\x10\n\x0cHVACStateOff\x10\x00\x12\x17\n\x13HVACStateHeatStage1\x10\x01\x12\x17\n\x13HVACStateCoolStage1\x10\x02\x12\x17\n\x13HVACStateHeatStage2\x10\x03\x12\x17\n\x13HVACStateCoolStage2\x10\x04\x62\x06proto3')
  ,
  dependencies=[nullabletypes__pb2.DESCRIPTOR,])

_FANMODE = _descriptor.EnumDescriptor(
  name='FanMode',
  full_name='xbospb.FanMode',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='FanAuto', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='FanOn', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='FanOff', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1357,
  serialized_end=1402,
)
_sym_db.RegisterEnumDescriptor(_FANMODE)

FanMode = enum_type_wrapper.EnumTypeWrapper(_FANMODE)
_HVACMODE = _descriptor.EnumDescriptor(
  name='HVACMode',
  full_name='xbospb.HVACMode',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='HVACModeOff', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACModeHeatOnly', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACModeCoolOnly', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACModeAuto', index=3, number=3,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1404,
  serialized_end=1493,
)
_sym_db.RegisterEnumDescriptor(_HVACMODE)

HVACMode = enum_type_wrapper.EnumTypeWrapper(_HVACMODE)
_HVACSTATE = _descriptor.EnumDescriptor(
  name='HVACState',
  full_name='xbospb.HVACState',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='HVACStateOff', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACStateHeatStage1', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACStateCoolStage1', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACStateHeatStage2', index=3, number=3,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='HVACStateCoolStage2', index=4, number=4,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1496,
  serialized_end=1625,
)
_sym_db.RegisterEnumDescriptor(_HVACSTATE)

HVACState = enum_type_wrapper.EnumTypeWrapper(_HVACSTATE)
FanAuto = 0
FanOn = 1
FanOff = 2
HVACModeOff = 0
HVACModeHeatOnly = 1
HVACModeCoolOnly = 2
HVACModeAuto = 3
HVACStateOff = 0
HVACStateHeatStage1 = 1
HVACStateCoolStage1 = 2
HVACStateHeatStage2 = 3
HVACStateCoolStage2 = 4



_URI = _descriptor.Descriptor(
  name='URI',
  full_name='xbospb.URI',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='namespace', full_name='xbospb.URI.namespace', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.URI.value', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=42,
  serialized_end=81,
)


_TRIPLE = _descriptor.Descriptor(
  name='Triple',
  full_name='xbospb.Triple',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='subject', full_name='xbospb.Triple.subject', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='predicate', full_name='xbospb.Triple.predicate', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='object', full_name='xbospb.Triple.object', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=83,
  serialized_end=182,
)


_XBOSIOTDEVICESTATE = _descriptor.Descriptor(
  name='XBOSIoTDeviceState',
  full_name='xbospb.XBOSIoTDeviceState',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='time', full_name='xbospb.XBOSIoTDeviceState.time', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='requestid', full_name='xbospb.XBOSIoTDeviceState.requestid', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error', full_name='xbospb.XBOSIoTDeviceState.error', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='thermostat', full_name='xbospb.XBOSIoTDeviceState.thermostat', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='meter', full_name='xbospb.XBOSIoTDeviceState.meter', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='light', full_name='xbospb.XBOSIoTDeviceState.light', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='evse', full_name='xbospb.XBOSIoTDeviceState.evse', index=6,
      number=7, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=185,
  serialized_end=381,
)


_XBOSIOTDEVICEACTUATION = _descriptor.Descriptor(
  name='XBOSIoTDeviceActuation',
  full_name='xbospb.XBOSIoTDeviceActuation',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='time', full_name='xbospb.XBOSIoTDeviceActuation.time', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='requestid', full_name='xbospb.XBOSIoTDeviceActuation.requestid', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='thermostat', full_name='xbospb.XBOSIoTDeviceActuation.thermostat', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='meter', full_name='xbospb.XBOSIoTDeviceActuation.meter', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='light', full_name='xbospb.XBOSIoTDeviceActuation.light', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='evse', full_name='xbospb.XBOSIoTDeviceActuation.evse', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=384,
  serialized_end=569,
)


_XBOSIOTCONTEXT = _descriptor.Descriptor(
  name='XBOSIoTContext',
  full_name='xbospb.XBOSIoTContext',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='time', full_name='xbospb.XBOSIoTContext.time', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='context', full_name='xbospb.XBOSIoTContext.context', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=571,
  serialized_end=634,
)


_THERMOSTAT = _descriptor.Descriptor(
  name='Thermostat',
  full_name='xbospb.Thermostat',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='temperature', full_name='xbospb.Thermostat.temperature', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='relative_humidity', full_name='xbospb.Thermostat.relative_humidity', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='override', full_name='xbospb.Thermostat.override', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='fan_state', full_name='xbospb.Thermostat.fan_state', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='fan_mode', full_name='xbospb.Thermostat.fan_mode', index=4,
      number=5, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='mode', full_name='xbospb.Thermostat.mode', index=5,
      number=6, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='state', full_name='xbospb.Thermostat.state', index=6,
      number=7, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='enabled_heat_stages', full_name='xbospb.Thermostat.enabled_heat_stages', index=7,
      number=8, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='enabled_cool_stages', full_name='xbospb.Thermostat.enabled_cool_stages', index=8,
      number=9, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=637,
  serialized_end=983,
)


_METER = _descriptor.Descriptor(
  name='Meter',
  full_name='xbospb.Meter',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='power', full_name='xbospb.Meter.power', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='voltage', full_name='xbospb.Meter.voltage', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='apparent_power', full_name='xbospb.Meter.apparent_power', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=985,
  serialized_end=1096,
)


_LIGHT = _descriptor.Descriptor(
  name='Light',
  full_name='xbospb.Light',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='state', full_name='xbospb.Light.state', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='brightness', full_name='xbospb.Light.brightness', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1098,
  serialized_end=1169,
)


_EVSE = _descriptor.Descriptor(
  name='EVSE',
  full_name='xbospb.EVSE',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='current_limit', full_name='xbospb.EVSE.current_limit', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='current', full_name='xbospb.EVSE.current', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='voltage', full_name='xbospb.EVSE.voltage', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='charging_time_left', full_name='xbospb.EVSE.charging_time_left', index=3,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='state', full_name='xbospb.EVSE.state', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1172,
  serialized_end=1355,
)

_TRIPLE.fields_by_name['subject'].message_type = _URI
_TRIPLE.fields_by_name['predicate'].message_type = _URI
_TRIPLE.fields_by_name['object'].message_type = _URI
_XBOSIOTDEVICESTATE.fields_by_name['thermostat'].message_type = _THERMOSTAT
_XBOSIOTDEVICESTATE.fields_by_name['meter'].message_type = _METER
_XBOSIOTDEVICESTATE.fields_by_name['light'].message_type = _LIGHT
_XBOSIOTDEVICESTATE.fields_by_name['evse'].message_type = _EVSE
_XBOSIOTDEVICEACTUATION.fields_by_name['thermostat'].message_type = _THERMOSTAT
_XBOSIOTDEVICEACTUATION.fields_by_name['meter'].message_type = _METER
_XBOSIOTDEVICEACTUATION.fields_by_name['light'].message_type = _LIGHT
_XBOSIOTDEVICEACTUATION.fields_by_name['evse'].message_type = _EVSE
_XBOSIOTCONTEXT.fields_by_name['context'].message_type = _TRIPLE
_THERMOSTAT.fields_by_name['temperature'].message_type = nullabletypes__pb2._DOUBLE
_THERMOSTAT.fields_by_name['relative_humidity'].message_type = nullabletypes__pb2._DOUBLE
_THERMOSTAT.fields_by_name['override'].message_type = nullabletypes__pb2._BOOL
_THERMOSTAT.fields_by_name['fan_state'].message_type = nullabletypes__pb2._BOOL
_THERMOSTAT.fields_by_name['fan_mode'].enum_type = _FANMODE
_THERMOSTAT.fields_by_name['mode'].enum_type = _HVACMODE
_THERMOSTAT.fields_by_name['state'].enum_type = _HVACSTATE
_THERMOSTAT.fields_by_name['enabled_heat_stages'].message_type = nullabletypes__pb2._INT32
_THERMOSTAT.fields_by_name['enabled_cool_stages'].message_type = nullabletypes__pb2._INT32
_METER.fields_by_name['power'].message_type = nullabletypes__pb2._DOUBLE
_METER.fields_by_name['voltage'].message_type = nullabletypes__pb2._DOUBLE
_METER.fields_by_name['apparent_power'].message_type = nullabletypes__pb2._DOUBLE
_LIGHT.fields_by_name['state'].message_type = nullabletypes__pb2._BOOL
_LIGHT.fields_by_name['brightness'].message_type = nullabletypes__pb2._INT64
_EVSE.fields_by_name['current_limit'].message_type = nullabletypes__pb2._DOUBLE
_EVSE.fields_by_name['current'].message_type = nullabletypes__pb2._DOUBLE
_EVSE.fields_by_name['voltage'].message_type = nullabletypes__pb2._DOUBLE
_EVSE.fields_by_name['charging_time_left'].message_type = nullabletypes__pb2._INT32
_EVSE.fields_by_name['state'].message_type = nullabletypes__pb2._BOOL
DESCRIPTOR.message_types_by_name['URI'] = _URI
DESCRIPTOR.message_types_by_name['Triple'] = _TRIPLE
DESCRIPTOR.message_types_by_name['XBOSIoTDeviceState'] = _XBOSIOTDEVICESTATE
DESCRIPTOR.message_types_by_name['XBOSIoTDeviceActuation'] = _XBOSIOTDEVICEACTUATION
DESCRIPTOR.message_types_by_name['XBOSIoTContext'] = _XBOSIOTCONTEXT
DESCRIPTOR.message_types_by_name['Thermostat'] = _THERMOSTAT
DESCRIPTOR.message_types_by_name['Meter'] = _METER
DESCRIPTOR.message_types_by_name['Light'] = _LIGHT
DESCRIPTOR.message_types_by_name['EVSE'] = _EVSE
DESCRIPTOR.enum_types_by_name['FanMode'] = _FANMODE
DESCRIPTOR.enum_types_by_name['HVACMode'] = _HVACMODE
DESCRIPTOR.enum_types_by_name['HVACState'] = _HVACSTATE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

URI = _reflection.GeneratedProtocolMessageType('URI', (_message.Message,), dict(
  DESCRIPTOR = _URI,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.URI)
  ))
_sym_db.RegisterMessage(URI)

Triple = _reflection.GeneratedProtocolMessageType('Triple', (_message.Message,), dict(
  DESCRIPTOR = _TRIPLE,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Triple)
  ))
_sym_db.RegisterMessage(Triple)

XBOSIoTDeviceState = _reflection.GeneratedProtocolMessageType('XBOSIoTDeviceState', (_message.Message,), dict(
  DESCRIPTOR = _XBOSIOTDEVICESTATE,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.XBOSIoTDeviceState)
  ))
_sym_db.RegisterMessage(XBOSIoTDeviceState)

XBOSIoTDeviceActuation = _reflection.GeneratedProtocolMessageType('XBOSIoTDeviceActuation', (_message.Message,), dict(
  DESCRIPTOR = _XBOSIOTDEVICEACTUATION,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.XBOSIoTDeviceActuation)
  ))
_sym_db.RegisterMessage(XBOSIoTDeviceActuation)

XBOSIoTContext = _reflection.GeneratedProtocolMessageType('XBOSIoTContext', (_message.Message,), dict(
  DESCRIPTOR = _XBOSIOTCONTEXT,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.XBOSIoTContext)
  ))
_sym_db.RegisterMessage(XBOSIoTContext)

Thermostat = _reflection.GeneratedProtocolMessageType('Thermostat', (_message.Message,), dict(
  DESCRIPTOR = _THERMOSTAT,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Thermostat)
  ))
_sym_db.RegisterMessage(Thermostat)

Meter = _reflection.GeneratedProtocolMessageType('Meter', (_message.Message,), dict(
  DESCRIPTOR = _METER,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Meter)
  ))
_sym_db.RegisterMessage(Meter)

Light = _reflection.GeneratedProtocolMessageType('Light', (_message.Message,), dict(
  DESCRIPTOR = _LIGHT,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Light)
  ))
_sym_db.RegisterMessage(Light)

EVSE = _reflection.GeneratedProtocolMessageType('EVSE', (_message.Message,), dict(
  DESCRIPTOR = _EVSE,
  __module__ = 'iot_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.EVSE)
  ))
_sym_db.RegisterMessage(EVSE)


# @@protoc_insertion_point(module_scope)
