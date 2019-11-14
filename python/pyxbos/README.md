# PyXBOS


## Driver Design Principles

A driver is a persistent process that communicates over WAVEMQ and is identified by a WAVE entity.
The purpose of a driver is to present a uniform interface to 1 or more underlying devices or services (for example: thermostats, BMS, lights, weather service).
The interface regularly *publishes* the current state of these devices, and can optionally receive commands for those devices.

Published state and received commands are formatted as protocol buffer messages, defined in [/proto](https://github.com/gtfierro/xboswave/tree/master/proto).
These protocol buffer messages define the fields that can be reported by devices, as well as the Brick classes for those fields; associating Brick classes with messages and fields lets us use Brick's definitions for clarification on whata field or message means.

## Anatomy of a Driver

I will implement a driver for the Philips Hue lighting system using the Python driver framework as a running example.

All Python drivers should subclass the `pyxbos.driver.Driver` class, which provides helpful methods + infrastructure for writing a driver.

### The `setup` Method

The `setup` method is called once when a driver is executed.
The `cfg` argument to the `setup` method is a Python dictionary containing the key-value pairs necessary for a driver's operation.
Aside from a small number of required keys (documented in the "Configuring a Driver" section below), there are no restrictions on what keys you can include in this dictionary, so knock yourself out!

For now, lets just consider the keys we need for the Philips Hue driver.
We will need the IP address of the "bridge" device which is capable of talking to the individual light bulbs in our deployment, as well as the username/password for the bridge's API.

```python
# sample configuration
cfg = {
    "bridge_address": "192.168.0.101",
    "username": "admin",
    "password": "default"
}
```

Here is the beginning of the driver implementation.
For now, assume that when the driver is executed, the `cfg` parameter of the `setup` function will be equal to the dictionary above.

The `setup` method takes in the configuration parameters and prepares the driver for execution.
This typically involves one or more of the following:

- setting up local API clients for the devices the driver needs to talk to
- initializing internal data structures needed for the driver's operation
- performing *discovery* to automatically establish the set of devices the driver is representing:
    - this is expected for drivers such as those talking to building management systems (BMS). A driver may create a BACnet client and run a "who-is" query to discover available devices. The driver can then loop through this device list to initialize internal data structures to "remember" the discovered devices.

```python
from pyxbos.driver import *
import logging
import time
from phue import Bridge

class HueDriver(Driver):
    def setup(self, cfg):
        # Initialize the Bridge client
        self.bridge = Bridge(cfg['bridge_address'],
                             username=cfg['username'],
                             password=cfg['password'])

        # Connect the client and prepare API access
        self.bridge.connect()
        self.bridge.get_api()

        # Initialize internal structure to track what lights we discover.
        # This allows us to look up lights by name.
        self.lights = {}
        # this loops through the lights that the Bridge discovered
        for l in self.bridge.lights:
            name = l.name.replace(' ','_') # remove spaces from name
            self.lights[name] = l
```

### The `read` Method

The `read` method is called at a regular rate (given by the driver configuration; see "Configuring a Driver") and is responsible for reporting the current state of the devices represented by the driver.
The `read` method takes a single optional parameter (`requestid`) that is used for actuation; we will ignore this for now.
The `read` method does not have a return argument.

The nominal behavior of a `read` method is to loop through the devices known by the driver (these are likely stored in an internal data structure established in the `setup` method), get their current state, and use the built-in `self.report(resource, message)` method to publish

`self.report` takes two arguments: a `resource`, which is a string indicating the WAVEMQ resource (topic) that the message is published on, and `msg` which is the `xbospb.XBOS` instance to be published.
The resource is hierarchical, where levels are delimited by the `/` character.
These resources do not have to capture spatial or functional context such as location or subsystem because this is adequately described in a Brick model.

A common resource breakdown is `<device type>/<device name>`; e.g. for a BMS this might be `ahu/ahu-1` or `vav/vav2-17`.

## Configuring a Driver

Drivers expect the following configuration keys

Key Name | Definition | Example
---------|------------|--------
`wavemq` | GRPC endpoint of the local WAVEMQ router | `127.0.0.1:4516` (default)
`namespace` | the base64-encoded public key of the WAVEMQ namespace that the driver operates on |  `GyDX55sFnbr9yCB-mPyXsy4kAUPUY8ftpWX62s6UcnvfIQ==`
`base_resource` | the common prefix for all WAVEMQ resources used in this driver | `mybuildings/bms`
`entity` | path to the entity file for this driver | `bms_driver.ent` (defaults to `$WAVE_DEFAULT_ENTITY`)
`id`  | the name of the driver | `library-bms-driver`
`rate` | rate at which `self.read` is called; polling rate of the driver | `3` (every 3 seconds)
