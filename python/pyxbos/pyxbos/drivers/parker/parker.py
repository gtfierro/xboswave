from pyxbos import *
from modbus_driver import Modbus_Driver
import os,sys
import json
import requests
import yaml
import argparse
import time
from inspect import getmembers


class ParkerDriver(Driver):
    def setup(self, cfg):
        with open('parker.yaml') as f:
            # use safe_load instead load for security reasons
            driverConfig = yaml.safe_load(f)
        self.modbus_device = Modbus_Driver(driverConfig)
        self.modbus_device.initialize_modbus()
        self.service_name = cfg['service_name']

    def read(self, requestid=None):
        output = self.modbus_device.get_data()

        # This is necessary because we need to bitwise operations to unpack the
        # the flags received from the modbus register
        regulator_flag_1 = output['regulator_flag_1']
        output['energy_saving_regulator_flag'] = bool(regulator_flag_1 & 0x0100)
        output['energy_saving_real_time_regulator_flag'] = bool(regulator_flag_1 & 0x0200)
        output['service_request_regulator_flag'] = bool(regulator_flag_1 & 0x0400)

        regulator_flag_2 = output['regulator_flag_2']
        output['on_standby_regulator_flag'] = bool(regulator_flag_2 & 0x0001)
        output['new_alarm_to_read_regulator_flag'] = bool(regulator_flag_2 & 0x0080)
        output['defrost_status_regulator_flag']	= bool(regulator_flag_2 & 0x0700)

        digital_io_status = output['digital_io_status']
        output['door_switch_input_status'] = bool(digital_io_status & 0x0001)
        output['multipurpose_input_status'] = bool(digital_io_status & 0x0002)
        output['compressor_status'] = bool(digital_io_status & 0x0100)
        output['output_defrost_status'] = bool(digital_io_status & 0x0200)
        output['fans_status'] = bool(digital_io_status & 0x0400)
        output['output_k4_status'] = bool(digital_io_status & 0x0800)

        digital_output_flags = output['digital_output_flags']
        output['energy_saving_status'] = bool(digital_output_flags & 0x0100)
        output['service_request_status'] =	bool(digital_output_flags & 0x0200)
        output['resistors_activated_by_aux_key_status'] = bool(digital_output_flags & 0x001)
        output['evaporator_valve_state'] = bool(digital_output_flags & 0x002)
        output['output_defrost_state'] = bool(digital_output_flags & 0x004)
        output['output_lux_state'] =	bool(digital_output_flags & 0x008)
        output['output_aux_state'] =	bool(digital_output_flags & 0x0010)
        output['resistors_state'] = bool(digital_output_flags & 0x0020)
        output['output_alarm_state'] =	bool(digital_output_flags & 0x0040)
        output['second_compressor_state'] =	bool(digital_output_flags & 0x0080)

        alarm_status = output['alarm_status']
        #print(format(output['alarm_status'], '#010b'))
        output['probe1_failure_alarm'] = bool(alarm_status & 0x0100)
        output['probe2_failure_alarm'] = bool(alarm_status & 0x0200)
        output['probe3_failure_alarm'] = bool(alarm_status & 0x0400)
        output['minimum_temperature_alarm'] = bool(alarm_status & 0x1000)
        output['maximum_temperture_alarm'] = bool(alarm_status & 0x2000)
        output['condensor_temperature_failure_alarm'] = bool(alarm_status & 0x4000)
        output['condensor_pre_alarm'] = bool(alarm_status & 0x8000)
        output['door_alarm'] = bool(alarm_status & 0x0004)
        output['multipurpose_input_alarm'] = bool(alarm_status & 0x0008)
        output['compressor_blocked_alarm'] = bool(alarm_status & 0x0010)
        output['power_failure_alarm'] = bool(alarm_status & 0x0020)
        output['rtc_error_alarm'] = bool(alarm_status & 0x0080)
        #print(output['rtc_error_alarm'])
        #print(format(output['rtc_error_alarm'], '#010b'))
        print(output)

        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                parker_state = parker_pb2.ParkerState(
                    compressor_working_hours  =   types.Double(value=output.get('compressor_working_hours',None)),
                    on_standby_status  =   types.Int64(value=output.get('on_standby_status',None)),
                    light_status  =   types.Int64(value=output.get('light_status',None)),
                    aux_output_status  =   types.Int64(value=output.get('aux_output_status',None)),
                    next_defrost_counter  =   types.Double(value=output.get('next_defrost_counter',None)),
                    door_switch_input_status  =   types.Int64(value=output.get('door_switch_input_status',None)),
                    multipurpose_input_status  =   types.Int64(value=output.get('multipurpose_input_status',None)),
                    compressor_status  =   types.Int64(value=output.get('compressor_status',None)),
                    output_defrost_status  =   types.Int64(value=output.get('output_defrost_status',None)),
                    fans_status  =   types.Int64(value=output.get('fans_status',None)),
                    output_k4_status  =   types.Int64(value=output.get('output_k4_status',None)),
                    cabinet_temperature  =   types.Double(value=output.get('cabinet_temperature',None)),
                    evaporator_temperature  =   types.Double(value=output.get('evaporator_temperature',None)),
                    auxiliary_temperature  =   types.Double(value=output.get('auxiliary_temperature',None)),
                    probe1_failure_alarm  =   types.Int64(value=output.get('probe1_failure_alarm',None)),
                    probe2_failure_alarm  =   types.Int64(value=output.get('probe2_failure_alarm',None)),
                    probe3_failure_alarm  =   types.Int64(value=output.get('probe3_failure_alarm',None)),
                    minimum_temperature_alarm  =   types.Int64(value=output.get('minimum_temperature_alarm',None)),
                    maximum_temperture_alarm  =   types.Int64(value=output.get('maximum_temperture_alarm',None)),
                    condensor_temperature_failure_alarm  =   types.Int64(value=output.get('condensor_temperature_failure_alarm',None)),
                    condensor_pre_alarm  =   types.Int64(value=output.get('condensor_pre_alarm',None)),
                    door_alarm  =   types.Int64(value=output.get('door_alarm',None)),
                    multipurpose_input_alarm  =   types.Int64(value=output.get('multipurpose_input_alarm',None)),
                    compressor_blocked_alarm  =   types.Int64(value=output.get('compressor_blocked_alarm',None)),
                    power_failure_alarm  =   types.Int64(value=output.get('power_failure_alarm',None)),
                    rtc_error_alarm  =   types.Int64(value=output.get('rtc_error_alarm',None)),
                    energy_saving_regulator_flag  =   types.Int64(value=output.get('energy_saving_regulator_flag',None)),
                    energy_saving_real_time_regulator_flag  =   types.Int64(value=output.get('energy_saving_real_time_regulator_flag',None)),
                    service_request_regulator_flag  =   types.Int64(value=output.get('service_request_regulator_flag',None)),
                    on_standby_regulator_flag  =   types.Int64(value=output.get('on_standby_regulator_flag',None)),
                    new_alarm_to_read_regulator_flag  =   types.Int64(value=output.get('new_alarm_to_read_regulator_flag',None)),
                    defrost_status_regulator_flag  =   types.Int64(value=output.get('defrost_status_regulator_flag',None)),
                    active_setpoint  =   types.Int64(value=output.get('active_setpoint',None)),
                    time_until_defrost  =   types.Int64(value=output.get('time_until_defrost',None)),
                    current_defrost_counter  =   types.Int64(value=output.get('current_defrost_counter',None)),
                    compressor_delay  =   types.Int64(value=output.get('compressor_delay',None)),
                    num_alarms_in_history  =   types.Int64(value=output.get('num_alarms_in_history',None)),
                    energy_saving_status  =   types.Int64(value=output.get('energy_saving_status',None)),
                    service_request_status  =   types.Int64(value=output.get('service_request_status',None)),
                    resistors_activated_by_aux_key_status  =   types.Int64(value=output.get('resistors_activated_by_aux_key_status',None)),
                    evaporator_valve_state  =   types.Int64(value=output.get('evaporator_valve_state',None)),
                    output_defrost_state  =   types.Int64(value=output.get('output_defrost_state',None)),
                    output_lux_state  =   types.Int64(value=output.get('output_lux_state',None)),
                    output_aux_state  =   types.Int64(value=output.get('output_aux_state',None)),
                    resistors_state  =   types.Int64(value=output.get('resistors_state',None)),
                    output_alarm_state  =   types.Int64(value=output.get('output_alarm_state',None)),
                    second_compressor_state  =   types.Int64(value=output.get('second_compressor_state',None)),
                    setpoint  =   types.Double(value=output.get('setpoint',None)),
                    r1  =   types.Double(value=output.get('r1',None)),
                    r2  =   types.Double(value=output.get('r2',None)),
                    r4  =   types.Double(value=output.get('r4',None)),
                    C0  =   types.Double(value=output.get('C0',None)),
                    C1  =   types.Double(value=output.get('C1',None)),
                    d0  =   types.Double(value=output.get('d0',None)),
                    d3  =   types.Double(value=output.get('d3',None)),
                    d5  =   types.Double(value=output.get('d5',None)),
                    d7  =   types.Double(value=output.get('d7',None)),
                    d8  =   types.Int64(value=output.get('d8',None)),
                    A0  =   types.Int64(value=output.get('A0',None)),
                    A1  =   types.Double(value=output.get('A1',None)),
                    A2  =   types.Int64(value=output.get('A2',None)),
                    A3  =   types.Int64(value=output.get('A3',None)),
                    A4  =   types.Double(value=output.get('A4',None)),
                    A5  =   types.Int64(value=output.get('A5',None)),
                    A6  =   types.Double(value=output.get('A6',None)),
                    A7  =   types.Double(value=output.get('A7',None)),
                    A8  =   types.Double(value=output.get('A8',None)),
                    A9  =   types.Double(value=output.get('A9',None)),
                    F0  =   types.Int64(value=output.get('F0',None)),
                    F1  =   types.Double(value=output.get('F1',None)),
                    F2  =   types.Int64(value=output.get('F2',None)),
                    F3  =   types.Double(value=output.get('F3',None)),
                    Hd1  =   types.Double(value=output.get('Hd1',None)),
                    Hd2  =   types.Double(value=output.get('Hd2',None)),
                    Hd3  =   types.Double(value=output.get('Hd3',None)),
                    Hd4  =   types.Double(value=output.get('Hd4',None)),
                    Hd5  =   types.Double(value=output.get('Hd5',None)),
                    Hd6  =   types.Double(value=output.get('Hd6',None))
                )
            )
        )
        print(self.report(self.service_name, msg))


if __name__ == '__main__':
    with open('parker.yaml') as f:
        # use safe_load instead load for security reasons
        driverConfig = yaml.safe_load(f)

    namespace = driverConfig['wavemq']['namespace']
    service_name = driverConfig['xbos']['service_name']
    #driver_cfg = "parker.yaml"
    print(driverConfig)

    xbos_cfg = {
        'wavemq': 'localhost:4516',
        'namespace': namespace,
        'base_resource': 'parker',
        'entity': 'parker.ent',
        'id': 'pyxbos-driver-parker-2',
        #'rate': 1800, # half hour
        'rate': 20, # 15 min
        'service_name': service_name
    }
    print(getmembers(iot_pb2))
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    #e = DarkSkyDriver(cfg)
    e = ParkerDriver(xbos_cfg)
    e.begin()
