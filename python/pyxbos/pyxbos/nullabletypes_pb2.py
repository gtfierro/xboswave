# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: nullabletypes.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='nullabletypes.proto',
  package='xbospb',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x13nullabletypes.proto\x12\x06xbospb\"\x16\n\x05Int32\x12\r\n\x05value\x18\x01 \x01(\x05\"\x16\n\x05Int64\x12\r\n\x05value\x18\x01 \x01(\x03\"\x17\n\x06Uint64\x12\r\n\x05value\x18\x01 \x01(\x04\"\x17\n\x06\x44ouble\x12\r\n\x05value\x18\x01 \x01(\x01\"\x15\n\x04\x42ool\x12\r\n\x05value\x18\x01 \x01(\x08\".\n\x07\x46\x61nMode\x12#\n\x05value\x18\x01 \x01(\x0e\x32\x14.xbospb.FanModeValue\"0\n\x08HVACMode\x12$\n\x05value\x18\x01 \x01(\x0e\x32\x15.xbospb.HVACModeValue\"2\n\tHVACState\x12%\n\x05value\x18\x01 \x01(\x0e\x32\x16.xbospb.HVACStateValue*2\n\x0c\x46\x61nModeValue\x12\x0b\n\x07\x46\x61nAuto\x10\x00\x12\t\n\x05\x46\x61nOn\x10\x01\x12\n\n\x06\x46\x61nOff\x10\x02*^\n\rHVACModeValue\x12\x0f\n\x0bHVACModeOff\x10\x00\x12\x14\n\x10HVACModeHeatOnly\x10\x01\x12\x14\n\x10HVACModeCoolOnly\x10\x02\x12\x10\n\x0cHVACModeAuto\x10\x03*\x86\x01\n\x0eHVACStateValue\x12\x10\n\x0cHVACStateOff\x10\x00\x12\x17\n\x13HVACStateHeatStage1\x10\x01\x12\x17\n\x13HVACStateCoolStage1\x10\x02\x12\x17\n\x13HVACStateHeatStage2\x10\x03\x12\x17\n\x13HVACStateCoolStage2\x10\x04\x62\x06proto3')
)

_FANMODEVALUE = _descriptor.EnumDescriptor(
  name='FanModeValue',
  full_name='xbospb.FanModeValue',
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
  serialized_start=302,
  serialized_end=352,
)
_sym_db.RegisterEnumDescriptor(_FANMODEVALUE)

FanModeValue = enum_type_wrapper.EnumTypeWrapper(_FANMODEVALUE)
_HVACMODEVALUE = _descriptor.EnumDescriptor(
  name='HVACModeValue',
  full_name='xbospb.HVACModeValue',
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
  serialized_start=354,
  serialized_end=448,
)
_sym_db.RegisterEnumDescriptor(_HVACMODEVALUE)

HVACModeValue = enum_type_wrapper.EnumTypeWrapper(_HVACMODEVALUE)
_HVACSTATEVALUE = _descriptor.EnumDescriptor(
  name='HVACStateValue',
  full_name='xbospb.HVACStateValue',
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
  serialized_start=451,
  serialized_end=585,
)
_sym_db.RegisterEnumDescriptor(_HVACSTATEVALUE)

HVACStateValue = enum_type_wrapper.EnumTypeWrapper(_HVACSTATEVALUE)
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



_INT32 = _descriptor.Descriptor(
  name='Int32',
  full_name='xbospb.Int32',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.Int32.value', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=31,
  serialized_end=53,
)


_INT64 = _descriptor.Descriptor(
  name='Int64',
  full_name='xbospb.Int64',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.Int64.value', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=55,
  serialized_end=77,
)


_UINT64 = _descriptor.Descriptor(
  name='Uint64',
  full_name='xbospb.Uint64',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.Uint64.value', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=79,
  serialized_end=102,
)


_DOUBLE = _descriptor.Descriptor(
  name='Double',
  full_name='xbospb.Double',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.Double.value', index=0,
      number=1, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
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
  serialized_start=104,
  serialized_end=127,
)


_BOOL = _descriptor.Descriptor(
  name='Bool',
  full_name='xbospb.Bool',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.Bool.value', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
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
  serialized_start=129,
  serialized_end=150,
)


_FANMODE = _descriptor.Descriptor(
  name='FanMode',
  full_name='xbospb.FanMode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.FanMode.value', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=152,
  serialized_end=198,
)


_HVACMODE = _descriptor.Descriptor(
  name='HVACMode',
  full_name='xbospb.HVACMode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.HVACMode.value', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=200,
  serialized_end=248,
)


_HVACSTATE = _descriptor.Descriptor(
  name='HVACState',
  full_name='xbospb.HVACState',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='value', full_name='xbospb.HVACState.value', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=250,
  serialized_end=300,
)

_FANMODE.fields_by_name['value'].enum_type = _FANMODEVALUE
_HVACMODE.fields_by_name['value'].enum_type = _HVACMODEVALUE
_HVACSTATE.fields_by_name['value'].enum_type = _HVACSTATEVALUE
DESCRIPTOR.message_types_by_name['Int32'] = _INT32
DESCRIPTOR.message_types_by_name['Int64'] = _INT64
DESCRIPTOR.message_types_by_name['Uint64'] = _UINT64
DESCRIPTOR.message_types_by_name['Double'] = _DOUBLE
DESCRIPTOR.message_types_by_name['Bool'] = _BOOL
DESCRIPTOR.message_types_by_name['FanMode'] = _FANMODE
DESCRIPTOR.message_types_by_name['HVACMode'] = _HVACMODE
DESCRIPTOR.message_types_by_name['HVACState'] = _HVACSTATE
DESCRIPTOR.enum_types_by_name['FanModeValue'] = _FANMODEVALUE
DESCRIPTOR.enum_types_by_name['HVACModeValue'] = _HVACMODEVALUE
DESCRIPTOR.enum_types_by_name['HVACStateValue'] = _HVACSTATEVALUE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Int32 = _reflection.GeneratedProtocolMessageType('Int32', (_message.Message,), dict(
  DESCRIPTOR = _INT32,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Int32)
  ))
_sym_db.RegisterMessage(Int32)

Int64 = _reflection.GeneratedProtocolMessageType('Int64', (_message.Message,), dict(
  DESCRIPTOR = _INT64,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Int64)
  ))
_sym_db.RegisterMessage(Int64)

Uint64 = _reflection.GeneratedProtocolMessageType('Uint64', (_message.Message,), dict(
  DESCRIPTOR = _UINT64,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Uint64)
  ))
_sym_db.RegisterMessage(Uint64)

Double = _reflection.GeneratedProtocolMessageType('Double', (_message.Message,), dict(
  DESCRIPTOR = _DOUBLE,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Double)
  ))
_sym_db.RegisterMessage(Double)

Bool = _reflection.GeneratedProtocolMessageType('Bool', (_message.Message,), dict(
  DESCRIPTOR = _BOOL,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Bool)
  ))
_sym_db.RegisterMessage(Bool)

FanMode = _reflection.GeneratedProtocolMessageType('FanMode', (_message.Message,), dict(
  DESCRIPTOR = _FANMODE,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.FanMode)
  ))
_sym_db.RegisterMessage(FanMode)

HVACMode = _reflection.GeneratedProtocolMessageType('HVACMode', (_message.Message,), dict(
  DESCRIPTOR = _HVACMODE,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.HVACMode)
  ))
_sym_db.RegisterMessage(HVACMode)

HVACState = _reflection.GeneratedProtocolMessageType('HVACState', (_message.Message,), dict(
  DESCRIPTOR = _HVACSTATE,
  __module__ = 'nullabletypes_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.HVACState)
  ))
_sym_db.RegisterMessage(HVACState)


# @@protoc_insertion_point(module_scope)
