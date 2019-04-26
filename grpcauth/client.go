package grpcauth

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
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
	proof           []byte
	namespace       string
	wave            pb.WAVEClient
}

func NewClientCredentials(perspective *pb.Perspective, agent string, namespace string, grpcservice string) (*ClientCredentials, error) {

	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to agent at %s", agent)
	}
	wave := pb.NewWAVEClient(conn)

	cc := &ClientCredentials{
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

	header, err := cc.GeneratePeerHeader()
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not generate header: %v", err)
	}
	_, err = conn.Write(header)
	if err != nil {
		rawConn.Close()
		return nil, nil, err
	}

	hdr, err := cc.ReadPeerHeader(conn)
	if err != nil {
		rawConn.Close()
		return nil, nil, errors.Wrap(err, "Could not read server header")
	}

	err = cc.VerifyServerHandshake(cc.namespace, hdr, cs.PeerCertificates[0].Signature)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not verify server handshake")
	}

	return conn, nil, nil
}

func (cc *ClientCredentials) VerifyServerHandshake(nsString string, hdr serverHeader, cert []byte) error {
	log.Info("Client verifying server handshake ", nsString)
	resp, err := cc.wave.VerifySignature(context.Background(), &pb.VerifySignatureParams{
		Signer:    hdr.entityHash,
		Signature: hdr.signature,
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
		ProofDER: hdr.proof,
		Subject:  hdr.entityHash,
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
	if !bytes.Equal(presp.Result.Subject, hdr.entityHash) {
		return errors.New("proof valid but for a different entity")
	}
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
	return &ClientCredentials{
		perspective:     cc.perspective,
		perspectiveHash: cc.perspectiveHash,
		namespace:       cc.namespace,
		wave:            cc.wave,
	}
}

func (cc *ClientCredentials) OverrideServerName(name string) error {
	return nil
}

func (cc *ClientCredentials) AddGRPCProofFile(filename string) (ns string, proof []byte, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil, errors.Wrap(err, "could not read designated routing file")
	}

	der := contents
	pblock, _ := pem.Decode(contents)
	if pblock != nil {
		der = pblock.Bytes
	}

	resp, err := cc.wave.VerifyProof(context.Background(), &pb.VerifyProofParams{
		ProofDER: der,
	})
	if err != nil {
		return "", nil, errors.Wrap(err, "could not verify dr file")
	}
	if resp.Error != nil {
		return "", nil, fmt.Errorf("could not verify dr file: %v", resp.Error.Message)
	}

	ns = base64.URLEncoding.EncodeToString(resp.Result.Policy.RTreePolicy.Namespace)
	//Check proof actually grants the right permissions:
	found := false
outer:
	for _, s := range resp.Result.Policy.RTreePolicy.Statements {
		if s.Resource == cc.grpcservice && bytes.Equal(s.GetPermissionSet(), []byte(XBOSPermissionSet)) {
			for _, perm := range s.Permissions {
				//TODO: need to MATCH the uri here for each of the uris, make sure we prove it
				if perm == GRPCCallPermission {
					found = true
					break outer
				}
			}
		}
	}

	if !found {
		return "", nil, fmt.Errorf("designated routing proof does not actually prove xbos:serve_grpc on any namespace")
	}
	cc.namespace = ns
	cc.proof = der

	return ns, der, nil
}

// client hash
// signature length
// signature (over proof)
// proof length
// proof
func (cc *ClientCredentials) GeneratePeerHeader() ([]byte, error) {
	hdr := bytes.Buffer{}
	if len(cc.perspectiveHash) != 34 {
		panic(cc.perspectiveHash)
	}
	//First: 34 byte entity hash
	hdr.Write(cc.perspectiveHash)
	//Second: signature of cert
	sigresp, err := cc.wave.Sign(context.Background(), &pb.SignParams{
		Perspective: cc.perspective,
		Content:     cc.proof,
	})
	if err != nil {
		return nil, err
	}
	if sigresp.Error != nil {
		return nil, errors.New(sigresp.Error.Message)
	}
	siglen := make([]byte, 2)
	sig := sigresp.Signature
	binary.LittleEndian.PutUint16(siglen, uint16(len(sig)))
	hdr.Write(siglen)
	hdr.Write(sig)

	//Third: the namespace proof for this namespace
	prooflen := make([]byte, 4)
	binary.LittleEndian.PutUint32(prooflen, uint32(len(cc.proof)))
	hdr.Write(prooflen)
	hdr.Write(cc.proof)
	return hdr.Bytes(), nil
}

func (cc *ClientCredentials) ReadPeerHeader(conn io.Reader) (serverHeader, error) {
	var (
		err error
		hdr serverHeader
	)
	entityHashBA := make([]byte, 34)
	_, err = io.ReadFull(conn, entityHashBA)
	if err != nil {
		return hdr, fmt.Errorf("could not read proof: %v\n", err)
	}

	signatureSizeBA := make([]byte, 2)
	_, err = io.ReadFull(conn, signatureSizeBA)
	if err != nil {
		return hdr, fmt.Errorf("could not read proof: %v\n", err)
	}

	signatureSize := binary.LittleEndian.Uint16(signatureSizeBA)
	signature := make([]byte, signatureSize)
	_, err = io.ReadFull(conn, signature)
	if err != nil {
		return hdr, fmt.Errorf("could not read proof: %v\n", err)
	}
	proofSizeBA := make([]byte, 4)
	_, err = io.ReadFull(conn, proofSizeBA)
	if err != nil {
		return hdr, fmt.Errorf("could not read proof: %v\n", err)
	}
	proofSize := binary.LittleEndian.Uint32(proofSizeBA)
	if proofSize > 10*1024*1024 {
		return hdr, fmt.Errorf("bad proof")
	}
	log.Debug("server read proof")
	proof := make([]byte, proofSize)
	_, err = io.ReadFull(conn, proof)
	if err != nil {
		return hdr, fmt.Errorf("could not read proof: %v\n", err)
	}

	hdr.entityHash = entityHashBA
	hdr.signature = signature
	hdr.proof = proof

	return hdr, nil
	////First verify the signature
	//log.Debug("Verify Server handshake")
	//err = cc.VerifyServerHandshake(cc.namespace, entityHashBA, signature, proof, cs.PeerCertificates[0].Signature)
	//if err != nil {
	//	return err
	//}
}

type serverHeader struct {
	entityHash []byte
	signature  []byte
	proof      []byte
}
