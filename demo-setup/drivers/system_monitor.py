from pyxbos import *
from pyxbos.drivers.system_monitor import systemmonitor
import logging
import os

if __name__ == '__main__':
    logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

    cfg = {
        'wavemq': os.environ.get('WAVEMQ_SITE_ROUTER'),
        'namespace': os.environ.get('XBOS_NAMESPACE'),
        'base_resource': os.environ.get('BASE_RESOURCE'),
        'entity':os.environ.get('WAVE_DEFAULT_ENTITY'),
        'id': 'system',
        'rate': 10,
    }
    print(cfg)
    e = systemmonitor.PsutilDriver(cfg)
    e.begin()
