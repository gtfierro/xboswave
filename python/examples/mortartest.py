from pyxbos import MortarClient
import pymortar
import time

client = MortarClient({
    'id': 'testclientid3',
    'namespace':  "GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==",
    'base_resource': 'mortar/s.grpcserver/mortar/i.grpc',
})
# client.qualify
resp = client.qualify([
    "SELECT ?zone WHERE { ?zone rdf:type brick:Electric_Meter };",
    "SELECT ?zone WHERE { ?zone rdf:type brick:Temperature_Sensor };"
])

req = pymortar.FetchRequest(
    sites=resp.sites,
    views=[
        pymortar.View(
            name="test1",
            definition="SELECT ?vav WHERE { ?vav rdf:type/rdfs:subClassOf* brick:Temperature_Sensor };",
        ),
        pymortar.View(
            name="meter",
            definition="SELECT ?meter WHERE { ?meter rdf:type/rdfs:subClassOf* brick:Electric_Meter };",
        ),
    ],
    dataFrames=[
        pymortar.DataFrame(
            name="meter_data",
            aggregation=pymortar.MEAN,
            window="5m",
            uuids=["b8166746-ba1c-5207-8c52-74e4700e4467"],
            #timeseries=[
            #    pymortar.Timeseries(
            #        view="meter",
            #        dataVars=["?meter"],
            #    )
            #]
        )
    ],
    time=pymortar.TimeParams(
        start="2019-01-01T00:00:00Z",
        end="2019-04-01T00:00:00Z",
    )
)
s = time.time()
res = client.fetch(req)
e = time.time()
print("took {0}".format(e-s))
print(res)
