package grpcauth

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/cloudflare/cfssl/log"
	pb "github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientCredentials struct {
	perspective     *pb.Perspective
	grpcservice     string
	perspectiveHash []byte
	namespace       string
	wave            pb.WAVEClient
}

func NewClientCredentials(perspective *pb.Perspective, agent string, namespace string, grpcservice string) (*WaveCredentials, error) {

	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to agent at %s", agent)
	}
	wave := pb.NewWAVEClient(conn)

	cc := &WaveCredentials{
		perspective: perspective,
		wave:        wave,
	}

	// learn the perspective hash
	iresp, err := cc.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: perspective.EntitySecret.DER,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not inspect perspective entity file")
	}
	if iresp.Error != nil {
		return nil, errors.Wrap(err, "could not inspect perspective entity file")
	}
	cc.perspectiveHash = iresp.Entity.Hash

	cc.namespace = namespace
	cc.grpcservice = grpcservice

	return cc, nil
}

func (cc *ClientCredentials) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	roots := x509.NewCertPool()

	log.Debug("start client handshake")
	conn := tls.Client(rawConn, &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            roots,
	})

	err := conn.Handshake()
	if err != nil {
		rawConn.Close()
		return nil, nil, err
	}
	cs := conn.ConnectionState()
	if len(cs.PeerCertificates) != 1 {
		fmt.Printf("peer connection weird response")
		rawConn.Close()
		return nil, nil, errors.New("Wrong certificates")
	}

	nsbin, err := base64.URLEncoding.DecodeString(cc.namespace)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(nsbin)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}

	log.Debug("read ent hash")
	entityHashBA := make([]byte, 34)
	_, err = io.ReadFull(conn, entityHashBA)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}
	log.Debug("read sig size")
	signatureSizeBA := make([]byte, 2)
	_, err = io.ReadFull(conn, signatureSizeBA)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}
	signatureSize := binary.LittleEndian.Uint16(signatureSizeBA)
	signature := make([]byte, signatureSize)
	log.Debug("read sig")
	_, err = io.ReadFull(conn, signature)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}
	log.Debug("read proof size")
	proofSizeBA := make([]byte, 4)
	_, err = io.ReadFull(conn, proofSizeBA)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}
	proofSize := binary.LittleEndian.Uint32(proofSizeBA)
	if proofSize > 10*1024*1024 {
		rawConn.Close()
		return nil, nil, fmt.Errorf("bad proof")
	}
	log.Debug("read proof")
	proof := make([]byte, proofSize)
	_, err = io.ReadFull(conn, proof)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	}

	//First verify the signature
	log.Debug("Verify Server handshake")
	err = cc.VerifyServerHandshake(cc.namespace, entityHashBA, signature, proof, cs.PeerCertificates[0].Signature)
	if err != nil {
		rawConn.Close()
		return nil, nil, err
	}

	return conn, nil, nil
}

func (cc *ClientCredentials) VerifyServerHandshake(nsString string, entityHash []byte, signature []byte, proof []byte, cert []byte) error {
	log.Info("Verifying server handshake", nsString)
	resp, err := cc.wave.VerifySignature(context.Background(), &pb.VerifySignatureParams{
		Signer:    entityHash,
		Signature: signature,
		Content:   cert,
	})
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.New(resp.Error.Message)
	}

	ns, err := base64.URLEncoding.DecodeString(nsString)
	if err != nil {
		return err
	}

	//Signature ok, verify proof
	presp, err := cc.wave.VerifyProof(context.Background(), &pb.VerifyProofParams{
		ProofDER: proof,
		Subject:  entityHash,
		RequiredRTreePolicy: &pb.RTreePolicy{
			Namespace: ns,
			Statements: []*pb.RTreePolicyStatement{
				{
					PermissionSet: []byte(XBOSPermissionSet),
					Permissions:   []string{GRPCServePermission},
					// grpc_package/ServiceName/* (all methods)
					// grpc_package/ServiceName/Method1 (only method 1)
					Resource: cc.grpcservice,
				},
			},
		},
	})

	if err != nil {
		return err
	}
	if presp.Error != nil {
		return errors.New(presp.Error.Message)
	}
	if !bytes.Equal(presp.Result.Subject, entityHash) {
		return errors.New("proof valid but for a different entity")
	}
	log.Info(">", ns)
	log.Infof("%+v", presp.Result)
	return nil
}

func (cc *ClientCredentials) ServerHandshake(net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, errors.New("Not Implemented")
}

func (cc *ClientCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	}
}
func (cc *ClientCredentials) Clone() credentials.TransportCredentials {
	return &WaveCredentials{
		perspective:     cc.perspective,
		perspectiveHash: cc.perspectiveHash,
		namespace:       cc.namespace,
		wave:            cc.wave,
	}
}

func (cc *ClientCredentials) OverrideServerName(name string) error {
	return nil
}
