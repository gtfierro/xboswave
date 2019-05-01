from pyxbos.driver import *
import logging
import time
from phue import Bridge

class HueLight:
    def __init__(self, phue_light, reportfunc, reporturi):
        self._report = reportfunc
        self.reporturi = reporturi
        self.light = phue_light

    def report(self, msg):
        self._report(self.reporturi, msg)

    def read(self,requestid=None):
        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                light = iot_pb2.Light(
                    state = types.Bool(value=self.light.on),
                    brightness = types.Int64(value=self.light.on*int(100 * (self.light.brightness / 254.))),
                )
            )
        )
        if requestid is not None:
            msg.XBOSIoTDeviceState.requestid = requestid
        return msg

    def write(self, uri, age, msg):
        if age > 30: return # cutoff actuation @ 30 seconds
        l = msg.XBOSIoTDeviceActuation.light
        if l.state is not None:
            self.light.on = l.state.value
        if self.light.on and l.brightness is not None:
            self.light.brightness = int(254 * (l.brightness.value / 100.))

        self.report(self.read(msg.XBOSIoTDeviceActuation.requestid))

class HueDriver(Driver):
    def setup(self, cfg):
        self.b = Bridge(cfg['hue_bridge'])
        self.b.connect()
        self.b.get_api()
        self.lights = {}
        for l in self.b.lights:
            self.lights[l.name.replace(' ','_')] = HueLight(l, self.report, l.name.replace(" ","_"))
        self._log.info("lights: {0}".format(self.b.lights))

    def write(self, uri, age, msg):
        """
        Dispatch the write to the correct object
        """
        if msg.XBOSIoTDeviceActuation is not None:
            l = msg.XBOSIoTDeviceActuation.light
            name = uri.split('/')[-1]
            self.lights[name].write(uri, age, msg)

    def read(self, requestid=None):
        """
        Read all sub devices and have them publish
        """
        for light in self.lights.values():
            light.report(light.read())

logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

if __name__ == '__main__':
    cfg = {
        'hue_bridge': '192.168.1.84',
        'wavemq': 'localhost:4516',
        'waved': 'localhost:410',
        'namespace': 'GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==',
        'base_resource': 'test/hue',
        'entity': 'gabehue.ent',
        'id': 'pyxbos-driver-hue-1',
        'rate': 10,
    }
    e = HueDriver(cfg)
    e.begin()
