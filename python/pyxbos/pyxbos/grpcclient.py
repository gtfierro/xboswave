import os
import logging
import random
import random
import time
import threading
import grpc
import pickle
import base64
from pymortar import Result
from pymortar import mortar_pb2 as pymortar
from pymortar import RAW, MEAN, MIN, MAX, COUNT, SUM
from pyxbos.helloworld_pb2 import *
from pyxbos.eapi_pb2 import *
from pyxbos.wavemq_pb2 import *
from pyxbos.wavemq_pb2_grpc import *
from pyxbos.grpcserver_pb2 import *
import asyncio

WaveBuiltinPSET = b"\x1b\x20\x19\x49\x54\xe8\x6e\xeb\x8f\x91\xff\x98\x3a\xcc\x56\xe6\xc8\x4a\xe2\x9a\x90\x7c\xe7\xe7\x63\x8e\x86\x57\xd5\x14\x99\xb1\x88\xa4"
WaveGlobalNamespace = b"\x1b\x20\xcf\x8d\x19\xd7\x9d\x23\x01\x38\x65\xbe\xf7\x57\xce\xa0\x4c\xde\xe5\xef\x4e\xde\xfc\x80\x8d\xd2\x1e\x4e\x00\x5e\x6f\x80\x47\xcc"

WaveBuiltinE2EE = "decrypt"

class GRPCClient:

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

    def call(self, method, argument, returntype):
        arg = google_dot_protobuf_dot_any__pb2.Any()
        arg.Pack(argument)

        queryid = random.randint(1,int(1<<32))

        sub = self.cl.Subscribe(SubscribeParams(
            perspective=self._perspective,
            namespace=self._namespace,
            uri=self._uri+'/signal/response',
            identifier=self._cfg['id'],
            expiry=10,
        ))
        
        call = UnaryCall(
            method=method,
            query_id=queryid,
            payload=arg,
        )

        po = PayloadObject(
            schema = "xbosproto/GRPCServer",
            content = call.SerializeToString(),
        )

        x = self.cl.Publish(PublishParams(
            perspective=self._perspective,
            namespace=self._namespace,
            uri=self._uri + '/slot/call',
            content=[po],
        ))
        for msg in sub:
            if len(msg.error.message) > 0:
                self._log.error("Get actuation message. Error {0}".format(msg.error.message))
                continue
            m = msg.message
            for po in m.tbs.payload:
                resp = UnaryResponse.FromString(po.content)
                actual = returntype()
                resp.payload.Unpack(actual)
                return actual
            break
        sub.cancel()

    def stream(self, method, argument, returntype):
        arg = google_dot_protobuf_dot_any__pb2.Any()
        arg.Pack(argument)

        queryid = random.randint(1,int(1<<32))

        print('subscribing')
        sub = self.cl.Subscribe(SubscribeParams(
            perspective=self._perspective,
            namespace=self._namespace,
            uri=self._uri+'/signal/response',
            identifier=self._cfg['id'],
            expiry=10,
        ))
        
        call = StreamingCall(
            method=method,
            query_id=queryid,
            payload=arg,
        )

        po = PayloadObject(
            schema = "xbosproto/GRPCServer",
            content = call.SerializeToString(),
        )

        print('make call')
        x = self.cl.Publish(PublishParams(
            perspective=self._perspective,
            namespace=self._namespace,
            uri=self._uri + '/slot/stream',
            content=[po],
        ))
        print('wait response')
        for msg in sub:
            if len(msg.error.message) > 0:
                self._log.error("Get actuation message. Error {0}".format(msg.error.message))
                continue
            m = msg.message
            for po in m.tbs.payload:
                resp = StreamingResponse.FromString(po.content)
                if resp.error != "":
                    print("ERROR:", resp.error)
                if not resp.finished:
                    actual = returntype()
                    resp.payload.Unpack(actual)
                    yield actual
                else:
                    return
        sub.cancel()


def b64decode(e):
    return base64.b64decode(e, altchars=bytes('-_', 'utf8'))
def b64encode(e):
    return base64.b64encode(e, altchars=bytes('-_', 'utf8'))

class MortarClient:
    def __init__(self, cfg):
        if 'wavemq' not in cfg:
            cfg['wavemq'] = 'localhost:4516'
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

        self._client = GRPCClient(cfg)

    def qualify(self, required_queries):
        arg = pymortar.QualifyRequest(required=required_queries)
        return self._client.call('Qualify', arg, pymortar.QualifyResponse)
    
    def fetch(self, request):
        res = Result()
        for a in self._client.stream('Fetch', request, pymortar.FetchResponse):
            if a.error != "":
                logging.error(a.error)
                break
            res._add(a)
        res._build()
        return res
