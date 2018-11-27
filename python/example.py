from pyxbos import *
import time
import random


class Example(Driver):
    def setup(self):
        print("setting up")

    def read(self):
        temp = 20 + random.random() * 70
        msg = xbos_pb2.XBOS(
            XBOSIoTDeviceState = iot_pb2.XBOSIoTDeviceState(
                time = int(time.time()*1e9),
                thermostat = iot_pb2.Thermostat(
                    temperature = temp,
                )
            )
        )
        yield (msg, "temp")

e = Example()
e.begin()
