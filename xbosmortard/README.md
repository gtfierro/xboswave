# xbosmortard

This is the client frontend to the `ingester` process which stores timeseries data retrieved through WAVEMQ subscriptions and the `hod` process which stores Brick models.

Default configuration uses InfluxDB and assumes the same database/measurement setup that the `ingester` process uses.
