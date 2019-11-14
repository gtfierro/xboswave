package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/immesys/wave/eapi"
	"github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
)

func (db *DB) LoadEntityFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "Could not read file")
	}

	resp, err := db.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: content,
	})
	if err != nil {
		return errors.Wrapf(err, "Could not inspect file %s", filename)
	}
	if resp.Error != nil {
		return fmt.Errorf("Could not inspect: %s", resp.Error.Message)
	}

	if resp.Entity == nil {
		return fmt.Errorf("%s was not an entity", filename)
	}
	ent := ParseEntity(resp.Entity)
	return db.insertEntity(ent)
}

func (db *DB) LoadAttestationFile(filename string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	resp, err := db.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: contents,
	})
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	if resp.Attestation == nil {
		return errors.New("Inspection failed (not an attestation?)")
	}

	resolve, err := db.wave.ResolveHash(context.Background(), &pb.ResolveHashParams{
		Perspective: db.perspective,
		Hash:        resp.Attestation.Hash,
	})
	if err != nil {
		return err
	}
	if resolve.Error != nil {
		return errors.New(resolve.Error.Message)
	}
	if resolve.Attestation == nil {
		return errors.New("Resolve failed (not an attestation?)")
	}
	att := ParseAttestation(resolve.Attestation)
	return db.insertAttestation(att)
}

func (db *DB) CreatePolicy(namespace, resource, pset string, indirections int, permissions []string) error {
	//namespace = resolveEntityNameOrHashOrFile(db.wave, db.perspective, namespace, "missing subject entity")
	//pset := resolveEntityNameOrHashOrFile(db.wave, db.perspective, pset, "missing subject entity")
	policy := &PolicyStatement{
		Namespace:     namespace,
		PermissionSet: pset,
		Indirections:  uint32(indirections),
		Permissions:   permissions,
		Resource:      resource,
	}
	return db.insertPolicy(policy)
}

func (db *DB) CreateAttestation(subjectHashOrFile string, ValidUntil time.Time, policies []PolicyStatement) error {
	if len(policies) == 0 {
		return fmt.Errorf("Need > 0 policies")
	}
	possiblename := db.getHashFromName(subjectHashOrFile)
	subject := resolveEntityNameOrHashOrFile(db.wave, db.perspective, possiblename, "missing subject entity")

	subjresp, err := db.wave.ResolveHash(context.Background(), &pb.ResolveHashParams{
		Hash: subject,
	})
	if err != nil {
		return fmt.Errorf("could not find subject location: %v\n", err)
	}
	if subjresp.Error != nil {
		return fmt.Errorf("could not find subject location: %v\n", subjresp.Error.Message)
	}

	params := &pb.CreateAttestationParams{
		Perspective:     db.perspective,
		BodyScheme:      eapi.BodySchemeWaveRef1,
		SubjectHash:     subject,
		SubjectLocation: subjresp.Location,
		ValidFrom:       time.Now().UnixNano() / 1e6,
		ValidUntil:      ValidUntil.UnixNano() / 1e6,
	}

	statements := []*pb.RTreePolicyStatement{}

	for _, policy := range policies {
		pset := resolveEntityNameOrHashOrFile(db.wave, db.perspective, policy.PermissionSet, "bad permission set")
		stmt := &pb.RTreePolicyStatement{
			PermissionSet: pset,
			Permissions:   policy.Permissions,
			Resource:      policy.Resource,
		}
		statements = append(statements, stmt)
	}

	ns := resolveEntityNameOrHashOrFile(db.wave, db.perspective, policies[0].Namespace, "bad permission set")
	params.Policy = &pb.Policy{
		RTreePolicy: &pb.RTreePolicy{
			Namespace:    ns,
			Indirections: policies[0].Indirections,
			Statements:   statements,
		},
	}

	resp, err := db.wave.CreateAttestation(context.Background(), params)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("error: %v\n", resp.Error.Message)
	}
	bl := pem.Block{
		Type:  eapi.PEM_ATTESTATION,
		Bytes: resp.DER,
	}
	stringhash := base64.URLEncoding.EncodeToString(resp.Hash)
	log.Infof("Created Attestation %s", stringhash)

	outfilename := fmt.Sprintf("att_%s.pem", stringhash)
	err = ioutil.WriteFile(outfilename, pem.EncodeToMemory(&bl), 0600)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not write attestation file"))
	}
	presp, err := db.wave.PublishAttestation(context.Background(), &pb.PublishAttestationParams{
		DER: resp.DER,
	})
	if err != nil {
		return err
	}
	if presp.Error != nil {
		return errors.New(presp.Error.Message)
	}
	log.Infof("Published Attestation %s", stringhash)

	return nil
}

// check that this attestation can be used in a real proof by
// seeing if we can build a proof to what it is granting
func (db *DB) Validate(att Attestation) error {

	for _, policy := range att.PolicyStatements {
		pset := resolveEntityNameOrHashOrFile(db.wave, db.perspective, policy.PermissionSet, "bad permission set")
		ns := resolveEntityNameOrHashOrFile(db.wave, db.perspective, policy.Namespace, "bad namesapce")
		params := &pb.BuildRTreeProofParams{
			Perspective: db.perspective,
			// subject is the perspective by default
			Namespace: ns,
			Statements: []*pb.RTreePolicyStatement{
				&pb.RTreePolicyStatement{
					PermissionSet: pset,
					Permissions:   policy.Permissions,
					Resource:      policy.Resource,
				},
			},
		}

		resp, err := db.wave.BuildRTreeProof(context.Background(), params)
		if err != nil {
			return err
		}
		if resp.Error != nil {
			return errors.New(resp.Error.Message)
		}

		vresp, err := db.wave.VerifyProof(context.Background(), &pb.VerifyProofParams{
			ProofDER: resp.ProofDER,
		})
		if err != nil {
			return err
		}
		if vresp.Error != nil {
			return errors.New(vresp.Error.Message)
		}
		//proof := vresp.Result
		//fmt.Printf("  Validity:\n")
		//fmt.Printf("   - Readable: %v\n", !proof.Attestation.Validity.NotDecrypted)
		//fmt.Printf("   - Revoked: %v\n", proof.Attestation.Validity.Revoked)
		//fmt.Printf("   - Malformed: %v\n", proof.Attestation.Validity.Malformed)
		//fmt.Printf("   - Subject invalid: %v\n", proof.Attestation.Validity.DstInvalid)
		//if !proof.Attestation.Validity.NotDecrypted {
		//	fmt.Printf("   - Valid: %v\n", proof.Attestation.Validity.Valid)
		//	fmt.Printf("   - Expired: %v\n", proof.Attestation.Validity.Expired)
		//	fmt.Printf("   - Attester invalid: %v\n", proof.Attestation.Validity.SrcInvalid)
		//}

	}
	return nil
}

func (db *DB) ListExpiring(within_next time.Duration) ([]Attestation, error) {
	expiring_before := time.Now().Add(within_next)

	f := &filter{
		expiring_before: &expiring_before,
	}

	atts, err := db.listAttestation(f)
	if err != nil {
		return nil, err
	}
	return atts, nil
}

func (db *DB) syncgraph() error {
	resp, err := db.wave.ResyncPerspectiveGraph(context.Background(), &pb.ResyncPerspectiveGraphParams{
		Perspective: db.perspective,
	})
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}
	srv, err := db.wave.WaitForSyncComplete(context.Background(), &pb.SyncParams{
		Perspective: db.perspective,
	})
	for {
		rv, err := srv.Recv()
		if err == io.EOF {
			break
		}
		fmt.Printf("Synchronized %d/%d entities\n", rv.CompletedSyncs, rv.TotalSyncRequests)
	}
	fmt.Printf("Perspective graph sync complete\n")
	return nil
}

func (db *DB) getNameFromHash(hash string) (name string) {
	// check wave, wavemq
	if hash == "GyAUM3SzL9J0OVT-R4b2z4bUA3IPXsRCNrZYwmoeaA9uAQ==" {
		return "wavemq"
	}
	bytehash, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		log.Debug(errors.Wrapf(err, "Hash %v was not base64", hash))
		return hash
	}
	resp, err := db.wave.ResolveReverseName(context.Background(), &pb.ResolveReverseNameParams{
		Perspective: db.perspective,
		Hash:        bytehash,
	})
	if err != nil {
		log.Debug(errors.Wrapf(err, "Could not resolve name of hash %s", hash))
		return hash
	}
	if resp.Error != nil {
		log.Debug(errors.Errorf("Could not resolve name of hash %s: %v", hash, resp.Error.Message))
		return hash
	}
	return resp.Name
}

func (db *DB) getHashFromName(name string) (hash string) {
	resp, err := db.wave.ResolveName(context.Background(), &pb.ResolveNameParams{
		Perspective: db.perspective,
		Name:        name,
	})
	if err != nil {
		log.Debug(errors.Wrap(err, "could not resolve hash for name"))
		return name
	}
	if resp.Error != nil {
		log.Debug(errors.Errorf("could not resolve hash for name: %v", resp.Error))
		return name
	}
	return base64.URLEncoding.EncodeToString(resp.Entity.Hash)
}

func (db *DB) getEntityFromFile(file string) (*Entity, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file")
	}

	resp, err := db.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: content,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not inspect file %s", file)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("Could not inspect: %s", resp.Error.Message)
	}

	if resp.Entity == nil {
		return nil, fmt.Errorf("%s was not an entity", file)
	}
	return ParseEntity(resp.Entity), nil
}
