from pyxbos.driver import *
import logging
import os
#; print os.uname()[1] 
import time
import psutil

class PsutilDriver(Driver):
    def setup(self, cfg):
        psutil.cpu_percent() # need to throw away first value
        self.hostname = os.uname()[1]
        self._log.info("# CPUS {0}".format(psutil.cpu_count()))

    def read(self, requestid=None):
        msg = xbos_pb2.XBOS(
            BasicServerStatus = system_monitor_pb2.BasicServerStatus(
                time = int(time.time()*1e9),
                hostname =              self.hostname,
                cpu_load =              [types.Double(value=x) for x in psutil.cpu_percent(interval=1, percpu=True)],
                phys_mem_available =    types.Int64(value=psutil.virtual_memory().available),
                disk_usage =            types.Double(value=psutil.disk_usage(path='/').percent),
                disk_available =        types.Double(value=psutil.disk_usage(path='/').free),
            )
        )
        self.report(self.hostname, msg)

if __name__ == '__main__':
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

    cfg = {
        'wavemq': 'localhost:4516',
        'namespace': 'GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==',
        'base_resource': 'test/system',
        'entity': 'system.ent',
        'id': 'system',
        'rate': 10,
    }
    e = PsutilDriver(cfg)
    e.begin()
