# Install Daemons

XBOS requires the `waved` and `wavemq` daemons to be installed and running.
We use the `wv` binary to communicate with the `waved` daemon and configure entities and permissions.

1. Fill out `env.sh` with any different configuration you want
2. Run `run.sh` to download and install the binaries. Also sets up basic config files and starts
   processes using systemd
    - **NOTE**: currently just works on amd64, systemd-based Linux distros
