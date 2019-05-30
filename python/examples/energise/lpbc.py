from pyxbos.process import run_loop
from pyxbos.drivers import pbc
import logging
logging.basicConfig(level="INFO", format='%(asctime)s - %(name)s - %(message)s')
import random

class democontroller(pbc.LPBCProcess):
    """
    To implement a LPBC, subclass pbc.LPBCProcess
    and implement the step() method as documented below
    """
    def __init__(self, cfg):
        super().__init__(cfg)
        self.measured_p = 1
        self.measured_q = 1
        self.vmag = None
        self.vang = None
        self.saturated = False

    def step(self, timestamp, vmag, vang, p_target, q_target):
        """
        Step is called every 'rate' seconds with the most recent Vmag and Vang from the upmu
        and the latest P and Q targets given by the SPBC.

        It runs its control loop to determine the actuation, performs it is 'self.control_on' is True
        and returns the status
        """
        self.vmag, self.vang = vmag, vang

        # do measurements
        self.measured_p = random.randint(0,100)
        self.measured_q = random.randint(0,100)

        p_diff = self.measured_p - p_target
        q_diff = self.measured_q - q_target

        print(f'controller called. P diff: {p_diff}, Q diff: {q_diff}')

        if self.control_on:
            print("DO CONTROL HERE")
        

        # return error message (default to empty string), p, q and boolean saturated value
        return ("error message", self.measured_p, self.measured_q, self.saturated)

cfg = {
        'namespace': "GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg==",
        'name': 'lpbc1', # name of lpbc
        'upmu': 'L1', # name + other info for uPMU
        'rate': 2, # number of seconds between calls to 'step'
        }
lpbc1 = democontroller(cfg)
run_loop()
