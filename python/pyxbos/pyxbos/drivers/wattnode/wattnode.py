from pyxbos import wattnode_pb2
from pyxbos.driver import *
from pyxbos.modbus_driver import Modbus_Driver
import yaml
import argparse
import time
from inspect import getmembers


class WattnodeDriver(Driver):
    def setup(self, cfg):
        config_file = cfg['config_file']

        self.modbus_device = Modbus_Driver(config_file=config_file, config_section='modbus')
        self.modbus_device.initialize_modbus()

        self.service_name_map = cfg['service_name_map']

    def read(self, requestid=None):
        for service_name in self.service_name_map:
            unit_id = self.service_name_map[service_name]

            try:

                output = self.modbus_device.get_data(unit=unit_id)
                msg = xbos_pb2.XBOS(
                    wattnode_state = wattnode_pb2.WattnodeState(
                        time = int(time.time()*1e9),
                        EnergySum  =   types.Double(value=output.get('EnergySum',None)),
                        EnergyPosSum  =   types.Double(value=output.get('EnergyPosSUm',None)),
                        EnergySumNR  =   types.Double(value=output.get('EnergySumNR',None)),
                        EnergyPosSumNr  =   types.Double(value=output.get('EnergyPosSumNr',None)),
                        PowerSum  =   types.Double(value=output.get('PowerSum',None)),
                        PowerA  =   types.Double(value=output.get('PowerA',None)),
                        PowerB  =   types.Double(value=output.get('PowerB',None)),
                        PowerC  =   types.Double(value=output.get('PowerC',None)),
                        VoltAvgLN  =   types.Double(value=output.get('VoltAvgLN',None)),
                        VoltA  =   types.Double(value=output.get('VoltA',None)),
                        VoltB  =   types.Double(value=output.get('VoltB',None)),
                        VoltC  =   types.Double(value=output.get('VoltC',None)),
                        VoltAvgLL  =   types.Double(value=output.get('VoltAvgLL',None)),
                        VoltAB  =   types.Double(value=output.get('VoltAB',None)),
                        VoltBC  =   types.Double(value=output.get('VoltBC',None)),
                        VoltAC  =   types.Double(value=output.get('VoltAC',None)),
                        Freq  =   types.Double(value=output.get('Freq',None)),
                        EnergyA  =   types.Double(value=output.get('EnergyA',None)),
                        EnergyB  =   types.Double(value=output.get('EnergyB',None)),
                        EnergyC  =   types.Double(value=output.get('EnergyC',None)),
                        EnergyPosA  =   types.Double(value=output.get('EnergyPosA',None)),
                        EnergyPosB  =   types.Double(value=output.get('EnergyPosB',None)),
                        EnergyPosC  =   types.Double(value=output.get('EnergyPosC',None)),
                        EnergyNegSum  =   types.Double(value=output.get('EnergyNegSum',None)),
                        EnergyNegSumNR  =   types.Double(value=output.get('EnergyNegSumNR',None)),
                        EnergyNegA  =   types.Double(value=output.get('EnergyNegA',None)),
                        EnergyNegB  =   types.Double(value=output.get('EnergyNegB',None)),
                        EnergyNegC  =   types.Double(value=output.get('EnergyNegC',None)),
                        EnergyReacSum  =   types.Double(value=output.get('EnergyReacSum',None)),
                        EnergyReacA  =   types.Double(value=output.get('EnergyReacA',None)),
                        EnergyReacB  =   types.Double(value=output.get('EnergyReacB',None)),
                        EnergyReacC  =   types.Double(value=output.get('EnergyReacC',None)),
                        EnergyAppSum  =   types.Double(value=output.get('EnergyAppSum',None)),
                        EnergyAppA  =   types.Double(value=output.get('EnergyAppA',None)),
                        EnergyAppB  =   types.Double(value=output.get('EnergyAppB',None)),
                        EnergyAppC  =   types.Double(value=output.get('EnergyAppC',None)),
                        PowerFactorAvg  =   types.Double(value=output.get('PowerFactorAvg',None)),
                        PowerFactorA  =   types.Double(value=output.get('PowerFactorA',None)),
                        PowerFactorB  =   types.Double(value=output.get('PowerFactorB',None)),
                        PowerFactorC  =   types.Double(value=output.get('PowerFactorC',None)),
                        PowerReacSum  =   types.Double(value=output.get('PowerReacSum',None)),
                        PowerReacA  =   types.Double(value=output.get('PowerReacA',None)),
                        PowerReacB  =   types.Double(value=output.get('PowerReacB',None)),
                        PowerReacC  =   types.Double(value=output.get('PowerReacC',None)),
                        PowerAppSum  =   types.Double(value=output.get('PowerAppSum',None)),
                        PowerAppA  =   types.Double(value=output.get('PowerAppA',None)),
                        PowerAppB  =   types.Double(value=output.get('PowerAppB',None)),
                        PowerAppC  =   types.Double(value=output.get('PowerAppC',None)),
                        CurrentA  =   types.Double(value=output.get('CurrentA',None)),
                        CurrentB  =   types.Double(value=output.get('CurrentB',None)),
                        CurrentC  =   types.Double(value=output.get('CurrentC',None)),
                        Demand  =   types.Double(value=output.get('Demand',None)),
                        DemandMin  =   types.Double(value=output.get('DemandMin',None)),
                        DemandMax  =   types.Double(value=output.get('DemandMax',None)),
                        DemandApp  =   types.Double(value=output.get('DemandApp',None)),
                        DemandA  =   types.Double(value=output.get('DemandA',None)),
                        DemandB  =   types.Double(value=output.get('DemandB',None)),
                        DemandC  =   types.Double(value=output.get('DemandC',None))
                    )
                )
                print(self.report(service_name, msg))
            except:
                print("error occured! reconnecting and continuing")
                self.modbus_device.reconnect()

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("config_file", help="config file with api key as well as namespace")
    args = parser.parse_args()
    config_file = args.config_file
    
    with open(config_file) as f:
        driverConfig = yaml.safe_load(f)

    xbosConfig = driverConfig['xbos']
    waved = xbosConfig.get('waved', 'localhost:777')
    wavemq = xbosConfig.get('wavemq', 'locahost:4516')
    namespace = xbosConfig.get('namespace')
    base_resource = xbosConfig.get('base_resource')
    service_name_map = xbosConfig.get('service_name_map')
    entity = xbosConfig.get('entity')
    rate = xbosConfig.get('rate')
    driver_id = xbosConfig.get('id', 'wattnode-driver')

    xbos_cfg = {
        'waved': waved,
        'wavemq': wavemq,
        'namespace': namespace,
        'base_resource': base_resource,
        'entity': entity,
        'id': driver_id,
        'rate': rate, 
        'service_name_map': service_name_map,
        'config_file': config_file
    }

    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    e = WattnodeDriver(xbos_cfg)
    e.begin()
