package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gtfierro/hoddb/hod"
)

// determining for a given graph if the edges together actually grant
// the permissions required

// Given an entity, determine if all of the resources it has or uses are
// reachable. Find a path to the namespace entity such that the intersection
// of permissions along that path is sufficient

type traversal struct {
	spec Spec
	hod  *hod.HodDB
	// start is the entity we want to start from
	start string
}

func newTraversal(spec Spec, start string) *traversal {

	t := &traversal{
		start: start,
		spec:  spec,
	}

	dir, err := ioutil.TempDir("", "_brickhod_")
	if err != nil {
		log.Fatal(err)
	}

	cfgStr := fmt.Sprintf(`
database:
    path: %s`, dir)
	cfg, err := hod.ReadConfigFromString(cfgStr)
	if err != nil {
		log.Fatal(err)
	}

	t.hod, err = hod.MakeHodDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	bundle := hod.FileBundle{
		GraphName: "test",
		TTLFile:   "test.ttl",
	}
	if err = t.hod.Load(bundle); err != nil {
		log.Fatal(err)
	}

	return t
}

// returns edges
func (t *traversal) findIncomingEdges(to string) (edges []edge) {
	for _, edge := range t.spec.Edges {
		if edge.To == to {
			edges = append(edges, edge)
		}
	}
	return
}

func (t *traversal) traverse() {

	// 1. make a list of all of the permissions we need
	permissions := t.permissionsRequired(t.start)

	// 2. for each of these permissions, trace the graph back towards
	// the namespace entity, taking the intersetion of the permissions along the way
	for _, perm := range permissions {
		terminates_at_root := false
		fmt.Println("---")
		fmt.Println("try to prove", perm, "for", t.start)
		var edges = newEdgeStack()
		for _, edge := range t.findIncomingEdges(t.start) {
			edges.push(edge)
		}

		for edges.length() > 0 {
			edge := edges.pop()

			if edge.Namespace != perm.Namespace {
				// skip edges with the wrong namespace
				continue
			}

			// intersect the permissions on this edge against
			// the permissions we are trying to build
			// We get
			granted, ok := RestrictBy(perm.Resource, edge.Resource)
			if !ok {
				// in this case, the permission granted on this
				// edge isn't sufficient for what we want to build
				// so we continue on
				fmt.Printf("bad edge %s (restricted to %s)\n", edge, granted)
				continue
			}

			if edge.From == perm.Namespace {
				// found the root authority
				terminates_at_root = true
				break
			}

			// find edges that go into this entity
			newEdges := t.findIncomingEdges(edge.From)
			for _, newEdge := range newEdges {
				edges.push(newEdge)
			}
		}
		if terminates_at_root {
			fmt.Printf("Permission %+v permitted\n", perm)
		} else {
			fmt.Printf("Permission %+v NOT GRANTED\n", perm)
		}
	}

}

// get the set of permissions that the entity requires.
// We get this by querying the graph for all of the resources used
// and owned by all of the processes that use this entity
func (t *traversal) permissionsRequired(entity string) []edge {
	var permissions []edge

	q, err := t.hod.ParseQuery(fmt.Sprintf(`SELECT ?ns ?uri ?ent ?proc ?res WHERE {
		?ent rdfs:label "%s".
		?proc xbos:hasEntity ?ent .
		?proc xbos:usesResource ?res .
		?res xbos:hasNamespace ?ns .
		?res xbos:hasURI ?uri .
	};`, entity), 0)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := t.hod.Select(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Variables)

	nsidx := -1
	uriidx := -1
	for idx, varname := range resp.Variables {
		if varname == "?ns" {
			nsidx = idx
		} else if varname == "?uri" {
			uriidx = idx
		}
	}

	for _, row := range resp.Rows {
		permissions = append(permissions, edge{
			To:          t.start,
			Namespace:   row.Values[nsidx].Value,
			Resource:    row.Values[uriidx].Value,
			Permissions: "subscribe",
			Pset:        "wavemq",
		})
	}

	q, _ = t.hod.ParseQuery(fmt.Sprintf(`SELECT ?ns ?uri ?ent ?proc WHERE {
		?ent rdfs:label "%s".
		?proc xbos:hasEntity ?ent .
		?proc xbos:hasResource ?res .
		?res xbos:hasNamespace ?ns .
		?res xbos:hasURI ?uri .
	};`, entity), 0)
	resp, err = t.hod.Select(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	nsidx = -1
	uriidx = -1
	for idx, varname := range resp.Variables {
		if varname == "?ns" {
			nsidx = idx
		} else if varname == "?uri" {
			uriidx = idx
		}
	}
	for _, row := range resp.Rows {
		permissions = append(permissions, edge{
			To:          t.start,
			Namespace:   row.Values[nsidx].Value,
			Resource:    row.Values[uriidx].Value,
			Permissions: "publish,subscribe",
			Pset:        "wavemq",
		})
	}

	for _, perm := range permissions {
		fmt.Printf("%+v\n", perm)
	}

	return permissions
}
