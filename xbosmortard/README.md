# xbosmortard

This is the client frontend to the `ingester` process which stores timeseries data retrieved through WAVEMQ subscriptions and the `hod` process which stores Brick models.

Default configuration uses InfluxDB and assumes the same database/measurement setup that the `ingester` process uses.

### Server Configuration

Setup WAVEMQ configuration in `mortarconfig.yml`

```yaml
WAVEMQ:
    SiteRouter: "localhost:4516"
    EntityFile: "myentity.ent"
    Namespace: "GyBzLC-FCBTB7cO8SzRCD-t2uI-RykjxgXy5s2t06Ddi4Q=="
    BaseURI: "mortar"
    ServerName: "mortarserver"
```

The service is available at `GyBzLC-FCBTB7cO8SzRCD-t2uI-RykjxgXy5s2t06Ddi4Q== mortar/s.grpcserver/mortarserver/i.grpc`.

You should of course configure the rest of the file to point to whatever db you have (InfluxDB or BTrDB)

### Client Configuration

Install/upgrade the `pyxbos` package

```python
client = MortarClient({
    'id': 'testclientid3',
    'namespace':  "GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==",
    'base_resource': 'mortar/s.grpcserver/mortarserver/i.grpc',
})
```

And now you can use the Mortar API in the same fashion. See the [example file](https://github.com/gtfierro/xboswave/blob/master/python/examples/mortartest.py)
