import grpc
import pickle
import base64

from .eapi_pb2 import *
from .wavemq_pb2 import *
from .wavemq_pb2_grpc import *
from . import xbos_pb2
from . import iot_pb2

import asyncio

WaveBuiltinPSET = b"\x1b\x20\x19\x49\x54\xe8\x6e\xeb\x8f\x91\xff\x98\x3a\xcc\x56\xe6\xc8\x4a\xe2\x9a\x90\x7c\xe7\xe7\x63\x8e\x86\x57\xd5\x14\x99\xb1\x88\xa4"
WaveGlobalNamespace = b"\x1b\x20\xcf\x8d\x19\xd7\x9d\x23\x01\x38\x65\xbe\xf7\x57\xce\xa0\x4c\xde\xe5\xef\x4e\xde\xfc\x80\x8d\xd2\x1e\x4e\x00\x5e\x6f\x80\x47\xcc"

WaveBuiltinE2EE = "decrypt"

# TODO:
# - need to load regular entity from disk using Python

class Driver:
    def __init__(self):

        self.cfg = {
            "wavemq": "localhost:4516",
            "waved": "localhost:410",
            "entity": "driver",
            "rate": 2,
            "namespace": "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==",
            "publish_uri": "driver/test/gabe/light1",
        }


        # connect to the waved agent
        self.connect()

        # load the wave entity
        self.ent, _ = self.createOrLoadEntity(self.cfg['entity'])
        self.perspective = Perspective(
            entitySecret=EntitySecret(DER=self.ent.SecretDER)
        )
        self.namespace = b64decode(self.cfg['namespace'])
        self.uri = self.cfg['publish_uri']

    def connect(self):
        # connect to wavemq agent
        #from .wavemq_pb2 import *
        #from .wavemq_pb2_grpc import WAVEMQStub
        wavemq_channel = grpc.insecure_channel(self.cfg['wavemq'])
        self.cl = WAVEMQStub(wavemq_channel)

    def begin(self):

        print("DRIVER entity is: ", b64encode(self.ent.hash))

        # call self.setup
        self.setup()

        # start read loop
        @asyncio.coroutine
        async def periodic():
            while True:
                for (data, name) in self.read():
                    po = PayloadObject(
                        schema = "xbosproto/XBOS",
                        content = data.SerializeToString(),
                    )
                    x = self.cl.Publish(PublishParams(
                        perspective=self.perspective,
                        namespace=self.namespace,
                        uri = self.uri+"/"+name,
                        content = [po],
                    ))
                await asyncio.sleep(self.cfg['rate'])
        loop = asyncio.get_event_loop()
        task = loop.create_task(periodic())
        loop.run_until_complete(task)

    def createOrLoadEntity(self, name):
        """
        Check if we have already created an entity (maybe we reset the notebook kernel)
        and load it. Otherwise create a new entity and persist it to disk
        """
        try:
            #from .eapi_pb2 import CreateEntityResponse, Perspective, EntitySecret
            f = open("entity-"+name, "rb")
            entf = pickle.load(f)
            f.close()
            ent = CreateEntityResponse(PublicDER=entf["pub"], SecretDER=entf["sec"], hash=entf["hash"])
            return ent, False
        except (IOError, FileNotFoundError) as e:
            from .wave.eapi_pb2 import CreateEntityParams, PublishEntityParams
            from .wave.eapi_pb2_grpc import WAVEStub
            wave_channel = grpc.insecure_channel(self.cfg['waved'])
            wv = WAVEStub(wave_channel)
            # TODO: what are default params
            ent = wv.CreateEntity(CreateEntityParams())
            if ent.error.code != 0:
                raise Exception(repr(ent.error))
            entf = {"pub":ent.PublicDER, "sec":ent.SecretDER, "hash":ent.hash}
            f = open("entity-"+name, "wb")
            pickle.dump(entf, f)
            f.close()
            resp = wv.PublishEntity(PublishEntityParams(DER=ent.PublicDER))
            if resp.error.code != 0:
                raise Exception(resp.error.message)
            return ent, True
def b64decode(e):
    return base64.b64decode(e, altchars=bytes('-_', 'utf8'))
def b64encode(e):
    return base64.b64encode(e, altchars=bytes('-_', 'utf8'))
