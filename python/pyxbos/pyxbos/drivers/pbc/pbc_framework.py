from pyxbos.process import XBOSProcess, b64decode, b64encode, schedule, run_loop
from pyxbos.xbos_pb2 import XBOS
from pyxbos.energise_pb2 import EnergiseMessage, LPBCStatus, LPBCCommand, SPBC, EnergiseError, EnergisePhasorTarget, Double
from pyxbos.c37_pb2 import Phasor, PhasorChannel
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

        # reference channels are URIs for the uPMU channels the SPBC
        # subscribes to. The SPBC framework maintains self.reference_phasors
        # to contain the most recent phasor measurements for each channel
        self.reference_phasors = {k: None for k in reference_channels}
        self._reference_channels = reference_channels
        for channel in reference_channels:
            upmu_uri = f"upmu/{channel}"
            self._log.info(f"Subscribing to {channel} as reference phasor")
            schedule(self.subscribe_extract(self.namespace, upmu_uri, ".C37DataFrame", self._upmucb))

        self.lpbcs = {}
        schedule(self.subscribe_extract(self.namespace, "lpbc/*", ".EnergiseMessage.LPBCStatus", self._lpbccb))

    def _upmucb(self, resp):
        """
        Called on every received C37 frame from the reference upmu channels.
        Stores the most recent phasors received:

        self.reference_phasors = {
            'flexlab1/L1': [
                {
                    "time": "1559231114799996800",
                    "angle": 193.30149788923268,
                    "magnitude": 0.038565948605537415
                },
                {
                    "time": "1559231114899996400",
                    "angle": 195.50249902851263,
                    "magnitude": 0.042079225182533264
                }
                ... etc
            ],
            'flexlab/L2': [
                {
                    "time": "1559231114799996800",
                    "angle": 220.30149788923268,
                    "magnitude": 10.038565948605537415
                },
                {
                    "time": "1559231114899996400",
                    "angle": 220.50249902851263,
                    "magnitude": 10.042079225182533264
                }
            ]
        }
        """
        upmu = resp.uri.lstrip('upmu/')
        self.reference_phasors[upmu] = resp.values[-1]['phasorChannels']['data']

    def _lpbccb(self, resp):
        """
        Caches the last message heard from each LPBC

        Each LPBC status looks like:
        {
            # local time of LPBC
            'time': 1559231114799996800,
            # phasor errors of LPBC
            'phasor_errors': {
                'angle': 1.12132,
                'magnitude': 31.12093090,
                # ... and/or ...
                'P': 1.12132,
                'Q': 31.12093090,
            },
            # true if P is saturated
            'p_saturated': True,
            # true if Q is saturated
            'q_saturated': True,
            # if p_saturated is True, expect the p max value
            'p_max': {'value': 1.4},
            # if q_saturated is True, expect the q max value
            'q_max': {'value': 11.4},
            # true if LPBc is doing control
            'do_control': True,
        }
        """
        self.lpbcs[resp.uri] = resp

    async def broadcast_target(self, nodeid, channels, vmags, vangs, kvbases=None):
        """
        Publishes SPBC V and delta for a particular node

        Args:
            nodeid (str): the name of the node we are publishing the target to
            channels (list of str): list of channel names for the node we are announcing targets to
            vmag (list of float): the 'V' target to be set for each channel
            vang (list of float): the 'delta' target to be set for each channel
            kvbase (list of float or None): the KV base for each channel
        """
        self._log.info(f"SPBC announcing channels {channels}, vmag {vmags}, vang {vangs} to node {nodeid}")
        # wrap value in nullable Double if provided

        targets = []
        for idx, channel in enumerate(channels):
            kvbase = Double(value=kvbases[idx]) if kvbases else None
            targets.append(
                EnergisePhasorTarget(
                    nodeID=nodeid,
                    channelName=channels[idx],
                    angle=vangs[idx],
                    magnitude=vmags[idx], 
                    kvbase=kvbase,
                )
            )

        await self.publish(self.namespace, f"spbc/{self.name}/node/{nodeid}", XBOS(
            EnergiseMessage = EnergiseMessage(
            SPBC=SPBC(
                time=int(datetime.utcnow().timestamp()*1e9),
                phasor_targets=targets
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
