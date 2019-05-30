from pyxbos.process import run_loop
from pyxbos.drivers import pbc
import logging
logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
import random

class democontroller(pbc.LPBCProcess):
    def __init__(self, cfg):
        super().__init__(cfg)
        self.measured_p = 1
        self.measured_q = 1
        self.vmag = None
        self.vang = None
        self.saturated = False

    def step(self, timestamp, vmag, vang, p_target, q_target):
        self.vmag, self.vang = vmag, vang

        # do measurements
        self.measured_p = random.randint(0,100)
        self.measured_q = random.randint(0,100)

        p_diff = self.measured_p - p_target
        q_diff = self.measured_q - q_target

        print(f'controller called. P diff: {p_diff}, Q diff: {q_diff}')

        if self.control_on:
            print("DO CONTROL HERE")
        

        return ("error message", self.measured_p, self.measured_q, self.saturated)

cfg = {
        'namespace': "GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg==",
        'name': 'lpbc1',
        'upmu': 'L1',
        }
lpbc1 = democontroller(cfg)
run_loop()
