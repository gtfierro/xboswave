from rdflib import Graph, Literal, BNode, Namespace, RDF, RDFS, OWL
from rdflib.namespace import XSD
import owlrl
from copy import deepcopy

XBOS = Namespace("https://xbos.io/schema/tmp/XBOS#")
DEP = Namespace("https://xbos.io/schema/tmp/mydeployment#")

def findResourceByURI(G, uri):
    q = f"SELECT ?res WHERE {{ ?res xbos:hasURI \"{uri}\" }}"
    l = list(G.query(q))
    return [row[0] for row in l]

G = Graph()
G.bind('rdf', RDF)
G.bind('rdfs', RDFS)
G.bind('owl', OWL)
G.bind('xbos', XBOS)
G.bind('dep', DEP)
A = RDF.type

# declare  classes
G.add( (XBOS.Entity, A, OWL.Class) )
G.add( (XBOS.Deployment, A, OWL.Class) )
G.add( (XBOS.Process, A, OWL.Class) )
G.add( (XBOS.Namespace, A, OWL.Class) )
G.add( (XBOS.Resource, A, OWL.Class) )
G.add( (XBOS.Controller, RDFS.subClassOf, XBOS.Process) )
G.add( (XBOS.uPMU, RDFS.subClassOf, OWL.Process) )
G.add( (XBOS.Supervisory_Phasor_Based_Controller, RDFS.subClassOf, XBOS.Controller) )
G.add( (XBOS.Local_Phasor_Based_Controller, RDFS.subClassOf, XBOS.Controller) )

# properties:
G.add( (XBOS.hasEntity, A, OWL.ObjectProperty) )
G.add( (XBOS.hasEntity, RDFS.domain, XBOS.Process) )
G.add( (XBOS.hasEntity, RDFS.range, XBOS.Entity) )

# hasResource is a property of a Process. Having a resource means
# that the property has the ability to publish on it
G.add( (XBOS.hasResource, A, OWL.ObjectProperty) )
G.add( (XBOS.hasResource, RDFS.domain, XBOS.Process) )
G.add( (XBOS.hasResource, RDFS.range, XBOS.Resource) )

# usesResource denotes which resources a process must
# subscribe to
G.add( (XBOS.usesResource, A, OWL.ObjectProperty) )
G.add( (XBOS.usesResource, RDFS.domain, XBOS.Process) )
G.add( (XBOS.usesResource, RDFS.range, XBOS.Resource) )

# TODO: different property for "clientOf", which implies
# a set of properties for publishing AND subscribing

# Resources may have multiple representations. One is as URI, given as
# a string literal. We may have URLs in the future
G.add( (XBOS.hasURI, A, OWL.DatatypeProperty) )
G.add( (XBOS.hasURI, RDFS.domain, XBOS.Resource) )
G.add( (XBOS.hasURI, RDFS.range, XSD.string) )

G.add( (XBOS.hasNamespace, A, OWL.DatatypeProperty) )
G.add( (XBOS.hasNamespace, RDFS.domain, XBOS.Resource) )
G.add( (XBOS.hasNamespace, RDFS.range, XSD.string) )

# instantiate processes
G.add( (DEP.dep, A, XBOS.Deployment) )
G.add( (DEP.dep, RDFS.label, Literal("Sample Energise Deployment")) )

upmus = ['uPMU_0','uPMU_123', 'uPMU_123P', 'uPMU_4']
for upmu in upmus:
    G.add( (DEP[upmu], A, XBOS.uPMU) )
    G.add( (DEP[upmu], RDFS.label, Literal(upmu)) )
    for phase in ['L1','L2','L3','C1','C2','C3']:
        res = DEP[f'{upmu}_{phase}_resource']
        G.add( (DEP[upmu], XBOS.hasResource, res) )
        G.add( (res, XBOS.hasURI, Literal(f"upmu/{upmu}/{phase}") ) )
        G.add( (res, XBOS.hasNamespace, Literal("energise") ) )

# lpbc usesResources
# static binding for now
lpbc_uses = {
    DEP.lpbc_675: [
        'upmu/uPMU_123P/L1',
        'upmu/uPMU_123P/C1',
    ],
    DEP.lpbc_790: [
        'upmu/uPMU_123/L1',
        'upmu/uPMU_123/C1',
    ],
}

for lpbc in ['lpbc_675','lpbc_790']:
    G.add( (DEP[lpbc], A, XBOS.Local_Phasor_Based_Controller) )
    G.add( (DEP[lpbc], RDFS.label, Literal(lpbc) ) )
    res = f"lpbc/{lpbc}/status"
    G.add( (DEP[f"{lpbc}_status"], A, XBOS.Resource) )
    G.add( (DEP[f"{lpbc}_status"], XBOS.hasURI, Literal(res)) )
    G.add( (DEP[f"{lpbc}_status"], XBOS.hasNamespace, Literal("energise")) )
    G.add( (DEP[lpbc], XBOS.hasResource, DEP[f"{lpbc}_status"]) )

    for uri in lpbc_uses[DEP[lpbc]]:
        for res in findResourceByURI(G, uri):
            G.add( (DEP[lpbc], XBOS.usesResource, res) )


## Define SPBCs
for spbc in ['spbc_0']:
    G.add( (DEP[spbc], A, XBOS.Supervisory_Phasor_Based_Controller) )
    G.add( (DEP[spbc], RDFS.label, Literal(spbc) ) )

    # bind SPBC to use LPBC status resources
    # SPBC needs all LPBCs for now
    pred = """SELECT ?lpbc ?res WHERE { ?lpbc a xbos:Local_Phasor_Based_Controller . ?lpbc xbos:hasResource ?res }"""
    for row in G.query(pred):
        print(row[1])
        G.add( (DEP[spbc], XBOS.usesResource, row[1]) )
        #TODO: add target topic: spbc HAS resource the target
        #G.add( (DEP[spbc], XBOS.hasResource, 

    # which uPMUs does the SPBC need?

    # static binding; the boring way
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_L1_resource) )
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_L2_resource) )
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_L3_resource) )
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_C1_resource) )
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_C2_resource) )
    G.add( (DEP[spbc], XBOS.usesResource, DEP.uPMU_0_C3_resource) )

    # a more interesting way: dynamic binding based on a query.
    # still have to know the PMU by name though
    upmu_0_resources = G.query(f"SELECT ?res WHERE {{ <{DEP.uPMU_0}> xbos:hasResource ?res }}")
    for res in upmu_0_resources:
        G.add( (DEP[spbc], XBOS.usesResource, res[0]) )

    # TODO: even more interesting: find uPMU through a query

# apply reasoner
Q = deepcopy(G)
owlrl.DeductiveClosure(owlrl.OWLRL_Semantics).expand(Q)

# infer the entities
all_processes = Q.query("""SELECT ?proc ?label WHERE { ?proc rdf:type xbos:Process . ?proc rdfs:label ?label}""")
for proc, label in all_processes:
    ent = DEP[f"{label}_entity"]
    G.add( (ent, A, XBOS.Entity) )
    G.add( (ent, RDFS.label, Literal(label)) )
    G.add( (proc, XBOS.hasEntity, ent) )

with open('test.ttl','wb') as f:
    f.write(G.serialize(format='turtle'))
