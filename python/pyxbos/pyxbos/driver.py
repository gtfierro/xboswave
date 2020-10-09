import os
import logging
import time
import threading
import grpc
import pickle
import base64

from . import nullabletypes_pb2 as types
from .eapi_pb2 import *
from .wavemq_pb2 import *
from .wavemq_pb2_grpc import *
from . import xbos_pb2
from . import iot_pb2
from . import system_monitor_pb2
from pyxbos.exceptions import *

import asyncio

class Driver:
    """Base class encapsulating driver report functionality"""
    def __init__(self, cfg):
        self._log = logging.getLogger(__name__)

        self._log.info("Reading config {0}".format(str(cfg)))

        # check defaults
        if 'wavemq' not in cfg:
            cfg['wavemq'] = 'localhost:4516'
        if 'waved' not in cfg:
            cfg['waved'] = 'localhost:410'
        if 'entity' not in cfg:
            if 'WAVE_DEFAULT_ENTITY' in os.environ:
                cfg['entity'] = os.environ['WAVE_DEFAULT_ENTITY']
            else:
                raise ConfigMissingError('entity', extra="And no WAVE_DEFAULT_ENTITY in environment")
        if 'id' not in cfg:
            raise ConfigMissingError('id')
        if 'namespace' not in cfg:
            raise ConfigMissingError('namespace')
        if 'base_resource' not in cfg:
            raise ConfigMissingError('base_resource')

        self._cfg = cfg

        # connect to the wavemq agent
        self._log.info("Connecting to wavemq agent at {0}".format(cfg['wavemq']))
        self.connect()
        self._log.info("Connected to wavemq")

        # load the wave entity
        self._log.info("Loading wave entity {0}".format(cfg['entity']))
        self._ent = open(self._cfg['entity'],'rb').read()
        self._perspective = Perspective(
            entitySecret=EntitySecret(DER=self._ent),
        )
        self._namespace = b64decode(self._cfg['namespace'])
        self._uri = self._cfg['base_resource']

    def connect(self):
        # connect to wavemq agent
        wavemq_channel = grpc.insecure_channel(self._cfg['wavemq'])
        self.cl = WAVEMQStub(wavemq_channel)

    def begin(self):
        # call self.setup
        self._log.info("Run driver setup")
        self.setup(self._cfg)

        # subscribe to the write uri
        writeuri = self._cfg['base_resource']+'/write/*'
        self._log.info("Subscribe to write URI {0}".format(writeuri))

        sub = self.cl.Subscribe(SubscribeParams(
            perspective=self._perspective,
            namespace=self._namespace,
            uri=writeuri,
            identifier=self._cfg['id'],
            expiry=120,
        ))

        loop = asyncio.get_event_loop()

        async def _doread(requestid=None):
            self.read(requestid=requestid)

        async def readloop():
            while True:
                await _doread()
                await asyncio.sleep(self._cfg['rate'])

        # this runs in a thread
        def writeloop():
            # create an event loop because we're in a new thread
            loop = asyncio.new_event_loop()
            self._log.info("write loop")
            for msg in sub:
                if len(msg.error.message) > 0:
                    self._log.error("Get actuation message. Error {0}".format(msg.error.message))
                    continue
                m = msg.message
                now = int(time.time()*1e9)
                # seconds
                since = (now - m.timestamps[-1]) / 1.e9
                #print('timestamps', m.timestamps, 'since', since)
                #print('drops', m.drops)
                #print('resource', m.tbs.uri)
                #print('pos', len(m.tbs.payload))
                for po in m.tbs.payload:
                    print('po', po.schema, len(po.content))
                    x = xbos_pb2.XBOS.FromString(po.content)
                    try:
                        self.write(m.tbs.uri, since, x)
                    except Exception as e:
                        print('error write', e)

        # start thread
        t = threading.Thread(target=writeloop)
        t.start()


        asyncio.ensure_future(readloop())
        try:
            loop.run_forever()
        finally:
            loop.close()

    def report(self, resource, msg):
        po = PayloadObject(
            schema = "xbosproto/XBOS",
            content = msg.SerializeToString(),
        )
        self._log.info("Publishing on %s", self._uri+"/"+resource)
        try:
            x = self.cl.Publish(PublishParams(
                perspective=self._perspective,
                namespace=self._namespace,
                uri = self._uri+"/"+resource,
                content = [po],
            ))
            if not x:
                self._log.error("Error reading: {0}".format(x))
                print('x>',x)
        except Exception as e:
            self._log.error("Error reading: {0}".format(e))

def b64decode(e):
    return base64.b64decode(e, altchars=bytes('-_', 'utf8'))
def b64encode(e):
    return base64.b64encode(e, altchars=bytes('-_', 'utf8'))
