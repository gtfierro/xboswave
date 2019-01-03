from pyxbos import *
import logging
import time
from phue import Bridge


class Example(Driver):
    def setup(self, cfg):
        print("setting up")

        self.b = Bridge(cfg['hue_bridge'])
        self.b.connect()
        self.b.get_api()
        self.lights = {}
        for l in self.b.lights:
            self.lights[l.name.replace(' ','_')] = l
        print(self.b.lights)

    def write(self, uri, age, msg):
        if msg.XBOSIoTDeviceActuation is not None:
            l = msg.XBOSIoTDeviceActuation.light
            name = uri.split('/')[-1]
            self.lights[name].on = l.state
            self.lights[name].brightness = int(254 * (l.brightness / 100.))

    def read(self):
        for light in self.b.lights:
            msg = xbos_pb2.XBOS(
                XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                    time = int(time.time()*1e9),
                    light = iot_pb2.Light(
                        state = light.on,
                        brightness = int(100 * (light.brightness / 254.)),
                    )
                )
            )
            yield (msg, light.name.replace(" ","_"))

logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

cfg = {
    'hue_bridge': '<address>',
    'wavemq': 'localhost:4516',
    'waved': 'localhost:410',
    'namespace': '<namespace>',
    'base_resource': 'hue/bridge1',
    'base_resource_write': 'hue/bridge1/cmd/*',
    'entity': '<entity path>',
    'id': 'abctest123098745879',
    'rate': 10,
    'write_expect': xbos_pb2.XBOS,
}
e = Example(cfg)
e.begin()
