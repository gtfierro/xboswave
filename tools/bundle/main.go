package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gtfierro/hoddb/hod"
	hodpb "github.com/gtfierro/hoddb/proto"
	"github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var config = flag.String("config", "hodconfig.yml", "Path to hodconfig.yml file")

type BrickModel struct {
	// name of the Brick model
	Name string
	// timestamp version
	Version time.Time
	// hash of model?
	Hash []byte
}

type Predicate struct {
	BrickModel *BrickModel
	Vars       []string
	Triples    []hodpb.Triple
}

type Attestation struct {
	// recipient of the grant
	SubjectHash []byte
	// validity bounds on attestation
	ValidFrom  time.Time
	ValidUntil time.Time
	Predicate  Predicate
}

func getConn(agent string) pb.WAVEClient {
	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to agent: %v\n", err)
	}
	client := pb.NewWAVEClient(conn)
	return client
}

func getPerspective(filename string) *pb.Perspective {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("could not read file %q: %v\n", filename, err)
	}
	return &pb.Perspective{
		EntitySecret: &pb.EntitySecret{
			DER: contents,
		},
	}
}

// TODO: do we want to compute the minimal set of edges;
// if we can grant fewer edges, we should.
// For example, if we grant all URIs in the model that share
// a prefix, we can just grant the prefix/*?
type Bundle struct {
	uris []string
}

func (b *Bundle) addURI(uri string) {
	b.uris = append(b.uris, uri)
}

type Bundler struct {
	client      pb.WAVEClient
	perspective *pb.Perspective
	hod         *hod.HodDB
}

func NewBundler(cfg *hod.Config) *Bundler {
	hod, err := hod.MakeHodDB(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "open hoddb"))
	}
	client := getConn(os.Getenv("WAVE_AGENT"))
	perspective := getPerspective(os.Getenv("WAVE_DEFAULT_ENTITY"))

	return &Bundler{
		client:      client,
		perspective: perspective,
		hod:         hod,
	}
}

func (b *Bundler) getURIs(query string) (Bundle, error) {
	parsed, err := b.hod.ParseQuery(query, 0)
	if err != nil {
		return Bundle{}, errors.Wrap(err, "could not parse query")
	}

	fmt.Printf("parsed %+v\n", parsed)
	var uri_vars []string
	for _, triple := range parsed.Where {
		if triple.Predicate[0].Namespace == `https://brickschema.org/schema/1.0.3/BrickFrame` && triple.Predicate[0].Value == "uri" {
			uri_vars = append(uri_vars, triple.Object.Value)
		}
	}

	res, err := b.hod.Select(context.Background(), parsed)
	if err != nil {
		return Bundle{}, errors.Wrap(err, "could not run query")
	}
	fmt.Printf("%+v\n", res)
	var uriidxs []int
	for idx, varname := range res.Variables {
		for _, urivar := range uri_vars {
			if varname == urivar {
				uriidxs = append(uriidxs, idx)
			}
		}
	}
	var bundle Bundle
	for _, row := range res.Rows {
		for _, idx := range uriidxs {
			bundle.addURI(row.Values[idx].Value)
		}
	}
	return bundle, nil
}

func main() {
	flag.Parse()

	// create the hod database
	cfg, err := hod.ReadConfig(*config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not load config file"))
	}

	bundler := NewBundler(cfg)

	bundle, err := bundler.getURIs(`SELECT ?uri WHERE { ?tstat rdf:type brick:Thermostat . ?tstat bf:uri ?uri };`)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not get bundle"))
	}

	fmt.Println(bundle)
}
