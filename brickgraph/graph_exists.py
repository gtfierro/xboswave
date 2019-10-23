from rdflib import Graph, Literal, BNode, Namespace, RDF, RDFS, OWL
from rdflib.namespace import XSD
import owlrl
import eapi_pb2, eapi_pb2_grpc
from grpc import insecure_channel
import logging
import time

WAVEMQPSET = bytes(b"\x1b\x20\x14\x33\x74\xb3\x2f\xd2\x74\x39\x54\xfe\x47\x86\xf6\xcf\x86\xd4\x03\x72\x0f\x5e\xc4\x42\x36\xb6\x58\xc2\x6a\x1e\x68\x0f\x6e\x01")
XBOS = Namespace("https://xbos.io/schema/tmp/XBOS#")
DEP = Namespace("https://xbos.io/schema/tmp/mydeployment#")
A = RDF.type

class GraphChecker:
    def __init__(self, cfg=None):
        """
        Config takes the following keys:

        waved (default localhost:410): address of waved agent
        entity (default $WAVE_DEFAULT_ENTITY): filepath to entity file of this process
        expiry (default 2 minutes): time we can remain disconnected from wavemq before undelivered
              messages in our subscriptions are dropped
        """


        self._log = logging.getLogger(__name__)
        if cfg is None:
            cfg = {}

        self._log.info("Reading config {0}".format(str(cfg)))

        # check defaults
        if 'waved' not in cfg:
            cfg['waved'] = 'localhost:410'
        if 'entity' not in cfg:
            if 'WAVE_DEFAULT_ENTITY' in os.environ:
                cfg['entity'] = os.environ['WAVE_DEFAULT_ENTITY']

        self._cfg = cfg

        # connect to the wavemq agent
        self._log.info("Connecting to waved agent at {0}".format(cfg['waved']))
        self._connect()
        self._log.info("Connected to waved")

        # load the wave entity
        self._log.info("Loading wave entity {0}".format(cfg['entity']))
        self._ent = open(self._cfg['entity'],'rb').read()
        self._perspective = eapi_pb2.Perspective(
            entitySecret=eapi_pb2.EntitySecret(DER=self._ent),
        )

        self.sync()

        self.G = Graph()
        self.G.bind('rdf', RDF)
        self.G.bind('rdfs', RDFS)
        self.G.bind('owl', OWL)
        self.G.bind('xbos', XBOS)
        self.G.bind('dep', DEP)
        self.G.parse(cfg['graph'], format='ttl')

    def _connect(self):
        # connect to wavemq agent
        wavemq_channel = insecure_channel(self._cfg['waved'])
        self._cl = eapi_pb2_grpc.WAVEStub(wavemq_channel)

    def sync(self):
        resp = self._cl.ResyncPerspectiveGraph(eapi_pb2.ResyncPerspectiveGraphParams(
            perspective=self._perspective,
        ))
        while True:
            r = self._cl.SyncStatus(eapi_pb2.SyncParams(
                perspective=self._perspective,
            ))
            print(f"Syncing: {r.completedSyncs}/{r.totalSyncRequests}")
            if r.totalSyncRequests == r.completedSyncs: break
            time.sleep(1)


    def generate_proofs_for(self, query, permissions):
        res = self.G.query(query)
        proof_requests = {}
        for row in res:
            (ent, name, ns, uri) = row
            # resolve ent name into a hash
            print(name.lower(), ns, uri)
            result = self._cl.ResolveName(eapi_pb2.ResolveNameParams(
                perspective=self._perspective,
                name=name.lower(),
            ))
            if result.error.message != '':
                raise Exception(result.error.message)
            subjectHash = result.entity.hash


            result = self._cl.ResolveName(eapi_pb2.ResolveNameParams(
                perspective=self._perspective,
                name=ns.lower(),
            ))
            if result.error.message != '':
                raise Exception(result.error.message)
            namespace = result.entity.hash


            statement = eapi_pb2.RTreePolicyStatement(
                permissionSet = WAVEMQPSET,
                permissions = permissions,
                resource = uri,
            )
            cli = f"wv rtprove --subject {name}.ent wavemq:{','.join(permissions)}@{ns}/{uri}"
            if ent not in proof_requests:
                proof_requests[name] = []
            proof_requests[name].append(cli)
        
        from mako.template import Template

        shell_template = Template("""
#!/bin/bash
echo "Checking proofs for ${entity}.ent"
% for line in lines:
${line}
% endfor""")

        for entity, proofs in proof_requests.items():
            print(f"Proofs for {entity}.ent")

            with open(f"check_{entity}.sh", "w") as f:
                f.write(shell_template.render(entity=entity, lines=proofs))



checker = GraphChecker({'entity': 'attgraph/gabe.ent', 'graph': 'test.ttl'})

sub_resources = """SELECT ?ent ?name ?namespace ?uri WHERE {
    ?proc xbos:hasEntity ?ent .
    ?ent rdf:type xbos:Entity .
    ?ent rdfs:label ?name .
    ?proc xbos:usesResource ?res .
    ?res xbos:hasNamespace ?namespace .
    ?res xbos:hasURI ?uri
}"""

checker.generate_proofs_for(sub_resources, ['subscribe'])

pub_resources = """SELECT ?ent ?name ?namespace ?uri WHERE {
    ?proc xbos:hasEntity ?ent .
    ?ent rdf:type xbos:Entity .
    ?ent rdfs:label ?name .
    ?proc xbos:hasResource ?res .
    ?res xbos:hasNamespace ?namespace .
    ?res xbos:hasURI ?uri
}"""

checker.generate_proofs_for(pub_resources, ['publish', 'subscribe'])