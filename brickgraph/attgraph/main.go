package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/immesys/wave/consts"
	"github.com/immesys/wave/eapi"
	"github.com/immesys/wave/eapi/pb"

	"github.com/BurntSushi/toml"
	"google.golang.org/grpc"
)

type GraphEngine struct {
	wave     pb.WAVEClient
	policies map[string]policy
	// map entity name -> contents
	entities map[string][]byte
	hashes   map[string][]byte
	spec     Spec
}

func GraphEngineFromSpecFile(filename string) *GraphEngine {
	var spec Spec
	if _, err := toml.DecodeFile(filename, &spec); err != nil {
		log.Fatal(err)
	}
	return GraphEngineFromSpec(spec)
}

func GraphEngineFromSpec(spec Spec) *GraphEngine {
	eng := &GraphEngine{
		policies: make(map[string]policy),
		entities: make(map[string][]byte),
		hashes:   make(map[string][]byte),
		spec:     spec,
		wave:     getConn("localhost:410"),
	}

	// add common pset aliases
	eng.hashes["wavemq"] = []byte("\x1b\x20\x14\x33\x74\xb3\x2f\xd2\x74\x39\x54\xfe\x47\x86\xf6\xcf\x86\xd4\x03\x72\x0f\x5e\xc4\x42\x36\xb6\x58\xc2\x6a\x1e\x68\x0f\x6e\x01")
	rv, _ := base64.URLEncoding.DecodeString(consts.WaveBuiltinPSET)
	eng.hashes["wave"] = rv
	eng.hashes["jedi"] = consts.JEDIBuiltinPSETByteArray

	// do some pre-processing

	// push parse 'edge' into namespace and resource
	for idx, att := range eng.spec.Edges {
		if att.Edge != "" {
			parts := strings.Split(att.Edge, ":")
			if len(parts) != 2 {
				log.Fatal(xerrors.Errorf("Edge %s invalid. Needs ns and resource", att.Edge))
			}
			att.Namespace = parts[0]
			att.Resource = parts[1]
			parts = strings.SplitN(att.Permissions, ":", -1)
			if len(parts) != 2 {
				log.Fatal(xerrors.Errorf("Permission %s invalid. Needs pset and perm list", att.Permissions))
			}
			att.Pset = parts[0]
			att.Permissions = parts[1]
			eng.spec.Edges[idx] = att
		}
	}

	for idx, pol := range eng.spec.Policies {
		if pol.Edge != "" {
			parts := strings.SplitN(pol.Edge, ":", -1)
			if len(parts) != 2 {
				log.Fatal(xerrors.Errorf("Policy %s edge %s invalid. Needs ns and resource", pol.Name, pol.Edge))
			}
			pol.Namespace = parts[0]
			pol.Resource = parts[1]
			parts = strings.SplitN(pol.Permissions, ":", -1)
			if len(parts) != 2 {
				log.Fatal(xerrors.Errorf("Policy %s permission %s invalid. Needs pset and perm list", pol.Name, pol.Permissions))
			}
			pol.Pset = parts[0]
			pol.Permissions = parts[1]
			eng.spec.Policies[idx] = pol
			fmt.Printf("Policy: %+v\n", pol)
		}
		eng.policies[pol.Name] = pol
	}

	// push policies into the actual edges
	for idx, att := range eng.spec.Edges {
		if att.Policy != "" {
			pol, found := eng.policies[att.Policy]
			if !found {
				log.Fatal(xerrors.Errorf("Policy %s not found. Make sure it is declared", att.Policy))
			}
			fmt.Println(pol)
			att.Namespace = pol.Namespace
			att.Resource = pol.Resource
			att.Permissions = pol.Permissions
			att.Pset = pol.Pset
			eng.spec.Edges[idx] = att
		}
		// push expiries too
		att := eng.spec.Edges[idx]
		if att.Expiry == nil {
			att.Expiry = &eng.spec.Graph.GrantExpiry
			eng.spec.Edges[idx] = att
		}
	}

	fmt.Printf("%+v\n", eng.spec)
	return eng
}

func (eng *GraphEngine) getEntityHash(content []byte) ([]byte, error) {
	resp, err := eng.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, xerrors.Errorf("could not inspect file: %s\n", resp.Error.Message)
	}
	if resp.Entity == nil {
		return nil, xerrors.Errorf("file was not an entity")
	}
	return resp.Entity.Hash, nil
}

func (eng *GraphEngine) ValidateConnectivity() (bool, []string) {

	// keep track of which entities are visited
	visited_entities := make(map[string]bool)
	for _, edge := range eng.spec.Edges {
		visited_entities[edge.From] = false
		visited_entities[edge.To] = false
	}

	// convenience method to get outgoing edges from an entity
	// in the context of this graph
	find1Hop := func(from string) ([]string, []string, []string) {
		var entities []string
		var permissions []string
		var uris []string
		for _, edge := range eng.spec.Edges {
			if edge.From == from {
				entities = append(entities, edge.To)
				permissions = append(entities, edge.Permissions)
				uris = append(entities, edge.Resource)
			}
		}
		return entities, permissions, uris
	}

	// do a BFS to check if the graph is fully connected across all of the namespaces
	for _, ns := range eng.spec.Graph.Namespaces {
		var nodes = newStack()
		var perms = newStack()
		var uris = newStack()
		nodes.push(ns)
		perms.push("publish,subcribe")
		uris.push("*")
		for nodes.length() > 0 {
			active_entity := nodes.pop()
			//active_perms := perms.pop()
			perms.pop()
			active_uri := uris.pop()

			if visited_entities[active_entity] {
				continue
			}

			visited_entities[active_entity] = true

			_entities, _permissions, _uris := find1Hop(active_entity)
			for idx, reachable_entity := range _entities {
				nodes.push(reachable_entity)
				restricted, ok := RestrictBy(active_uri, _uris[idx])
				if !ok {
					break
				}
				uris.push(restricted)
				//TODO: restrict permissions
				perms.push(_permissions[idx])
			}
		}

	}

	// check which entities were visited, keeping track of unreachable nodes
	all_visited := true
	var unreachable []string
	for ent, reachable := range visited_entities {
		if !reachable {
			all_visited = false
			unreachable = append(unreachable, ent)
		}
	}

	return all_visited, unreachable
}

// do the following for all entities (namespaces, from, to, permissionsets) to
// prepare for fulfilling the graph:
// - if {entity name}.ent exists, pull in the file contents
// - if {entity name}.ent does not exist, create it
//
// Next have each entity register names for all other entities
// If (publish) is true, publish the entities after they are all
// created and have them name each other; otherwise, keep them offline
// and do not create names.
func (eng *GraphEngine) PrepareEntities(publish bool) error {
	all_entities := eng.spec.Entities()
	fmt.Println(all_entities)

	any_new := false

	for _, ent := range all_entities {
		filename := fmt.Sprintf("%s.ent", ent)
		if fileDoesNotExist(filename) {
			any_new = true
			// create file
			resp, err := eng.wave.CreateEntity(context.Background(), &pb.CreateEntityParams{
				ValidFrom:  time.Now().UnixNano() / 1e6,
				ValidUntil: time.Now().Add(eng.spec.Graph.EntityExpiry.Duration).UnixNano() / 1e6,
				RevocationLocation: &pb.Location{
					AgentLocation: "default",
				},
			})
			if err != nil {
				return xerrors.Errorf("Could not call CreateEntity: %w", err)
			}
			if resp.Error != nil {
				return xerrors.New(resp.Error.Message)
			}
			bl := pem.Block{
				Type:  eapi.PEM_ENTITY_SECRET,
				Bytes: resp.SecretDER,
			}
			err = ioutil.WriteFile(filename, pem.EncodeToMemory(&bl), 0600)
			if err != nil {
				return xerrors.Errorf("Could not write entity to %s: %w", filename, err)
			}
			stringhash := base64.URLEncoding.EncodeToString(resp.Hash)
			fmt.Println("Created entity", ent, stringhash)

			// publish if required
			if publish {
				presp, err := eng.wave.PublishEntity(context.Background(), &pb.PublishEntityParams{
					DER: resp.PublicDER,
					Location: &pb.Location{
						AgentLocation: "default",
					},
				})
				if err != nil {
					return xerrors.Errorf("Could not call PublishEntity: %w", err)
				}
				if presp.Error != nil {
					return xerrors.Errorf("Could not publish entity: %w", err)
				}
			}
			// finish by loading entity bytes from file
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return xerrors.Errorf("Could not read entity file %s: %w", filename, err)
			}
			eng.entities[ent] = content
			// we have the hash from the parsed object
			eng.hashes[ent] = resp.Hash
		} else {
			fmt.Println("Entity", ent, "already exists locally")
			// read existing file
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return xerrors.Errorf("Could not read entity file %s: %w", filename, err)
			}
			// save bytes
			eng.entities[ent] = content
			// compute hash
			hash, err := eng.getEntityHash(content)
			if err != nil {
				return xerrors.Errorf("Could not get entity hash from %s: %w", filename, err)
			}
			eng.hashes[ent] = hash
		}

	}

	if !publish || !any_new {
		return nil
	}

	// handle naming: all entities name each other
	// TODO: for now, this computes all (n-1)*(n-1) pairwise namings
	// each time it is run.
	for namingentity, entitybytes := range eng.entities {
		for namedentity := range eng.entities {
			//if namingentity == namedentity {
			//	continue
			//}

			resp, err := eng.wave.CreateNameDeclaration(context.Background(), &pb.CreateNameDeclarationParams{
				Perspective: &pb.Perspective{
					EntitySecret: &pb.EntitySecret{
						DER: entitybytes,
					},
				},
				Name:       strings.ToLower(namedentity),
				Subject:    eng.hashes[namedentity],
				ValidFrom:  time.Now().UnixNano() / 1e6,
				ValidUntil: time.Now().Add(5*365*24*time.Hour).UnixNano() / 1e6,
			})
			if err != nil {
				return xerrors.Errorf("Could not call CreateNameDeclaration: %w", err)
			}
			if resp.Error != nil {
				return xerrors.Errorf("Could not create name declaration: %s", resp.Error.Message)
			}
			fmt.Println(namingentity, "named", namedentity)

		}
	}

	return nil
}

func (eng *GraphEngine) PrepareEdges() error {
	for _, e := range eng.spec.Edges {
		fmt.Printf("%+v\n", e)
	}

	//create + publish attestations
	for _, att := range eng.spec.Edges {
		params := &pb.CreateAttestationParams{
			Perspective: &pb.Perspective{
				EntitySecret: &pb.EntitySecret{DER: eng.entities[att.From]},
			},
			SubjectHash: eng.hashes[att.To],
			ValidFrom:   time.Now().UnixNano() / 1e6,
			ValidUntil:  time.Now().Add(att.Expiry.Duration).UnixNano() / 1e6,
			Policy: &pb.Policy{
				RTreePolicy: &pb.RTreePolicy{
					Namespace:    eng.hashes[att.Namespace],
					Indirections: uint32(att.TTL),
					Statements: []*pb.RTreePolicyStatement{
						{
							PermissionSet: eng.hashes[att.Pset],
							Permissions:   strings.SplitN(att.Permissions, ",", -1),
							Resource:      att.Resource,
						},
					},
				},
			},
		}
		resp, err := eng.wave.CreateAttestation(context.Background(), params)
		if err != nil {
			return err
		}
		if resp.Error != nil {
			log.Println(att.From, att.To)
			return xerrors.Errorf(resp.Error.Message)
		}

		presp, err := eng.wave.PublishAttestation(context.Background(), &pb.PublishAttestationParams{
			DER: resp.DER,
		})
		if err != nil {
			return err
		}
		if presp.Error != nil {
			return xerrors.Errorf(presp.Error.Message)
		}
		stringhash := base64.URLEncoding.EncodeToString(resp.Hash)
		fmt.Println("Published Attestation", stringhash)
	}

	return nil
}

func getConn(agent string) pb.WAVEClient {
	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to agent: %v\n", err)
	}
	client := pb.NewWAVEClient(conn)
	return client
}

// TODO:
// need separate commands for
// - validate
// - make entities
// - make attestations
// - (or all 3 steps)
// - OPTIONAL: dump RDF representation
// - INSTEAD: generate instructions so entities/edges can be done non-centrally
// - terminal nodes can be hashes; we don't create these, just grant to the hash
func main() {
	g := GraphEngineFromSpecFile("energise.toml")
	all_visited, unreachable := g.ValidateConnectivity()
	fmt.Println("Fully connected?", all_visited, "unreachable:", unreachable)
	fmt.Println("Prepare entities")

	//t := newTraversal(g.spec, "lpbc_675")
	//t.traverse()
	for _, ent := range g.spec.Entities() {
		t := newTraversal(g.spec, ent)
		t.traverse()
	}

	return
	if err := g.PrepareEntities(true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Prepare edges")
	if err := g.PrepareEdges(); err != nil {
		log.Fatal(err)
	}
	g.ToRDF()
}

func fileDoesNotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// Copied verbatim from bosswave
// RestrictBy takes a topic, and a permission, and returns the intersection
// that represents the from topic restricted by the permission. It took a
// looong time to work out this logic...
func RestrictBy(from string, by string) (string, bool) {
	fp := strings.Split(from, "/")
	bp := strings.Split(by, "/")
	fout := make([]string, 0, len(fp)+len(bp))
	bout := make([]string, 0, len(fp)+len(bp))
	var fsx, bsx int
	for fsx = 0; fsx < len(fp) && fp[fsx] != "*"; fsx++ {
	}
	for bsx = 0; bsx < len(bp) && bp[bsx] != "*"; bsx++ {
	}
	fi, bi := 0, 0
	fni, bni := len(fp)-1, len(bp)-1
	emit := func() (string, bool) {
		for i := 0; i < len(bout); i++ {
			fout = append(fout, bout[len(bout)-i-1])
		}
		return strings.Join(fout, "/"), true
	}
	//phase 1
	//emit matching prefix
	for ; fi < len(fp) && bi < len(bp); fi, bi = fi+1, bi+1 {
		if fp[fi] != "*" && (fp[fi] == bp[bi] || (bp[bi] == "+" && fp[fi] != "*")) {
			fout = append(fout, fp[fi])
		} else if fp[fi] == "+" && bp[bi] != "*" {
			fout = append(fout, bp[bi])
		} else {
			break
		}
	}
	//phase 2
	//emit matching suffix
	for ; fni >= fi && bni >= bi; fni, bni = fni-1, bni-1 {
		if bp[bni] != "*" && (fp[fni] == bp[bni] || (bp[bni] == "+" && fp[fni] != "*")) {
			bout = append(bout, fp[fni])
		} else if fp[fni] == "+" && bp[bni] != "*" {
			bout = append(bout, bp[bni])
		} else {
			break
		}
	}
	//phase 3
	//emit front
	if fi < len(fp) && fp[fi] == "*" {
		for ; bi < len(bp) && bp[bi] != "*" && bi <= bni; bi++ {
			fout = append(fout, bp[bi])
		}
	} else if bi < len(bp) && bp[bi] == "*" {
		for ; fi < len(fp) && fp[fi] != "*" && fi <= fni; fi++ {
			fout = append(fout, fp[fi])
		}
	}
	//phase 4
	//emit back
	if fni >= 0 && fp[fni] == "*" {
		for ; bni >= 0 && bp[bni] != "*" && bni >= bi; bni-- {
			bout = append(bout, bp[bni])
		}
	} else if bni >= 0 && bp[bni] == "*" {
		for ; fni >= 0 && fp[fni] != "*" && fni >= fi; fni-- {
			bout = append(bout, fp[fni])
		}
	}
	//phase 5
	//emit star if they both have it
	if fi == fni && fp[fi] == "*" && bi == bni && bp[bi] == "*" {
		fout = append(fout, "*")
		return emit()
	}
	//Remove any stars
	if fi < len(fp) && fp[fi] == "*" {
		fi++
	}
	if bi < len(bp) && bp[bi] == "*" {
		bi++
	}
	if (fi == fni+1 || fi == len(fp)) && (bi == bni+1 || bi == len(bp)) {
		return emit()
	}
	return "", false
}
