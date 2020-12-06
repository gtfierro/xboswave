from pyxbos import flexstat_pb2
from pyxbos.driver import *
import yaml
import argparse
import time
import BAC0


class FlexstatDriver(Driver):
    def setup(self, cfg):
        config_file = cfg['config_file']
        with open(config_file) as f:
            driverConfig = yaml.safe_load(f)

        self.service_name_map = cfg['service_name_map']
        self.thermostat_config = driverConfig['thermostat_config']

        bacnet_mask = self.thermostat_config.get('bacnet_network_mask')
        bacnet_router_address = self.thermostat_config.get('bacnet_router_address', None)
        bbmd_ttl = self.thermostat_config.get('bbmd_ttl', None)
        if bacnet_router_address == None or bbmd_ttl == None:
            self.bacnet = BAC0.connect(ip=bacnet_mask)
        else:
            self.bacnet = BAC0.connect(ip=bacnet_mask, bbmdAddress=bacnet_router_address, bbmdTTL=bbmd_ttl)
        self.point_map = self.thermostat_config['point_map']

        self.device_map = {}
        for service_name in self.service_name_map:
            ip = self.service_name_map[service_name].get('ip')
            device_id = self.service_name_map[service_name].get('device_id', 1)

            device = BAC0.device(address=ip, device_id=device_id, network=self.bacnet)
            self.device_map[service_name] = device

    def read(self, requestid=None):
        for service_name in self.device_map:
            device = self.device_map[service_name]

            try:
                measurements = {}
                for point in self.point_map:
                    bacnet_point_name = self.point_map[point]
                    val = device[bacnet_point_name].value
                    if type(val) == str:
                        if val == "active":
                            val = True
                        else:
                            val = False
                    measurements[point] = val
                print(measurements)

                msg = xbos_pb2.XBOS(
                    flexstat_state=flexstat_pb2.FlexstatState(
                        time=int(time.time() * 1e9),
                        space_temp_sensor=types.Double(value=measurements.get('space_temp_sensor', None)),
                        minimum_proportional=types.Double(value=measurements.get('minimum_proportional', None)),
                        active_cooling_setpt=types.Double(value=measurements.get('active_cooling_setpt', None)),
                        active_heating_setpt=types.Double(value=measurements.get('active_heating_setpt', None)),
                        unocc_cooling_setpt=types.Double(value=measurements.get('unocc_cooling_setpt', None)),
                        unocc_heating_setpt=types.Double(value=measurements.get('unocc_heating_setpt', None)),
                        occ_min_clg_setpt=types.Double(value=measurements.get('occ_min_clg_setpt', None)),
                        occ_max_htg_setpt=types.Double(value=measurements.get('occ_max_htg_setpt', None)),
                        override_timer=types.Double(value=measurements.get('override_timer', None)),
                        occ_cooling_setpt=types.Double(value=measurements.get('occ_cooling_setpt', None)),
                        occ_heating_setpt=types.Double(value=measurements.get('occ_heating_setpt', None)),
                        current_mode_setpt=types.Double(value=measurements.get('current_mode_setpt', None)),
                        ui_setpt=types.Double(value=measurements.get('ui_setpt', None)),
                        cooling_need=types.Double(value=measurements.get('cooling_need', None)),
                        heating_need=types.Double(value=measurements.get('heating_need', None)),
                        unocc_min_clg_setpt=types.Double(value=measurements.get('unocc_min_clg_setpt', None)),
                        unocc_max_htg_setpt=types.Double(value=measurements.get('unocc_max_htg_setpt', None)),
                        min_setpt_diff=types.Double(value=measurements.get('min_setpt_diff', None)),
                        min_setpt_limit=types.Double(value=measurements.get('min_setpt_limit', None)),
                        space_temp=types.Double(value=measurements.get('space_temp', None)),
                        cooling_prop=types.Double(value=measurements.get('cooling_prop', None)),
                        heating_prop=types.Double(value=measurements.get('heating_prop', None)),
                        cooling_intg=types.Double(value=measurements.get('cooling_intg', None)),
                        heating_intg=types.Double(value=measurements.get('heating_intg', None)),
                        fan=types.Int64(value=measurements.get('fan', None)),
                        occupancy_mode=types.Int64(value=measurements.get('occupancy_mode', None)),
                        setpt_override_mode=types.Int64(value=measurements.get('setpt_override_mode', None)),
                        fan_alarm=types.Int64(value=measurements.get('fan_alarm', None)),
                        fan_need=types.Int64(value=measurements.get('fan_need', None)),
                        heating_cooling_mode=types.Int64(value=measurements.get('heating_cooling_mode', None)),
                        occ_fan_auto_on=types.Int64(value=measurements.get('occ_fan_auto_on', None)),
                        unocc_fan_auto_on=types.Int64(value=measurements.get('unocc_fan_auto_on', None)),
                        fan_status=types.Int64(value=measurements.get('fan_status', None))
                    )
                )
                print(self.report(service_name, msg))
            except:
                print("error for thermostat {0} !! continuing!".format(service_name))


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
    thermostat_config = xbosConfig.get('thermostat_config')

    xbos_cfg = {
        'waved': waved,
        'wavemq': wavemq,
        'namespace': namespace,
        'base_resource': base_resource,
        'entity': entity,
        'id': driver_id,
        'rate': rate,
        'service_name_map': service_name_map,
        'config_file': config_file,
        'thermostat_config': thermostat_config
    }

    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
    e = FlexstatDriver(xbos_cfg)
    e.begin()
