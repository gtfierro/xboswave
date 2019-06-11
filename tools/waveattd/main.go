package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/immesys/wave/eapi/pb"
	logrus "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

// this is a "stringy" version of pb.Attestation that is easier to
type Attestation struct {
	id int
	// base64 encoded
	Attester string
	// base64 encoded
	Subject string
	// base64 encoded
	Hash             string
	ValidFrom        time.Time
	ValidUntil       time.Time
	PolicyStatements []PolicyStatement
}

func ParseAttestation(att *pb.Attestation) *Attestation {
	if att.Body == nil {
		return nil
	}
	return &Attestation{
		Attester:         base64.URLEncoding.EncodeToString(att.Body.AttesterHash),
		Subject:          base64.URLEncoding.EncodeToString(att.SubjectHash),
		Hash:             base64.URLEncoding.EncodeToString(att.Hash),
		ValidFrom:        time.Unix(0, att.Body.ValidFrom*1e6),
		ValidUntil:       time.Unix(0, att.Body.ValidUntil*1e6),
		PolicyStatements: ParsePolicyStatement(att.Body),
	}
}

type PolicyStatement struct {
	id int
	// base64 encoded
	Namespace string
	// base64 encoded
	PermissionSet string
	Indirections  uint32
	Permissions   []string
	Resource      string
}

func ParsePolicyStatement(body *pb.AttestationBody) (stmts []PolicyStatement) {
	policy := body.Policy.RTreePolicy
	for i := 0; i < len(policy.Statements); i++ {
		stmt := PolicyStatement{
			Namespace:     base64.URLEncoding.EncodeToString(policy.Namespace),
			Indirections:  policy.Indirections,
			Permissions:   policy.Statements[i].Permissions,
			PermissionSet: base64.URLEncoding.EncodeToString(policy.Statements[i].PermissionSet),
			Resource:      policy.Statements[i].Resource,
		}
		stmts = append(stmts, stmt)
	}
	return
}

type Config struct {
	Path        string
	Agent       string
	Perspective string
}

// Required for the whole tool
// - perspective (entity file or contents)

// Required for an attestation
// - subject name/hash
// - namespace name/hash
// - permission set + permissions
// - resource
// - number of indirections
// - expiry

// API
//
// Naming
// - create alias for a hash (namespace, entity)
//
// Portions of attestations
// - list namespaces ( alias => hash )
// - list attestation subjects
// - list resources + namespaces
//
// Attesting
// - create attestation
// - renew attestation
//
// Listing
// - list all attestions
// - order by expiry
// - list by namespace, resource, permission, subject, etc

//func (db *DB) GrantAttestation(attester, passphrase string) {
//	perspective := getPerspective(attester, passphrase, "missing attesting entity secret\n")
//	resp, err := db.wave.ResyncPerspectiveGraph(context.Background(), &pb.ResyncPerspectiveGraphParams{
//		Perspective: perspective,
//	})
//	if err != nil {
//		fmt.Printf("error: %v\n", err)
//		os.Exit(1)
//	}
//	if resp.Error != nil {
//		fmt.Printf("error: %v\n", resp.Error.Message)
//		os.Exit(1)
//	}
//	srv, err := db.wave.WaitForSyncComplete(context.Background(), &pb.SyncParams{
//		Perspective: perspective,
//	})
//	for {
//		rv, err := srv.Recv()
//		if err == io.EOF {
//			break
//		}
//		fmt.Printf("Synchronized %d/%d entities\n", rv.CompletedSyncs, rv.TotalSyncRequests)
//	}
//	fmt.Printf("Perspective graph sync complete\n")
//}

//ParseDuration is a little like the existing time.ParseDuration
//but adds days and years because its really annoying not having that
// from Michael Andersen
func ParseDuration(s string) (*time.Duration, error) {
	if s == "" {
		return nil, nil
	}
	pat := regexp.MustCompile(`^(\d+y)?(\d+d)?(\d+h)?(\d+m)?(\d+s)?$`)
	res := pat.FindStringSubmatch(s)
	if res == nil {
		return nil, fmt.Errorf("Invalid duration")
	}
	res = res[1:]
	sec := int64(0)
	for idx, mul := range []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1} {
		if res[idx] != "" {
			key := res[idx][:len(res[idx])-1]
			v, e := strconv.ParseInt(key, 10, 64)
			if e != nil { //unlikely
				return nil, e
			}
			sec += v * mul
		}
	}
	rv := time.Duration(sec) * time.Second
	return &rv, nil
}

const DATABASE_LOCATION = "WAVEATTD_DB_LOCATION"
const WAVE_ENTITY = "WAVE_DEFAULT_ENTITY"
const WAVE_AGENT = "WAVE_AGENT"

func main() {

	var (
		location string
		entity   string
		agent    string
		found    bool
	)

	location, found = os.LookupEnv(DATABASE_LOCATION)
	if !found {
		location = "waveattd.sqlite3"
	}
	entity, found = os.LookupEnv(WAVE_ENTITY)
	if !found {
		log.Fatal("Set WAVE_DEFAULT_ENTITY")
	}
	agent, found = os.LookupEnv(WAVE_AGENT)
	if !found {
		agent = "localhost:410"
	}

	log.Info("╒ WAVEATTD_DB_LOCATION: ", location)
	log.Info("╞ WAVE_DEFAULT_ENTITY: ", entity)
	log.Info("╘ WAVE_AGENT: ", agent)

	cfg := &Config{
		Path:        location,
		Agent:       agent,
		Perspective: entity,
	}
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go db.watch(".")

	db.RunShell()

	_ = db
}
