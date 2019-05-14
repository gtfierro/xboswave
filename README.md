# XBOS(2) 

[![Build Status](https://travis-ci.org/gtfierro/xboswave.svg?branch=master)](https://travis-ci.org/gtfierro/xboswave)

## Looking for a Quick Start?

Take a gander at the [demo installation](https://github.com/gtfierro/xboswave/tree/master/demo-setup)

## Ingester

Requirements:
1. go [install](https://golang.org/doc/install)
2. bzr `sudo apt-get install bzr`

### BTrDB Setup

How do get started:

1. Install [BtrDB dev machine](https://docs.smartgrid.store/development-environment.html) OR [InfluxDB](https://docs.influxdata.com/influxdb/v1.7/introduction/)
2. Install [WAVEMQ](https://github.com/immesys/wavemq):
3. Create entity (docs coming)
    - read [wave docs](https://github.com/immesys/wave) for now
4. Give perms
    - read [wave docs](https://github.com/immesys/wave) for now
5. Build and run ingester:
    ```
    cd ingester
    make run
    ```
