from pyxbos import *
import time
from phue import Bridge


class Example(Driver):
    def setup(self):
        print("setting up")

        self.b = Bridge('192.168.1.84')
        self.b.connect()
        self.b.get_api()



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
