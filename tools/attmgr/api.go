package main

import (
	"context"
	"io/ioutil"

	"github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
)

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
