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
        print('write',uri, age)
        if msg.XBOSIoTDeviceActuation is not None:
            l = msg.XBOSIoTDeviceActuation.light
            name = uri.split('/')[-1]
            self.lights[name].on = l.state
            self.lights[name].brightness = l.brightness


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

e = Example()
e.begin()
