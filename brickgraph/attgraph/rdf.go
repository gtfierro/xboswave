package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	rdf "github.com/gtfierro/hoddb/turtle"
	"os"
)

var _nsuuid = uuid.MustParse("77b5594a-d019-11e9-8e19-bb3cebfd2292")

// take a specification and transform it into an RDF graph
// start from the graph engine so we can use actual hashes

func fmtTriple(t rdf.Triple) string {
	return fmt.Sprintf("%s\t%s\t%s .\n", fmtURI(t.Subject), fmtURI(t.Predicate), fmtURI(t.Object))
}

func fmtURI(u rdf.URI) string {
	if u.Namespace != "" {
		return fmt.Sprintf("%s:%s", u.Namespace, u.Value)
	}
	return fmt.Sprintf("\"%s\"", u.Value)
}

func (eng *GraphEngine) ToRDF() error {

	rdf_namespaces := map[string]string{
		"owl":     "http://www.w3.org/2002/07/owl#",
		"rdf":     "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
		"rdfs":    "http://www.w3.org/2000/01/rdf-schema#",
		"brick":   "https://brickschema.org/schema/1.1.0/Brick#",
		"wave":    "https://xbos.io/ontologies/0.0.1/WAVE#",
		"mygraph": "https://xbos.io/ontologies/tmp/mygraph#",
	}

	var triples = make(map[rdf.Triple]struct{})
	var hashes = make(map[string]string)

	addTriple := func(t rdf.Triple) {
		triples[t] = struct{}{}
	}

	for entityname := range eng.entities {
		hashes[entityname] = base64.URLEncoding.WithPadding(' ').EncodeToString(eng.hashes[entityname])
	}

	// declare all entities, named by hash
	for entityname := range eng.entities {
		// <hash> a wave:entity
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + hashes[entityname]),
			Predicate: rdf.ParseURI("rdf:type"),
			Object:    rdf.ParseURI("wave:Entity"),
		})

		// <hash> wave:name "name"
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + hashes[entityname]),
			Predicate: rdf.ParseURI("wave:name"),
			Object:    rdf.URI{Value: entityname},
		})
	}

	// declare edges
	for _, edge := range eng.spec.Edges {
		edgeuuid := uuid.New().String()
		var data []byte
		data = append(data, eng.hashes[edge.Namespace]...)
		data = append(data, eng.hashes[edge.Pset]...)
		data = append(data, []byte(edge.Permissions)...)
		uriuuid := uuid.NewSHA1(_nsuuid, data).String()
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("rdf:type"),
			Object:    rdf.ParseURI("wave:Attestation"),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("wave:Attester"),
			Object:    rdf.ParseURI("mygraph:" + hashes[edge.From]),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("wave:subject"),
			Object:    rdf.ParseURI("mygraph:" + hashes[edge.To]),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("wave:hasURI"),
			Object:    rdf.ParseURI("mygraph:" + uriuuid),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("wave:permissions"),
			Object:    rdf.URI{Value: edge.Permissions},
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + edgeuuid),
			Predicate: rdf.ParseURI("wave:pset"),
			Object:    rdf.ParseURI("mygraph:" + hashes[edge.Pset]),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + uriuuid),
			Predicate: rdf.ParseURI("rdf:type"),
			Object:    rdf.ParseURI("wave:URI"),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + uriuuid),
			Predicate: rdf.ParseURI("wave:namespace"),
			Object:    rdf.ParseURI("mygraph:" + hashes[edge.Namespace]),
		})
		addTriple(rdf.Triple{
			Subject:   rdf.ParseURI("mygraph:" + uriuuid),
			Predicate: rdf.ParseURI("wave:resource"),
			Object:    rdf.URI{Value: edge.Resource},
		})
	}

	// dump to turtle
	f, err := os.Create("dump.ttl")
	if err != nil {
		return err
	}
	for abb, ns := range rdf_namespaces {
		fmt.Fprintf(f, "@prefix %s: <%s> .\n", abb, ns)
	}
	for t := range triples {
		fmt.Fprintf(f, fmtTriple(t))
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
