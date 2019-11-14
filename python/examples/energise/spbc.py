from pyxbos.process import run_loop, schedule
from pyxbos.drivers import pbc
import logging
import random
logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')

class myspbc(pbc.SPBCProcess):
    def __init__(self, cfg):
        super().__init__(cfg)
        self.p = 10
        self.q = 20

        # recalculate every 30 seconds
        schedule(self.call_periodic(30, self.compute_and_announce))

    async def compute_and_announce(self):
        # do expensive compute to get new P and Q
        for lpbc, status in self.lpbcs.items():
            print('LPBC status:', lpbc,':', status)
        self.p = 100 * random.random()
        self.q = 100 * random.random()
        await self.broadcast_target(self.p, self.q)

cfg = {
        'namespace': "GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg=="
}
spbc_instance = myspbc(cfg)
run_loop()
