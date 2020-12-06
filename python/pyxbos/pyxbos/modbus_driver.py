# coding: utf-8

# In[61]:


#!/usr/bin/env python
"""
Adapted from code sample code from the Pymodbus Examples
***Created by Andrew Shephard
***Maintained by Anand Prakash <akprakash@lbl.gov>
"""

#---------------------------------------------------------------------------#
# import the required server implementation
#---------------------------------------------------------------------------#
from pymodbus.client.sync import ModbusTcpClient
from pymodbus.client.sync import ModbusSerialClient

from pymodbus.constants import Endian
from pymodbus.payload import BinaryPayloadDecoder
from pymodbus.payload import BinaryPayloadBuilder
import configparser
import yaml
import logging


class Modbus_Driver(object):
    def __init__(self, config_file, config_section='modbus', **kwargs):
        # Use a config section if the config file is being shared with other
        # parts of a project. **kwargs can contain a variable amount of
        if (isinstance(config_file,str)):

            with open(config_file) as f:
                modbusConfig = yaml.safe_load(f)
        else:
            modbusConfig = config_file

        modbus_section = config_section

        self.BYTE_ORDER_DICT = {}
        self.WORD_ORDER_DICT = {}
        self.input_register_dict = {}
        self.holding_register_dict = {}
        self.coil_register_dict = {}
        self.discrete_register_dict = {}

        self.input_registers = {}
        self.holding_registers = {}
        self.coil_registers = {}
        self.discrete_registers = {}
        self.MODBUS_TYPE = modbusConfig[modbus_section]['modbus_type']
        # Check to see if unit id is a list, if it is then set flag that it is a
        # list
        self.UNIT_ID = modbusConfig[modbus_section]['UNIT_ID']
        if isinstance(self.UNIT_ID, list):
            self.UNIT_ID_LIST = self.UNIT_ID
            #Set default UNIT_ID as first UNIT_ID in list
            self.UNIT_ID = int(self.UNIT_ID_LIST[0])
        else:
            # Make a unit id list from the non-list definition for compatibility
            # reasons of previous configs. This also eliminates the possibility
            # of error in calling get_data_all_devices() on a config with a non
            # list definition
            self.UNIT_ID_LIST = []
            self.UNIT_ID_LIST.append(self.UNIT_ID)
        # Start logging if enabled in config
        self.LOGGING_FLAG = modbusConfig[modbus_section]['enable_logging']
        if self.LOGGING_FLAG == False:
            #Start client logging for trouble shooting
            logging.basicConfig()
            log = logging.getLogger()
            log.setLevel(logging.ERROR)

        # Start appropriate client based on the type specified in the config
        if self.MODBUS_TYPE == 'serial':
            self.METHOD = modbusConfig[modbus_section]['method']
            self.SERIAL_PORT = modbusConfig[modbus_section]['serial_port']
            self.STOPBITS = modbusConfig[modbus_section]['stopbits']
            self.BYTESIZE = modbusConfig[modbus_section]['bytesize']
            self.PARITY = modbusConfig[modbus_section]['parity']
            self.BAUDRATE = modbusConfig[modbus_section]['baudrate']
        elif self.MODBUS_TYPE == 'tcp':
            self.IP_ADDRESS = modbusConfig[modbus_section]['ip']
            self.PORT = modbusConfig[modbus_section]['port']
        else:
            print("Invalid modbus type")
            exit

        # Set the byte order as big or little endian
        if modbusConfig[modbus_section]['byte_order'] == 'big':
            self.BYTE_ORDER = Endian.Big
            self.BYTE_ORDER_DICT[self.UNIT_ID] = Endian.Big
        elif modbusConfig[modbus_section]['byte_order'] == 'little':
            self.BYTE_ORDER = Endian.Little
            self.BYTE_ORDER_DICT[self.UNIT_ID] = Endian.Little
        else:
            print("invalid byte order") # change to except later
            exit()
        # Set the word order as big or little endian
        if modbusConfig[modbus_section]['word_order'] == 'big':
            self.WORD_ORDER = Endian.Big
            self.WORD_ORDER_DICT[self.UNIT_ID] = Endian.Big

        elif modbusConfig[modbus_section]['word_order'] == 'little':
            self.WORD_ORDER = Endian.Little
            self.WORD_ORDER_DICT[self.UNIT_ID] = Endian.Little
        else:
            print("invalid byte order") # change to except later
            exit()

        # Read in all registers specified in the YAML config
        self.coil_register_dict = modbusConfig[modbus_section]['coil_registers']
        self.discrete_register_dict = modbusConfig[modbus_section]['discrete_registers']
        self.holding_register_dict = modbusConfig[modbus_section]['holding_registers']
        self.input_register_dict = modbusConfig[modbus_section]['input_registers']

        self.coil_registers[self.UNIT_ID] = self.coil_register_dict
        self.discrete_registers[self.UNIT_ID] = self.discrete_register_dict
        self.holding_registers[self.UNIT_ID] = self.holding_register_dict
        self.input_registers[self.UNIT_ID] = self.input_register_dict
        #print(self.holding_registers)

        # Add single device that is either specified in config_section parameter
        # or a single device config file
        for current_device in self.UNIT_ID_LIST:
            self.coil_registers[current_device] = self.coil_register_dict
            self.discrete_registers[current_device] = self.discrete_register_dict
            self.holding_registers[current_device] = self.holding_register_dict
            self.input_registers[current_device] = self.input_register_dict
            # Set the byte order as big or little endian
            if modbusConfig[modbus_section]['byte_order'] == 'big':
                self.BYTE_ORDER_DICT[current_device] = Endian.Big
            elif modbusConfig[modbus_section]['byte_order'] == 'little':
                self.BYTE_ORDER_DICT[current_device] = Endian.Little
            else:
                print("invalid byte order") # change to except later
                exit()
            # Set the word order as big or little endian
            if modbusConfig[modbus_section]['word_order'] == 'big':
                self.WORD_ORDER_DICT[current_device] = Endian.Big
            elif modbusConfig[modbus_section]['word_order'] == 'little':
                self.WORD_ORDER_DICT[current_device] = Endian.Little
            else:
                print("invalid word order") # change to except later
                exit()

        # Apply register offset if specified
        self.OFFSET_REGISTERS = modbusConfig[modbus_section]['OFFSET_REGISTERS']
        for key in self.holding_register_dict:
            self.holding_register_dict[key][0] -= self.OFFSET_REGISTERS

        # Add devices that were specified with **kwargs
        for device_name, modbus_section in kwargs.items():
            # The Device ID is used as the key in a dictionary for all settings
            # that could potentially differ between devices. Since all of the
            # functions already have been updated to take in a UNIT_ID this
            # can be used to retrieve the appropriate setting for the device.

            # TODO Handle case where the config section has a list of the same
            # device.

            current_device = modbusConfig[modbus_section]['UNIT_ID']
            #print(type(current_device))
            # TODO make this a for loop for each ID
            #current_device = current_device[0]
            self.UNIT_ID_LIST.append(int(current_device))

            if modbusConfig[modbus_section]['byte_order'] == 'big':
                self.BYTE_ORDER_DICT[current_device] = Endian.Big
            elif modbusConfig[modbus_section]['byte_order'] == 'little':
                self.BYTE_ORDER_DICT[current_device] = Endian.Little
            # Set the word order as big or little endian
            if modbusConfig[modbus_section]['word_order'] == 'big':
                self.WORD_ORDER_DICT[current_device] = Endian.Big
            elif modbusConfig[modbus_section]['word_order'] == 'little':
                self.WORD_ORDER_DICT[current_device] = Endian.Little
            else:
                print("invalid word order") # change to except later
                exit()
            # Read in all registers specified in the YAML config
            self.coil_register_dict = modbusConfig[modbus_section]['coil_registers']
            self.discrete_register_dict = modbusConfig[modbus_section]['discrete_registers']
            self.holding_register_dict = modbusConfig[modbus_section]['holding_registers']
            self.input_register_dict = modbusConfig[modbus_section]['input_registers']

            self.coil_registers[current_device] = self.coil_register_dict
            self.discrete_registers[current_device] = self.discrete_register_dict
            self.holding_registers[current_device] = self.holding_register_dict
            self.input_registers[current_device] = self.input_register_dict
            #print(self.holding_register_dict)
            '''
            # Read in all registers specified in the YAML config
            self.coil_registers[current_device] = modbusConfig[modbus_section]['coil_registers']
            self.discrete_registers[current_device] = modbusConfig[modbus_section]['discrete_registers']
            self.holding_registers[current_device] = modbusConfig[modbus_section]['holding_registers']
            self.input_register_dict[current_device] = modbusConfig[modbus_section]['input_registers']
            '''

            # Apply register offset if specified
            # TODO fix this for one device as well as multiple
            """
            self.OFFSET_REGISTERS_DICT[current_device] = modbusConfig[modbus_section]['OFFSET_REGISTERS']
            for key in self.holding_register_dict:
                self.holding_register_dict[key][0] -= self.OFFSET_REGISTERS
            """
            #print(self.holding_registers)



    def initialize_modbus(self):
        """
        initalize correct client according to type specified in config:
            'tcp' or 'serial'
        """
        if self.MODBUS_TYPE == 'serial':
            self.client= ModbusSerialClient(
                    method      = self.METHOD,
                    port        = self.SERIAL_PORT,
                    stopbits    = self.STOPBITS,
                    bytesize    = self.BYTESIZE, 
                    parity      = self.PARITY,
                    baudrate    = self.BAUDRATE
                )
            connection = self.client.connect()

        if self.MODBUS_TYPE == 'tcp':
            self.client = ModbusTcpClient(self.IP_ADDRESS,port=self.PORT)
        '''
        rr = self.read_register_raw(0x601,1,247)
        decoder = BinaryPayloadDecoder.fromRegisters(
                rr.registers,
                byteorder=self.BYTE_ORDER,
                wordorder=self.WORD_ORDER)
        output = decoder.decode_16bit_int()
        print(output)
        '''
        #rr = self.read_register_raw(1001,2,7)
        '''decoder = BinaryPayloadDecoder.fromRegisters(
                rr.registers,
                byteorder=self.BYTE_ORDER,
                wordorder=self.WORD_ORDER)
        '''


    def reconnect(self):
        try:
            self.client.close()
        finally:
            self.initialize_modbus()

    def write_single_register(self,register,value, unit=None):
        """
        :param register: address of reigster to write
        :param value: Unsigned short
        :returns: Status of write
        """
        if (unit is None):
            unit = self.UNIT_ID
        response = self.client.write_register(register,value,unit)
        return response

    def write_data(self,register,value):
        response = self.client.write_register(register,value,unit= self.UNIT_ID)
        return response

    def write_register(self,register_name,value, unit=None):
        """
        :param register_name: register key from holding register dictionary
            generated by yaml config
        :param value: value to write to register
        :returns: -- Nothing
        """
        # TODO add the ability to discern which settings will be appropriate for
        # the device that is being written to
        if (unit is None):
            unit = self.UNIT_ID
        '''
        builder = BinaryPayloadBuilder(byteorder=self.BYTE_ORDER,
            wordorder=self.WORD_ORDER_DICT[unit])
        '''
        builder = BinaryPayloadBuilder(byteorder=self.BYTE_ORDER_DICT[unit],
            wordorder=self.WORD_ORDER_DICT[unit])
        # This will change depending on the device that is being connected
        # potentially so it has to be correleated to the device ID

        if (self.holding_register_dict[register_name][1] == '8int'):
            builder.add_8bit_int(value)
        elif (self.holding_register_dict[register_name][1] == '8uint'):
            builder.add_8bit_uint(value)
        elif (self.holding_register_dict[register_name][1] == '16int'):
            builder.add_16bit_int(value)
        elif (self.holding_register_dict[register_name][1] == '16uint'):
            builder.add_16bit_uint(value)
        elif (self.holding_register_dict[register_name][1] == '32int'):
            builder.add_32bit_int(value)
        elif (self.holding_register_dict[register_name][1] == '32uint'):
            builder.add_32bit_uint(value)
        elif (self.holding_register_dict[register_name][1] == '32float'):
            builder.add_32bit_float(value)
        elif (self.holding_register_dict[register_name][1] == '64int'):
            builder.add_64bit_int(value)
        elif (self.holding_register_dict[register_name][1] == '64uint'):
            builder.add_64bit_uint(value)
        elif (self.holding_register_dict[register_name][1] == '64float'):
            builder.add_64bit_float(value)
        else:
            print("Bad type")
            exit()
        payload = builder.build()
        self.client.write_registers(self.holding_register_dict[register_name][0],
            payload, skip_encode=True, unit = self.UNIT_ID)

    def write_coil(self,register,value, unit=None):
        """
        :param register_name: register key from holding register dictionary
            generated by yaml config
        :param value: value to write to register
        :returns:
        """
        # TODO mention what type the value needs to be for value
        if (unit is None):
            unit = self.UNIT_ID

        response = self.client.write_coil(register,value,unit)
        return response

    def read_coil(self,register, unit=None):
        """
        :param register: coil register address to read
        :returns: value stored in coil register
        """
        # TODO mention what type the value needs to be for value
        if (unit is None):
            unit = self.UNIT_ID

        rr = self.client.read_coils(register, 1, unit=unit)
        return rr.bits[0]

    def read_discrete(self,register,unit=None):
        """
        :param register: discrete register address to read
        :returns: value stored in coil register
        """
        if (unit is None):
            unit = self.UNIT_ID

        rr = self.client.read_discrete_inputs(register, count=1,unit=unit)
        return rr.bits[0]

    def read_register_raw(self,register,length, unit=None):
        """
        :param register: base holding register address to read
        :param length: amount of registers to read to encompass all of the data necessary
            for the type
        :returns: A deferred response handle
        """
        if (unit is None):
            unit = self.UNIT_ID

        response = self.client.read_holding_registers(register,length,unit=unit)
        return response

    def read_input_raw(self,register,length, unit=None):
        """
        :param register: base input register address to read
        :param length: amount of registers to read to encompass all of the data necessary
            for the type
        :returns: A deferred response handle
        """
        if (unit is None):
            unit = self.UNIT_ID

        response = self.client.read_input_registers(register,length,unit=unit)
        return response

    def decode_register(self,register,type, unit=None):
        #print(unit)
        #print(type(unit))
        """
        :param register: holding register address to retrieve
        :param type: type to interpret the registers retrieved as
        :returns: data in the type specified

        Based on the type provided, this function retrieves the values contained
        in the register address specfied plus the amount necessary to encompass
        the the type. For example, if 32int is specified with an address of 200
        the registers accessed would be 200 and 201.

        The types accepted are listed in the table below along with their length
        |   Type          | Length (registers) |
        | ------------- |:------------------:|
        |        ignore |                  1 |
        |          8int |                  1 |
        |         8uint |                  1 |
        |         16int |                  1 |
        |        16uint |                  1 |
        |         32int |                  2 |
        |        32uint |                  2 |
        |       32float |                  2 |
        |         64int |                  4 |
        |        64uint |                  4 |
        |       64float |                  4 |
        """
        if (unit is None):
            unit = self.UNIT_ID
        #omitting string for now since it requires a specified length
        if type == '8int':
            rr = self.read_register_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_8bit_int()

        elif type == '8uint':
            rr = self.read_register_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_8bit_uint()
        elif type == '16int':
            rr = self.read_register_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_16bit_int()
        elif type == '16uint':
            rr = self.read_register_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_16bit_uint()
        elif type == '32int':
            rr = self.read_register_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_int()
        elif type == '32uint':
            rr = self.read_register_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_uint()
        elif type == '32float':
            rr = self.read_register_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_float()
        elif type == '64int':
            rr = self.read_register_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_int()
        elif type == '64uint':
            rr = self.read_register_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_uint()
        elif type == 'ignore':
            rr = self.read_register_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.skip_bytes(8)
        elif type == '64float':
            rr = self.read_register_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_float()
        else:
            print("Wrong type specified")
            exit()

        return output

    def decode_input_register(self,register,type, unit=None):
        """
        :param register: input register address to retrieve
        :param type: type to interpret the registers retrieved as
        :returns: data in the type specified

        Based on the type provided, this function retrieves the values contained
        in the register address specfied plus the amount necessary to encompass
        the the type. For example, if 32int is specified with an address of 200
        the registers accessed would be 200 and 201.

        The types accepted are listed in the table below along with their length
        |   Type          | Length (registers) |
        | ------------- |:------------------:|
        |        ignore |                  1 |
        |          8int |                  1 |
        |         8uint |                  1 |
        |         16int |                  1 |
        |        16uint |                  1 |
        |         32int |                  2 |
        |        32uint |                  2 |
        |       32float |                  2 |
        |         64int |                  4 |
        |        64uint |                  4 |
        |       64float |                  4 |
        """
        if (unit is None):
            unit = self.UNIT_ID
        #omitting string for now since it requires a specified length
        if type == '8int':
            rr = self.read_input_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_8bit_int()

        elif type == '8uint':
            rr = self.read_input_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit][unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_8bit_uint()
        elif type == '16int':
            rr = self.read_input_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_16bit_int()
        elif type == '16uint':
            rr = self.read_input_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_16bit_uint()
        elif type == '32int':
            rr = self.read_input_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_int()
        elif type == '32uint':
            rr = self.read_input_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_uint()
        elif type == '32float':
            rr = self.read_input_raw(register,2,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_32bit_float()
        elif type == '64int':
            rr = self.read_input_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_int()
        elif type == '64uint':
            rr = self.read_input_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_uint()
        elif type == 'ignore':
            rr = self.read_input_raw(register,1,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.skip_bytes(8)
        elif type == '64float':
            rr = self.read_input_raw(register,4,unit)
            decoder = BinaryPayloadDecoder.fromRegisters(
                    rr.registers,
                    byteorder=self.BYTE_ORDER_DICT[unit],
                    wordorder=self.WORD_ORDER_DICT[unit])
            output = decoder.decode_64bit_float()
        else:
            print("Wrong type specified")
            exit()

        return output

    def read_register(self,register_name):
        response = self.decode_register(self.holding_register_dict[register_name][0],
            self.holding_register_dict[register_name][1])
        return response

    def read_input_raw(self,register_name):
        response = self.decode_input_register(self.holding_register_dict[register_name][0],
            self.holding_register_dict[register_name][1])
        return response


    def get_data(self,unit=None):
        """
        :returns: Dictionary containing the value retrieved for each register
        contained in the YAML config file, register names cannot be repeated
        or the register will be overwritten
        """
        output = {}

        if unit is None:
            unit = self.UNIT_ID

        for key in self.coil_registers[unit]:
            output[key] = self.read_coil(self.coil_registers[unit][key][0],unit)

        for key in self.discrete_registers[unit]:
            output[key] = self.read_discrete(self.discrete_registers[unit][key][0],unit)

        for key in self.input_registers[unit]:
            output[key] = self.decode_input_register(self.input_registers[unit][key][0],self.input_registers[unit][key][1],unit)
            
        for key in self.holding_registers[unit]:
            if (len(self.holding_registers[unit][key]) == 3):
                # Check Read/Write Flag
                if (self.holding_registers[unit][key][2].find('R') != -1):
                    output[key] = self.decode_register(self.holding_registers[unit][key][0],self.holding_registers[unit][key][1],unit)
            else:
                # Register list does not contain a Read/Write Flag assume R
                output[key] = self.decode_register(self.holding_registers[unit][key][0],self.holding_registers[unit][key][1],unit)

        return output

    def get_data_all_devices(self):
        reg_data_dict = {}
        cnt = 1
        for dev_id in self.UNIT_ID_LIST:
            new_key = str(dev_id)
            if str(dev_id) in reg_data_dict:
                new_key = new_key + '_' + str(cnt)
                cnt += 1
            reg_data_dict[new_key] = self.get_data(dev_id)
        return reg_data_dict
    def kill_modbus(self):
        """
        Closes connection with Modbus Slave
        """
        self.client.close()
