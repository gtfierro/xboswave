from pyxbos.process import XBOSProcess, b64decode, b64encode, schedule, run_loop
from pyxbos.xbos_pb2 import XBOS
from pyxbos.energise_pb2 import EnergiseMessage, LPBCStatus, LPBCCommand, SPBC, EnergiseError
from pyxbos.c37_pb2 import Phasor
from datetime import datetime
from collections import deque

class SPBCProcess(XBOSProcess):
    def __init__(self, cfg):
        super().__init__(cfg)
        if 'namespace' not in cfg:
            raise ConfigMissingError('namespace')
        elif 'name' not in cfg:
            raise ConfigMissingError('name')
        self.namespace = b64decode(cfg['namespace'])
        self._log.info(f"initialized SPBC: {cfg}")
        self.name = cfg['name']

        self.lpbcs = {}
        schedule(self.subscribe_extract(self.namespace, "lpbc/*", ".EnergiseMessage.LPBCStatus", self._lpbccb))

    def _lpbccb(self, resp):
        self.lpbcs[resp.uri] = resp

    async def broadcast_target(self, p, q):
        self._log.info(f"SPBC announcing p {p} q {q}")
        await self.publish(self.namespace, "spbc/{self.name}", XBOS(
            EnergiseMessage = EnergiseMessage(
            SPBC=SPBC(
                time=int(datetime.utcnow().timestamp()*1e9),
                phasor_target=Phasor(P=p, Q=q)
                )
            )))


class LPBCProcess(XBOSProcess):
    def __init__(self, cfg):
        super().__init__(cfg)
        if 'namespace' not in cfg:
            raise ConfigMissingError('namespace')
        self.namespace = b64decode(cfg['namespace'])
        if 'upmu' not in cfg:
            raise ConfigMissingError('upmu')
        if 'name' not in cfg:
            raise ConfigMissingError('name')
        if 'spbc' not in cfg:
            raise ConfigMissingError('spbc')
        if 'rate' not in cfg:
            raise ConfigMissingError('rate')
        self.upmu = cfg['upmu']
        self.name = cfg['name']
        self.spbc = cfg['spbc']
        self._rate = int(cfg['rate'])
        self.last_upmu_reading = None
        self.last_spbc_command = None
        self.control_on = False

        schedule(self.subscribe_extract(self.namespace, f"upmu/{self.upmu}", ".C37DataFrame", self._upmucb))
        schedule(self.subscribe_extract(self.namespace, "spbc/{self.spbc}", ".EnergiseMessage.SPBC", self._spbccb))
	
        schedule(self.call_periodic(self._rate, self._trigger, runfirst=False))

        self._log.info(f"initialized LPBC: {cfg}")

    def _upmucb(self, resp):
        self.last_upmu_reading = resp
    def _spbccb(self, resp):
        self.last_spbc_command = resp

    async def _trigger(self):
        try:
            if self.last_upmu_reading is None or len(self.last_upmu_reading.values) == 0:
                self._log.warning(f"LPBC {self.name} has not received a upmu reading")
                return
            if self.last_spbc_command is None or len(self.last_spbc_command.values) == 0:
                self._log.warning(f"LPBC {self.name} has not received an SPBC command")
                return
            c37_frame = self.last_upmu_reading.values[-1] # most recent
            spbc_cmd = self.last_spbc_command.values[-1] # most recent
            targets = spbc_cmd['phasorTarget']
            await self.do_trigger(c37_frame, targets['P'], targets['Q'])
        except IndexError:
            return # no upmu readings

    async def do_trigger(self, c37_frame, p_target, q_target):
        self._log.info(f"""LPBC {self.name} received call at {datetime.now()}:
    C37 data frame: {c37_frame}
    SPBC targets: P: {p_target} Q: {q_target}
    Enable control?: {self.control_on}
        """)
        error, mp, mq, sat = self.step(c37_frame, p_target, q_target)

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
