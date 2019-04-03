package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	"github.com/immesys/wave/consts"
	"github.com/immesys/wave/eapi"
	"github.com/immesys/wave/eapi/pb"
	"google.golang.org/grpc"
)

// this is a "stringy" version of pb.Attestation that is easier to
type Attestation struct {
	// perspective is implicit (its who is running this program)

	// base64 encoded
	Subject string
	// base64 encoded
	Hash             string
	ValidFrom        time.Time
	ValidUntil       time.Time
	PolicyStatements []PolicyStatement
}

func ParseAttestation(att *pb.Attestation) Attestation {
	return Attestation{
		Subject:          base64.URLEncoding.EncodeToString(att.SubjectHash),
		Hash:             base64.URLEncoding.EncodeToString(att.Hash),
		ValidFrom:        time.Unix(0, att.Body.ValidFrom*1e6),
		ValidUntil:       time.Unix(0, att.Body.ValidUntil*1e6),
		PolicyStatements: ParsePolicyStatement(att.Body),
	}
}

type PolicyStatement struct {
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

func resolveEntityNameOrHashOrFile(conn pb.WAVEClient, perspective *pb.Perspective, in string, msg string) (hash []byte) {
	f, err := ioutil.ReadFile(in)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error opening file %q: %v\n", in, err)
			os.Exit(1)
		}
		//Resolve as name/hash
		if len(in) == 48 && strings.Index(in, ".") == -1 {
			//Resolve as hash
			rv, err := base64.URLEncoding.DecodeString(in)
			if err != nil {
				fmt.Printf("bad base64: %q\n", in)
				os.Exit(1)
			}
			return rv
		}
		//Resolve as name
		if in == "wave" {
			//Hardcoded builtin PSET
			rv, _ := base64.URLEncoding.DecodeString(consts.WaveBuiltinPSET)
			return rv
		} else if in == "wavemq" {
			return []byte("\x1b\x20\x14\x33\x74\xb3\x2f\xd2\x74\x39\x54\xfe\x47\x86\xf6\xcf\x86\xd4\x03\x72\x0f\x5e\xc4\x42\x36\xb6\x58\xc2\x6a\x1e\x68\x0f\x6e\x01")
		}

		resp, err := conn.ResolveName(context.Background(), &pb.ResolveNameParams{
			Perspective: perspective,
			Name:        in,
		})
		if err != nil {
			fmt.Printf("could not resolve name: %v\n", err)
			os.Exit(1)
		}
		if resp.Error != nil {
			fmt.Printf("could not resolve name %q: %s\n", in, resp.Error.Message)
			os.Exit(1)
		}
		return resp.Entity.Hash
	}
	//Resolve as file
	resp, err := conn.Inspect(context.Background(), &pb.InspectParams{
		Content: f,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	if resp.Error != nil {
		fmt.Printf("could not inspect file: %s\n", resp.Error.Message)
		os.Exit(1)
	}
	if resp.Entity != nil {
		return resp.Entity.Hash
	}
	fmt.Printf(msg)
	os.Exit(1)
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

func getPerspective(file string, passphrase string, msg string) *pb.Perspective {
	if file != "" {
		pass := []byte(passphrase)
		if len(pass) == 0 {
			fmt.Printf("passphrase for entity secret: ")
			var err error
			pass, err = gopass.GetPasswdMasked()
			if err != nil {
				fmt.Printf("could not read passphrase: %v\n", err)
				os.Exit(1)
			}
		}
		pder := loadEntitySecretDER(file)
		perspective := &pb.Perspective{
			EntitySecret: &pb.EntitySecret{
				DER:        pder,
				Passphrase: pass,
			},
		}
		return perspective
	} else {
		fmt.Printf(msg)
		os.Exit(1)
		return nil
	}
}

func loadEntitySecretDER(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("could not read file %q: %v\n", filename, err)
		os.Exit(1)
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		fmt.Printf("file %q is not a PEM file\n", filename)
		os.Exit(1)
	}
	if block.Type != eapi.PEM_ENTITY_SECRET {
		fmt.Printf("PEM is not an entity secret\n")
		os.Exit(1)
	}
	return block.Bytes
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

func main() {
	cfg := &Config{
		Path:        "attestations.sqlite3",
		Agent:       "localhost:410",
		Perspective: "test.ent",
	}
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("load att1.pem", db.LoadAttestationFile("att1.pem"))
	db.watch(".")

	_ = db
}
