# XBOS Demo Setup

These scripts + environment set up a minimal, single-machine installation of XBOS. Not recommended for production deployment.

The installation consists of the following services:

- a `waved` daemon to fulfill the WAVE API for the `wv` utility
- a `wavemq` message bus daemon
- a `ingester` daemon to archive data published on the message bus
    - the `ingester` daemon is backed by an InfluxDB instance
- a system monitor driver periodically publishing CPU/Mem/Disk stats to the message bus

The installation operates on a single WAVE namespace and creates an administrative entity with full permissions to help manage the namespace.
Each of the services and drivers above has its own individual WAVE entity and is granted appropriate permissions.


```
#Informal tree of delegations

namespace.ent
    |
    +--> route on * to router.ent
    |
    +--> pub, sub on * to admin.ent
                            |
                            +--> pub, sub on drivers/systemmonitor/* to system_monitor_driver.ent
                            |
                            +--> sub on * to ingester.ent
```

Each process runs in its own Docker container.

The `wv` invocations to create the grants and entities can be found in `run.sh`.

### Setting Up

Requires:
- git
- docker
    - This is easier if you [add your user to the `docker` group to avoid using sudo](https://docs.docker.com/install/linux/linux-postinstall/)
- [`wv`](https://github.com/gtfierro/xboswave/tree/master/demo-setup/bin) utility

```bash
git clone github.com/gtfierro/xboswave
cd xboswave/demo-setup
source environment.sh
./run.sh
```

### Tearing Down

This will stop and remove all Docker containers.
The entities and InfluxDB data dir will remain on disk

```
# sourcing environment.sh not necessary
./stop.sh
```
