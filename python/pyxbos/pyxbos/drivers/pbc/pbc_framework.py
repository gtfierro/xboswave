from pyxbos.process import XBOSProcess, b64decode, b64encode, schedule, run_loop
from pyxbos.xbos_pb2 import XBOS
from pyxbos.nullabletypes_pb2 import Double
from pyxbos.energise_pb2 import EnergiseMessage, LPBCStatus, LPBCCommand, SPBC, EnergiseError, EnergisePhasorTarget, ChannelStatus, ActuatorCommand
from pyxbos.c37_pb2 import Phasor, PhasorChannel
from datetime import datetime
from functools import partial
from collections import deque
import asyncio
import traceback

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
        self._reference_channels = cfg['reference_channels']
        self.reference_phasors = {k: None for k in self._reference_channels}
        for channel in self._reference_channels:
            upmu_uri = f"upmu/{channel}"
            self._log.info(f"Subscribing to {channel} as reference phasor")
            schedule(self.subscribe_extract(self.namespace, upmu_uri, ".C37DataFrame", self._upmucb, "spbc_reference"))

        self.lpbcs = {}
        schedule(self.subscribe_extract(self.namespace, "lpbc/*", ".EnergiseMessage.LPBCStatus", self._lpbccb, "lpbc_status"))

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
        upmu = resp.uri[5:]
        self.reference_phasors[upmu] = resp.values[-1]['phasorChannels'][0]['data']
        for chan in resp.values[-1]['scalarChannels']:
            if chan['channelName'] == 'FREQ':
                key = 'freq'
            elif chan['channelName'] == 'DFREQ':
                key = 'dfreq'
            else:
                continue
            for upmu_name, phasor in self.reference_phasors.items():
                for idx, d in enumerate(phasor):
                    assert chan['data'][idx]['time'] == d['time']
                    d[key] = chan['data'][idx].get('value')


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
            'pSaturated': True,
            # true if Q is saturated
            'qSaturated': True,
            # if pSaturated is True, expect the p max value
            'pMax': 1.4,
            # if qSaturated is True, expect the q max value
            'qMax': 11.4,
            # true if LPBc is doing control
            'do_control': True,
        }
        """
        if len(resp.values) == 0:
            return
        statuses = resp.values[-1]
        if statuses is None:
            return
        timestamp = statuses['time']
        for status in statuses['statuses']:
            if status['nodeID'] not in self.lpbcs:
                self.lpbcs[status['nodeID']] = {}
            if 'pSaturated' not in status:
                status['pSaturated'] = False
            if 'qSaturated' not in status:
                status['qSaturated'] = False
            if 'pMax' in status:
                status['pMax'] = status['pMax']['value']
            if 'qMax' in status:
                status['qMax'] = status['qMax']['value']
            self.lpbcs[status['nodeID']][status['channelName']] = status
        #self.lpbcs[resp.uri] = resp

    async def broadcast_target(self, nodeid, channels, vmags, vangs, kvbases=None, kvabases=None):
        """
        Publishes SPBC V and delta for a particular node

        Args:
            nodeid (str): the name of the node we are publishing the target to
            channels (list of str): list of channel names for the node we are announcing targets to
            vmag (list of float): the 'V' target to be set for each channel
            vang (list of float): the 'delta' target to be set for each channel
            kvbases (list of float or None): the KV base for each channel
            kvabases (list of float or None): the KVA base for each channel
        """
        self._log.info(f"SPBC announcing channels {channels}, vmag {vmags}, vang {vangs} to node {nodeid}")
        # wrap value in nullable Double if provided

        targets = []
        for idx, channel in enumerate(channels):
            kvbase = Double(value=kvbases[idx]) if kvbases is not None else None
            kvabase = Double(value=kvabases[idx]) if kvabases is not None else None
            targets.append(
                EnergisePhasorTarget(
                    nodeID=nodeid,
                    channelName=channels[idx],
                    angle=vangs[idx],
                    magnitude=vmags[idx],
                    kvbase=kvbase,
                    KVAbase=kvabase,
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
        if 'local_channels' not in cfg:
            raise ConfigMissingError('local_channels')
        if 'reference_channels' not in cfg:
            raise ConfigMissingError('reference_channels')
        if 'name' not in cfg:
            raise ConfigMissingError('name')
        if 'spbc' not in cfg:
            raise ConfigMissingError('spbc')
        if 'rate' not in cfg:
            raise ConfigMissingError('rate')

        # locks
        self._reference_phasor_lock = asyncio.Lock()
        self._local_phasor_lock = asyncio.Lock()

        self.local_channels = cfg['local_channels']
        self.reference_channels = cfg['reference_channels']

        self.name = cfg['name']
        self.spbc = cfg['spbc']
        self._rate = int(cfg['rate'])

        # local buffers for phasor data
        # key: channel name, value: phasor data
        self.local_phasor_data = {}
        self.reference_phasor_data = {}
        self.last_spbc_command = None
        self.control_on = False

        for local_channel in self.local_channels:
            self.local_phasor_data[local_channel] = []
            cb = partial(self._local_upmucb, local_channel)
            schedule(self.subscribe_extract(self.namespace, f"upmu/{local_channel}", ".C37DataFrame", cb, f"local_channel_{local_channel}"))
        for reference_channel in self.reference_channels:
            self.reference_phasor_data[reference_channel] = []
            cb = partial(self._reference_upmucb, reference_channel)
            schedule(self.subscribe_extract(self.namespace, f"upmu/{reference_channel}", ".C37DataFrame", cb, f"reference_phasor_{reference_channel}"))

        # TODO: listen to SPBC
        print(f"spbc/{self.spbc}/node/{self.name}")
        schedule(self.subscribe_extract(self.namespace, f"spbc/{self.spbc}/node/{self.name}", ".EnergiseMessage.SPBC", self._spbccb, "spbc_sub"))

        schedule(self.call_periodic(self._rate, self._trigger, runfirst=False))

        self._log.info(f"initialized LPBC: {cfg}")

    def _local_upmucb(self, channel, resp):
        """Stores the most recent local upmu reading"""
        frame = resp.values[-1]['phasorChannels'][0]
        #self._log.info(f"got {len(frame['data'])} values on local")
        if channel not in self.local_phasor_data:
            self.local_phasor_data[channel] = []
        data = frame['data']
        for chan in resp.values[-1]['scalarChannels']:
            if chan['channelName'] == 'FREQ':
                key = 'freq'
            elif chan['channelName'] == 'DFREQ':
                key = 'dfreq'
            else:
                continue
            for idx, d in enumerate(data):
                assert chan['data'][idx]['time'] == d['time']
                d[key] = chan['data'][idx].get('value')
        self.local_phasor_data[channel].extend(data)

    def _reference_upmucb(self, channel, resp):
        """Stores the most recent reference upmu reading"""
        if len(resp.values) == 0:
            self._log.error("no content in UPMU message")
            return
        if len(resp.values[-1]['phasorChannels']) == 0:
            self._log.error("no phasor channels in UPMU message")
            return
        frame = resp.values[-1]['phasorChannels'][0]
        if channel not in self.reference_phasor_data:
            self.reference_phasor_data[channel] = []
        data = frame['data']
        for chan in resp.values[-1]['scalarChannels']:
            if chan['channelName'] == 'FREQ':
                key = 'freq'
            elif chan['channelName'] == 'DFREQ':
                key = 'dfreq'
            else:
                continue
            for idx, d in enumerate(data):
                assert chan['data'][idx]['time'] == d['time']
                d[key] = chan['data'][idx].get('value')
        self.reference_phasor_data[channel].extend(data)

    def _spbccb(self, resp):
        """Stores the most recent SPBC command"""
        if len(resp.values) == 0:
            self._log.error("no content in SPBC message")
            return
        resp = resp.values[-1]
        resp['phasor_targets'] = resp.pop('phasorTargets')
        self.last_spbc_command = resp

    def _received_local_phasor_data(self):
        return sum(map(len, self.local_phasor_data.values())) > 0
    def _received_reference_phasor_data(self):
        return sum(map(len, self.reference_phasor_data.values())) > 0

    async def _trigger(self):
        try:
            if not self._received_local_phasor_data():
                self._log.warning(f"LPBC {self.name} has not received a local upmu reading")
                #return
            if not self._received_reference_phasor_data():
                self._log.warning(f"LPBC {self.name} has not received a reference upmu reading")
                #return
            if self.last_spbc_command is None:
                self._log.warning(f"LPBC {self.name} has not received an SPBC command")
                #return
            #spbc_cmd = self.last_spbc_command.values[-1] # most recent
            #targets = spbc_cmd['phasor_targets']
            async with self._local_phasor_lock:
                local_phasors = [self.local_phasor_data.pop(local_channel) for local_channel in self.local_channels]
            async with self._reference_phasor_lock:
                reference_phasors = [self.reference_phasor_data.pop(reference_channel) for reference_channel in self.reference_channels]

            phasor_targets = self.last_spbc_command

            # rebuild buffers
            async with self._local_phasor_lock:
                for local_channel in self.local_channels:
                    self.local_phasor_data[local_channel] = []

            async with self._reference_phasor_lock:
                for reference_channel in self.reference_channels:
                    self.reference_phasor_data[reference_channel] = []

            await self.do_trigger(local_phasors, reference_phasors, phasor_targets)
        except Exception as e:
            self._log.error(f"Error in processing trigger: {traceback.format_exc()}")

    async def do_trigger(self, local_phasors, reference_phasors, phasor_targets):
        self._log.info(f"""LPBC {self.name} received call at {datetime.now()}:
    Local phasor has {len(local_phasors)} channels
    Reference phasor has {len(reference_phasors)} channels
    SPBC targets: {phasor_targets}
        """)
        status = self.step(local_phasors, reference_phasors, phasor_targets)
        if status is None:
            return

        for required in ['p_max','q_max','phases','phasor_errors','p_saturated','q_saturated']:
            if required not in status:
                raise Exception(f"Need {required} key in status dictionary")

        statuses = []
        for idx, phase_name in enumerate(status.pop('phases')):
            p_max = Double(value=status['p_max'][idx]) if status['p_max'][idx] else None
            q_max = Double(value=status['q_max'][idx]) if status['q_max'][idx] else None
            channel_status = ChannelStatus(
                nodeID=self.name,
                channelName=phase_name,
                phasor_errors=Phasor(
                    magnitude=status['phasor_errors']['V'][idx],
                    angle=status['phasor_errors']['delta'][idx],
                ),
                p_saturated=status['p_saturated'][idx],
                q_saturated=status['q_saturated'][idx],
                p_max=p_max,
                q_max=q_max,
            )
            statuses.append(channel_status)

        msg = XBOS(
                EnergiseMessage=EnergiseMessage(
                    LPBCStatus=LPBCStatus(
                        time=int(datetime.utcnow().timestamp()*1e9),
                        statuses=statuses,
                    )
                ))
        await self.publish(self.namespace, f"lpbc/{self.name}", msg)

    def log_actuation(self, actuation):
        """
        Publish an EnergiseMessage.ActuatorCommand message containing the parameters that were
        sent to and received from the inverters, etc.
        The 'actuation' argument has the following structure, and is expected to be a dictionary:

            {
                "phases": ["a","b","c"],
                "P_cmd": [10.1, 20.2, 30.3],
                "Q_cmd": [10.9, 20.9, 30.9],
                "P_act": [.1, .2, .3],
                "Q_act": [.1, .2, .3],
                "P_PV": [11.1,22.2,33.3],
                "Batt_cmd": [99.1, 99.2, 99.3],
                "pf_ctrl": [8.7,6.5,5.4]
            }

        All of the values are lists of floats and should have the same length as the 'phases' key.
        """
        components = ['phases','P_cmd','Q_cmd','P_act','Q_act','P_PV','Batt_cmd','pf_ctrl']
        for required in components:
            if required not in actuation:
                raise Exception(f"Need {required} key in dictionary")
        num_components = len(actuation['phases'])
        for component in components:
            if len(actuation[component]) != num_components:
                raise Exception(f"Field {component} needs to be of length {num_components} to match number of phases")

        msg = XBOS(
            EnergiseMessage=EnergiseMessage(
                ActuatorCommand=ActuatorCommand(
                    time=int(datetime.utcnow().timestamp()*1e9),
                    nodeID=self.name,
                    phases=actuation["phases"],
                    P_cmd=actuation["P_cmd"],
                    Q_cmd=actuation["Q_cmd"],
                    P_act=actuation["P_act"],
                    Q_act=actuation["Q_act"],
                    P_PV=actuation["P_PV"],
                    Batt_cmd=actuation["Batt_cmd"],
                    pf_ctrl=actuation["pf_ctrl"],
                )
            )
        )
        schedule(self.publish(self.namespace, f"lpbc/{self.name}/actuation", msg))
