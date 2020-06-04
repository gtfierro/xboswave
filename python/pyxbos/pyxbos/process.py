"""
Simple wrapper for a control process
"""

import asyncio
import traceback
import xxhash
import toml
from datetime import datetime
import logging
import base64
import os
import jq
import uuid
from aiogrpc import insecure_channel
from google.protobuf.json_format import MessageToDict
from pyxbos.exceptions import *
from . import xbos_pb2
from .eapi_pb2 import *
from .wavemq_pb2 import *
from .wavemq_pb2_grpc import *

class XBOSProcess:
    """
    Base class wrapping interaction with WAVE, WAVEMQ
    """
    def __init__(self, cfg=None):
        """
        Config takes the following keys:

        wavemq (default localhost:4516): address of wavemq site router
        waved (default localhost:410): address of waved agent
        entity (default $WAVE_DEFAULT_ENTITY): filepath to entity file of this process
        expiry (default 2 minutes): time we can remain disconnected from wavemq before undelivered
              messages in our subscriptions are dropped
        """
        self._log = logging.getLogger(__name__)
        if cfg is None:
            cfg = {}

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
        if 'expiry' not in cfg:
            cfg['expiry'] = 10 # default to 10 second expiry

        self._cfg = cfg

        # connect to the wavemq agent
        self._log.info("Connecting to wavemq agent at {0}".format(cfg['wavemq']))
        self._connect()
        self._log.info("Connected to wavemq")

        # load the wave entity
        self._log.info("Loading wave entity {0}".format(cfg['entity']))
        self._ent = open(self._cfg['entity'],'rb').read()
        self._perspective = Perspective(
            entitySecret=EntitySecret(DER=self._ent),
        )
        self._subscription_expiry = cfg['expiry']

        # associate (ns, uri) => identifier
        self._subscription_ids = {}

    def _get_identifier(self, namespace, resource, extra=None):
        """
        Assumes namespace is b64 encoded
        """
        h = xxhash.xxh32()
        h.update(self._ent)
        h.update(namespace)
        h.update(resource)
        if extra:
            h.update(extra)
        identifier = b64encode(h.digest())
        return identifier

#        key = (namespace, resource)
#        while self._subscription_ids.get(key)
#        if self._subscription_ids.get(key) is None:
#            self._subscription_ids[key] = identifier
#        else:
#            h.
#        return self._subscription_ids[key]

    def _connect(self):
        # connect to wavemq agent
        wavemq_channel = insecure_channel(self._cfg['wavemq'])
        self._cl = WAVEMQStub(wavemq_channel)

    async def subscribe_msg(self, namespace, resource, callback, name=None):
        """
        callback takes mqpb_pb.SubscriptionMessage as an argument
        """
        ns = ensure_b64decode(namespace)
        async for msg in self._cl.Subscribe(SubscribeParams(
                perspective=self._perspective,
                namespace=ns,
                uri=resource,
                identifier=self._get_identifier(namespace, resource, extra=name),
                expiry=self._subscription_expiry,
            )
            ):
            try:
                callback(msg)
            except:
                self._log.error(f"Error in processing callback: {traceback.format_exc()}")

    async def subscribe_extract(self, namespace, resource, path, callback, name=None):
        """
        extracts the submessage at the given path
        callback returns a Response object:
        - uri
        - namespace (base64)
        - sent timestamp
        - submessage (list)
        """
        def cb(msg):
            uri = msg.message.tbs.uri
            namespace = ensure_b64encode(msg.message.tbs.namespace)
            sent_timestamp = msg.message.timestamps
            if len(sent_timestamp) == 0:
                sent_timestamp = datetime.now()
            else:
                sent_timestamp = datetime.utcfromtimestamp(sent_timestamp[0])
            values = []
            for po in msg.message.tbs.payload:
                x = xbos_pb2.XBOS.FromString(po.content)
                x = MessageToDict(x)
                values.append(jq.jq(path).transform(x))
            callback(Response(namespace, uri, sent_timestamp, values))
        await self.subscribe_msg(namespace, resource, cb, name=name)

    async def publish(self, namespace, resource, *msgs):
        """publishes msgs in list as payload objects"""
        pos = []
        for msg in msgs:
            pos.append(PayloadObject(
                schema = "xbosproto/XBOS",
                content = msg.SerializeToString(),
                ))
        namespace = ensure_b64decode(namespace)
        try:
            x = await self._cl.Publish(PublishParams(
                perspective=self._perspective,
                namespace=namespace,
                uri = resource,
                content = pos,
                ))
            if not x:
                self._log.error("Error publishing: {0}".format(x))
        except Exception as e:
            self._log.error("Error publishing: {0}".format(e))


    async def call_periodic(self, seconds, cb, *args, runfirst=True):
        """
        Run asynchronous function every n seconds. If runfirst is true, we run
        the function once before starting the timer
        """
        if runfirst:
            try:
                schedule(cb(*args))
            except:
                self._log.error(f"Error in processing callback: {traceback.format_exc()}")
        while True:
            await asyncio.sleep(seconds)
            try:
                schedule(cb(*args))
            except:
                self._log.error(f"Error in processing callback: {traceback.format_exc()}")

class Response():
    def __init__(self, ns, uri, ts, values):
        self.ns = ns
        self.uri = uri
        self.ts = ts
        self.values = values

    def __repr__(self):
        return f"Response<({self.ns}, {self.uri}, {self.ts}, {self.values}"

def b64decode(e):
    return base64.b64decode(e, altchars=bytes('-_', 'utf8'))
def ensure_b64decode(e):
    return e if isinstance(e, bytes) else b64decode(e)

def ensure_b64encode(e):
    return b64encode(e) if isinstance(e, bytes) else e
def b64encode(e):
    return base64.b64encode(e, altchars=bytes('-_', 'utf8'))

def config_from_file(filename):
    """
    Returns XBOS process configuration from a file
    """
    return toml.load(open(filename))

def schedule(f):
    """
    Runs task asynchronously (subscribe, publish)
    """
    fun = asyncio.ensure_future(f)
    def handle_exception(f=None):
        exc = fun.exception()
        if exc is not None:
            logging.error(exc)
    fun.add_done_callback(handle_exception)

def run_loop():
    """
    Block forever, running async tasks
    """
    loop = asyncio.get_event_loop()
    loop.run_forever()
    loop.close()
