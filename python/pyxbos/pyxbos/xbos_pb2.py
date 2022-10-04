# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: xbos.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import hamilton_pb2 as hamilton__pb2
from . import iot_pb2 as iot__pb2
from . import dentmeter_pb2 as dentmeter__pb2
from . import system_monitor_pb2 as system__monitor__pb2
from . import parker_pb2 as parker__pb2
from . import wattnode_pb2 as wattnode__pb2
from . import flexstat_pb2 as flexstat__pb2
from . import c37_pb2 as c37__pb2
from . import energise_pb2 as energise__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='xbos.proto',
  package='xbospb',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\nxbos.proto\x12\x06xbospb\x1a\x0ehamilton.proto\x1a\tiot.proto\x1a\x0f\x64\x65ntmeter.proto\x1a\x14system_monitor.proto\x1a\x0cparker.proto\x1a\x0ewattnode.proto\x1a\x0e\x66lexstat.proto\x1a\tc37.proto\x1a\x0e\x65nergise.proto\"\xe5\x05\n\x04XBOS\x12*\n\x0cHamiltonData\x18\x32 \x01(\x0b\x32\x14.xbospb.HamiltonData\x12\x38\n\x13HamiltonBRLinkStats\x18\x33 \x01(\x0b\x32\x1b.xbospb.HamiltonBRLinkStats\x12\x34\n\x11HamiltonBRMessage\x18\x34 \x01(\x0b\x32\x19.xbospb.HamiltonBRMessage\x12\x36\n\x12XBOSIoTDeviceState\x18\x64 \x01(\x0b\x32\x1a.xbospb.XBOSIoTDeviceState\x12>\n\x16XBOSIoTDeviceActuation\x18\x65 \x01(\x0b\x32\x1e.xbospb.XBOSIoTDeviceActuation\x12.\n\x0eXBOSIoTContext\x18\x66 \x01(\x0b\x32\x16.xbospb.XBOSIoTContext\x12/\n\x0e\x44\x65ntMeterState\x18\x96\x01 \x01(\x0b\x32\x16.xbospb.DentMeterState\x12*\n\x0cparker_state\x18\x97\x01 \x01(\x0b\x32\x13.xbospb.ParkerState\x12.\n\x0ewattnode_state\x18\x98\x01 \x01(\x0b\x32\x15.xbospb.WattnodeState\x12.\n\x0e\x66lexstat_state\x18\x99\x01 \x01(\x0b\x32\x15.xbospb.FlexstatState\x12\x45\n\x1a\x66lexstat_actuation_message\x18\x9a\x01 \x01(\x0b\x32 .xbospb.FlexstatActuationMessage\x12\x35\n\x11\x42\x61sicServerStatus\x18\xc8\x01 \x01(\x0b\x32\x19.xbospb.BasicServerStatus\x12+\n\x0c\x43\x33\x37\x44\x61taFrame\x18\xfa\x01 \x01(\x0b\x32\x14.xbospb.C37DataFrame\x12\x31\n\x0f\x45nergiseMessage\x18\xfb\x01 \x01(\x0b\x32\x17.xbospb.EnergiseMessage\"x\n\x08Resource\x12$\n\ttransport\x18\x01 \x01(\x0e\x32\x11.xbospb.Transport\x12\x11\n\tnamespace\x18\x02 \x01(\t\x12\x0f\n\x07service\x18\x03 \x01(\t\x12\x10\n\x08instance\x18\x04 \x01(\t\x12\x10\n\x08location\x18\x05 \x01(\t*!\n\tTransport\x12\n\n\x06WAVEMQ\x10\x00\x12\x08\n\x04GRPC\x10\x01\x62\x06proto3')
  ,
  dependencies=[hamilton__pb2.DESCRIPTOR,iot__pb2.DESCRIPTOR,dentmeter__pb2.DESCRIPTOR,system__monitor__pb2.DESCRIPTOR,parker__pb2.DESCRIPTOR,wattnode__pb2.DESCRIPTOR,flexstat__pb2.DESCRIPTOR,c37__pb2.DESCRIPTOR,energise__pb2.DESCRIPTOR,])

_TRANSPORT = _descriptor.EnumDescriptor(
  name='Transport',
  full_name='xbospb.Transport',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='WAVEMQ', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GRPC', index=1, number=1,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1027,
  serialized_end=1060,
)
_sym_db.RegisterEnumDescriptor(_TRANSPORT)

Transport = enum_type_wrapper.EnumTypeWrapper(_TRANSPORT)
WAVEMQ = 0
GRPC = 1



_XBOS = _descriptor.Descriptor(
  name='XBOS',
  full_name='xbospb.XBOS',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='HamiltonData', full_name='xbospb.XBOS.HamiltonData', index=0,
      number=50, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='HamiltonBRLinkStats', full_name='xbospb.XBOS.HamiltonBRLinkStats', index=1,
      number=51, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='HamiltonBRMessage', full_name='xbospb.XBOS.HamiltonBRMessage', index=2,
      number=52, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='XBOSIoTDeviceState', full_name='xbospb.XBOS.XBOSIoTDeviceState', index=3,
      number=100, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='XBOSIoTDeviceActuation', full_name='xbospb.XBOS.XBOSIoTDeviceActuation', index=4,
      number=101, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='XBOSIoTContext', full_name='xbospb.XBOS.XBOSIoTContext', index=5,
      number=102, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='DentMeterState', full_name='xbospb.XBOS.DentMeterState', index=6,
      number=150, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='parker_state', full_name='xbospb.XBOS.parker_state', index=7,
      number=151, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='wattnode_state', full_name='xbospb.XBOS.wattnode_state', index=8,
      number=152, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='flexstat_state', full_name='xbospb.XBOS.flexstat_state', index=9,
      number=153, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='flexstat_actuation_message', full_name='xbospb.XBOS.flexstat_actuation_message', index=10,
      number=154, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='BasicServerStatus', full_name='xbospb.XBOS.BasicServerStatus', index=11,
      number=200, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='C37DataFrame', full_name='xbospb.XBOS.C37DataFrame', index=12,
      number=250, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='EnergiseMessage', full_name='xbospb.XBOS.EnergiseMessage', index=13,
      number=251, type=11, cpp_type=10, label=1,
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
  serialized_start=162,
  serialized_end=903,
)


_RESOURCE = _descriptor.Descriptor(
  name='Resource',
  full_name='xbospb.Resource',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='transport', full_name='xbospb.Resource.transport', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='namespace', full_name='xbospb.Resource.namespace', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='service', full_name='xbospb.Resource.service', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='instance', full_name='xbospb.Resource.instance', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='location', full_name='xbospb.Resource.location', index=4,
      number=5, type=9, cpp_type=9, label=1,
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
  serialized_start=905,
  serialized_end=1025,
)

_XBOS.fields_by_name['HamiltonData'].message_type = hamilton__pb2._HAMILTONDATA
_XBOS.fields_by_name['HamiltonBRLinkStats'].message_type = hamilton__pb2._HAMILTONBRLINKSTATS
_XBOS.fields_by_name['HamiltonBRMessage'].message_type = hamilton__pb2._HAMILTONBRMESSAGE
_XBOS.fields_by_name['XBOSIoTDeviceState'].message_type = iot__pb2._XBOSIOTDEVICESTATE
_XBOS.fields_by_name['XBOSIoTDeviceActuation'].message_type = iot__pb2._XBOSIOTDEVICEACTUATION
_XBOS.fields_by_name['XBOSIoTContext'].message_type = iot__pb2._XBOSIOTCONTEXT
_XBOS.fields_by_name['DentMeterState'].message_type = dentmeter__pb2._DENTMETERSTATE
_XBOS.fields_by_name['parker_state'].message_type = parker__pb2._PARKERSTATE
_XBOS.fields_by_name['wattnode_state'].message_type = wattnode__pb2._WATTNODESTATE
_XBOS.fields_by_name['flexstat_state'].message_type = flexstat__pb2._FLEXSTATSTATE
_XBOS.fields_by_name['flexstat_actuation_message'].message_type = flexstat__pb2._FLEXSTATACTUATIONMESSAGE
_XBOS.fields_by_name['BasicServerStatus'].message_type = system__monitor__pb2._BASICSERVERSTATUS
_XBOS.fields_by_name['C37DataFrame'].message_type = c37__pb2._C37DATAFRAME
_XBOS.fields_by_name['EnergiseMessage'].message_type = energise__pb2._ENERGISEMESSAGE
_RESOURCE.fields_by_name['transport'].enum_type = _TRANSPORT
DESCRIPTOR.message_types_by_name['XBOS'] = _XBOS
DESCRIPTOR.message_types_by_name['Resource'] = _RESOURCE
DESCRIPTOR.enum_types_by_name['Transport'] = _TRANSPORT
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

XBOS = _reflection.GeneratedProtocolMessageType('XBOS', (_message.Message,), dict(
  DESCRIPTOR = _XBOS,
  __module__ = 'xbos_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.XBOS)
  ))
_sym_db.RegisterMessage(XBOS)

Resource = _reflection.GeneratedProtocolMessageType('Resource', (_message.Message,), dict(
  DESCRIPTOR = _RESOURCE,
  __module__ = 'xbos_pb2'
  # @@protoc_insertion_point(class_scope:xbospb.Resource)
  ))
_sym_db.RegisterMessage(Resource)


# @@protoc_insertion_point(module_scope)
