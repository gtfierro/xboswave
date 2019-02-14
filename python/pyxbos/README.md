# PyXBOS


## Driver Framework

The Python driver framework provides a simple skeleton for maintaining an interface between a
device/service and WAVEMQ.

Drivers are instantiated with a configuration. This configuration contains *at least*:

- `wavemq` (default `localhost:4516`): the GRPC endpoint of the local WAVEMQ site router
- `entity` (defaults to `$WAVE_DEFAULT_ENTITY`): the WAVE entity representing the driver instance.
  The driver's ability to publish and interact is dependent on this entity having the proper
  permissions to do so.
- `namespace`: the public key of the namespace the driver is interacting on
- `base_resource`: the resource prefix of where the driver interacts within the namespace

Booting a driver invokes the following elements in order:

- validate configuration
- connect to local WAVEMQ site router
- load entity file
- call the driver's `setup` method (driver-specific instantiation)
- subscribe to actuation topic
- start "read loop"
- start "write loop"

The "read loop" is called periodically (in default configuration). The rate at which it is called
is configured using the `rate` configuration variable (integer seconds between calls). The read loop
simply calls the driver's `read()` function at that regular interval.

---

Thoughts on setting this up.
Because we have devices at URIs, we should bind together the read and write methods for a device.
This way its easier to trigger a report on a corresponding read and whatnot.
