from pyxbos.process import XBOSProcess, b64decode, b64encode, schedule, run_loop
from pyxbos.xbos_pb2 import XBOS
from pyxbos.energise_pb2 import EnergiseMessage, LPBCStatus, LPBCCommand, SPBC, EnergiseError, EnergisePhasorTarget, Double
from pyxbos.c37_pb2 import Phasor
from datetime import datetime
from collections import deque

class ConfigMissingError(Exception): pass

class SPBCProcess(XBOSProcess):
    """
    Wrapper process for supervisory phasor-based control in Python.

    An SPBC subscribes to a set of LPBCs and receives status mesages from them.
    These status messages consist of an error quantity and saturation state. 
    In the current implementation, the SPBC subscribes to all LPBCs that it has
    permission to see.

    At regular intervals, the SPBC publishes a V + delta target for *each* node,
    represented by an LPBC.

    TODO: how does the SPBC get the most recent C37 reference phasor reading?
    """

    def __init__(self, cfg):
        super().__init__(cfg)
        if 'namespace' not in cfg:
            raise ConfigMissingError('namespace')
        if 'name' not in cfg:
            raise ConfigMissingError('name')
        if 'reference_channels' not in cfg:
            raise ConfigMissingError('reference_channels')
        self.namespace = b64decode(cfg['namespace'])
        self._log.info(f"initialized SPBC: {cfg}")
        self.name = cfg['name']

        # TODO: document this configuration option
        self._reference_channels = reference_channels
        for channel in reference_channels:
            upmu_uri = f"upmu/{channel}"
            self._log.info(f"Subscribing to {channel} as reference phasor")
            schedule(self.subscribe_extract(self.namespace, upmu_uri, ".C37DataFrame", self._upmucb))

        self.lpbcs = {}
        schedule(self.subscribe_extract(self.namespace, "lpbc/*", ".EnergiseMessage.LPBCStatus", self._lpbccb))

    def _upmucb(self, c37_frame):
        """
        """
        # TODO: handle upmu

        pass

    def _lpbccb(self, resp):
        """
        Caches the last message heard from each LPBC
        """
        self.lpbcs[resp.uri] = resp

    async def broadcast_target(self, nodeid, vmag, vang, kvbase=None):
        """
        Publishes SPBC V and delta for a particular node

        Args:
            nodeid (str): the name of the node we are publishing the target to
            vmag (float): the 'V' target to be set
            vang (float): the 'delta' target to be set
            kvbase (float): the KV base
        """
        self._log.info(f"SPBC announcing vmag {vmag}, vang {vang} to node {nodeid}")
        # wrap value in nullable Double if provided
        kvbase = Double(value=kvbase) if kvbase else None
        await self.publish(self.namespace, f"spbc/{self.name}/node/{nodeid}", XBOS(
            EnergiseMessage = EnergiseMessage(
            SPBC=SPBC(
                time=int(datetime.utcnow().timestamp()*1e9),
                phasor_target=EnergisePhasorTarget(nodeID=nodeid,magnitude=vmag,angle=vang,kvbase=kvbase)
                )
            )))


class LPBCProcess(XBOSProcess):
    """
    Wrapper process for local phasor-based control in Python.
    Requires a uPMU
    """
    def __init__(self, cfg):
        super().__init__(cfg)
        if 'namespace' not in cfg:
            raise ConfigMissingError('namespace')
        self.namespace = b64decode(cfg['namespace'])
        if 'local_upmu' not in cfg:
            raise ConfigMissingError('local_upmu')
        if 'reference_upmu' not in cfg:
            raise ConfigMissingError('reference_upmu')
        if 'name' not in cfg:
            raise ConfigMissingError('name')
        if 'spbc' not in cfg:
            raise ConfigMissingError('spbc')
        if 'rate' not in cfg:
            raise ConfigMissingError('rate')

        self.local_upmu = cfg['local_upmu']
        self.reference_upmu = cfg['reference_upmu']

        self.name = cfg['name']
        self.spbc = cfg['spbc']
        self._rate = int(cfg['rate'])
        self.last_local_upmu_reading = None
        self.last_reference_upmu_reading = None
        self.last_spbc_command = None
        self.control_on = False

        schedule(self.subscribe_extract(self.namespace, f"upmu/{self.local_upmu}", ".C37DataFrame", self._local_upmucb))
        schedule(self.subscribe_extract(self.namespace, f"upmu/{self.reference_upmu}", ".C37DataFrame", self._reference_upmucb))
        #schedule(self.subscribe_extract(self.namespace, f"spbc/{self.spbc}/*", ".EnergiseMessage.SPBC", self._spbccb))
        schedule(self.subscribe_extract(self.namespace, "spbc/*", ".EnergiseMessage.SPBC", self._spbccb))
	
        schedule(self.call_periodic(self._rate, self._trigger, runfirst=False))

        self._log.info(f"initialized LPBC: {cfg}")

    def _local_upmucb(self, resp):
        """Stores the most recent local upmu reading"""
        self.last_local_upmu_reading = resp
    def _reference_upmucb(self, resp):
        """Stores the most recent reference upmu reading"""
        self.last_reference_upmu_reading = resp
    def _spbccb(self, resp):
        """Stores the most recent SPBC command"""
        self.last_spbc_command = resp

    async def _trigger(self):
        try:
            if self.last_local_upmu_reading is None or len(self.last_local_upmu_reading.values) == 0:
                self._log.warning(f"LPBC {self.name} has not received a local upmu reading")
                return
            if self.last_reference_upmu_reading is None or len(self.last_reference_upmu_reading.values) == 0:
                self._log.warning(f"LPBC {self.name} has not received a reference upmu reading")
                return
            if self.last_spbc_command is None or len(self.last_spbc_command.values) == 0:
                self._log.warning(f"LPBC {self.name} has not received an SPBC command")
                return
            local_c37_frame = self.last_local_upmu_reading.values[-1] # most recent
            reference_c37_frame = self.last_reference_upmu_reading.values[-1] # most recent
            spbc_cmd = self.last_spbc_command.values[-1] # most recent
            targets = spbc_cmd['phasorTarget']
            await self.do_trigger(local_c37_frame, reference_c37_frame, 1,0)#targets['P'], targets['Q'])
        except IndexError:
            return # no upmu readings

    async def do_trigger(self, local_c37_frame, reference_c37_frame, p_target, q_target):
        self._log.info(f"""LPBC {self.name} received call at {datetime.now()}:
    Local C37 data frame from {local_c37_frame['stationName']}
    Reference C37 data frame from {reference_c37_frame['stationName']}
    SPBC targets: P: {p_target} Q: {q_target}
    Enable control?: {self.control_on}
        """)
        error, mp, mq, sat = self.step(local_c37_frame, reference_c37_frame, p_target, q_target)

        msg = XBOS(
                EnergiseMessage=EnergiseMessage(
                    LPBCStatus=LPBCStatus(
                        time=int(datetime.utcnow().timestamp()*1e9),
                        phasors=Phasor(P=mp, Q=mq),
                        saturated=sat,
                        do_control=self.control_on,
                    )
                ))
        await self.publish(self.namespace, f"lpbc/{self.name}", msg)
