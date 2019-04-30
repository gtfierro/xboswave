package grpcauth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/log"
	pb "github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const XBOSPermissionSet_b64 = "GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA=="
const XBOSPermissionSet = "\x1b\x20\xb9\xc1\x45\x06\x28\xe3\x7a\xb8\x2e\x20\xc6\xe1\xfa\x4e\x9c\xd4\xf6\xfb\xee\x28\x71\x9e\x25\xad\x63\x52\xbe\x06\xfa\x3f\xcd\x6c"

const GRPCServePermission = "serve_grpc"
const GRPCCallPermission = "call_grpc"

// type TransportCredentials interface {
// 	// ClientHandshake does the authentication handshake specified by the corresponding
// 	// authentication protocol on rawConn for clients. It returns the authenticated
// 	// connection and the corresponding auth information about the connection.
// 	// Implementations must use the provided context to implement timely cancellation.
// 	// gRPC will try to reconnect if the error returned is a temporary error
// 	// (io.EOF, context.DeadlineExceeded or err.Temporary() == true).
// 	// If the returned error is a wrapper error, implementations should make sure that
// 	// the error implements Temporary() to have the correct retry behaviors.
// 	//
// 	// If the returned net.Conn is closed, it MUST close the net.Conn provided.
// 	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
// 	// ServerHandshake does the authentication handshake for servers. It returns
// 	// the authenticated connection and the corresponding auth information about
// 	// the connection.
// 	//
// 	// If the returned net.Conn is closed, it MUST close the net.Conn provided.
// 	ServerHandshake(net.Conn) (net.Conn, AuthInfo, error)
// 	// Info provides the ProtocolInfo of this TransportCredentials.
// 	Info() ProtocolInfo
// 	// Clone makes a copy of this TransportCredentials.
// 	Clone() TransportCredentials
// 	// OverrideServerName overrides the server name used to verify the hostname on the returned certificates from the server.
// 	// gRPC internals also use it to override the virtual hosting name if it is set.
// 	// It must be called before dialing. Currently, this is only used by grpclb.
// 	OverrideServerName(string) error
// }

type WaveCredentials struct {
	perspective     *pb.Perspective
	info            map[string]grpc.ServiceInfo
	grpcservice     string
	perspectiveHash []byte
	proof           []byte
	namespace       string
	wave            pb.WAVEClient
}

func NewServerCredentials(perspective *pb.Perspective, agent string) (*WaveCredentials, error) {

	conn, err := grpc.Dial(agent, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to agent at %s", agent)
	}
	wave := pb.NewWAVEClient(conn)

	wc := &WaveCredentials{
		perspective: perspective,
		wave:        wave,
	}

	// learn the perspective hash
	iresp, err := wc.wave.Inspect(context.Background(), &pb.InspectParams{
		Content: perspective.EntitySecret.DER,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not inspect perspective entity file")
	}
	if iresp.Error != nil {
		return nil, errors.Wrap(err, "could not inspect perspective entity file")
	}
	wc.perspectiveHash = iresp.Entity.Hash

	return wc, nil
}

func (wc *WaveCredentials) AddServiceInfo(server *grpc.Server) {
	wc.info = server.GetServiceInfo()

	var uris []string

	// form a list of <package name>/<service name>/<method name> URIs
	for pkg_svc_name, svc_info := range wc.info {
		uri_pkg_svc_name := strings.Replace(pkg_svc_name, ".", "/", -1)
		for _, method_info := range svc_info.Methods {
			// TODO: get rid of this hack
			wc.grpcservice = uri_pkg_svc_name + "/*"
			uris = append(uris, uri_pkg_svc_name+"/"+method_info.Name)
			log.Info("GRPC Resource: ", uris[len(uris)-1])
		}
	}
}

func (wc *WaveCredentials) AddGRPCProofFile(filename string) (ns string, proof []byte, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil, errors.Wrap(err, "could not read designated routing file")
	}

	der := contents
	pblock, _ := pem.Decode(contents)
	if pblock != nil {
		der = pblock.Bytes
	}

	resp, err := wc.wave.VerifyProof(context.Background(), &pb.VerifyProofParams{
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
		log.Info("matches resource? ", s.Resource == wc.grpcservice, " ", s.Resource, " ", wc.grpcservice, " matches permset? ", bytes.Equal(s.GetPermissionSet(), []byte(XBOSPermissionSet)))
		if s.Resource == wc.grpcservice && bytes.Equal(s.GetPermissionSet(), []byte(XBOSPermissionSet)) {
			for _, perm := range s.Permissions {
				log.Info("match perm? ", perm == GRPCServePermission, " ", perm, " ", GRPCServePermission)
				//TODO: need to MATCH the uri here for each of the uris, make sure we prove it
				if perm == GRPCServePermission {
					found = true
					break outer
				}
			}
		}
	}

	if !found {
		return "", nil, fmt.Errorf("designated routing proof does not actually prove xbos:serve_grpc on any namespace")
	}
	wc.namespace = ns
	wc.proof = der

	return ns, der, nil
}

func (wc *WaveCredentials) ServerTransportCredentials() credentials.TransportCredentials {
	return wc
}

func (wc *WaveCredentials) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, errors.New("NOT IMPLEMENTED")
}

func (wc *WaveCredentials) VerifyClientHandshake(nsString string, hdr clientHeader) error {
	log.Info("Server verifying server handshake ", nsString)
	resp, err := wc.wave.VerifySignature(context.Background(), &pb.VerifySignatureParams{
		Signer:    hdr.entityHash,
		Signature: hdr.signature,
		Content:   hdr.proof,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	if resp.Error != nil {
		log.Error(resp.Error.Message)
		return errors.New(resp.Error.Message)
	}

	ns, err := base64.URLEncoding.DecodeString(nsString)
	if err != nil {
		log.Error(err)
		return err
	}

	//Signature ok, verify proof
	presp, err := wc.wave.VerifyProof(context.Background(), &pb.VerifyProofParams{
		ProofDER: hdr.proof,
		Subject:  hdr.entityHash,
		RequiredRTreePolicy: &pb.RTreePolicy{
			Namespace: ns,
			Statements: []*pb.RTreePolicyStatement{
				{
					PermissionSet: []byte(XBOSPermissionSet),
					Permissions:   []string{GRPCCallPermission},
					// grpc_package/ServiceName/* (all methods)
					// grpc_package/ServiceName/Method1 (only method 1)
					Resource: wc.grpcservice, // TODO: replace this with the name, etc of the GRPC service
				},
			},
		},
	})

	if err != nil {
		log.Error(err)
		return err
	}
	if presp.Error != nil {
		log.Error(presp.Error.Message)
		return errors.New(presp.Error.Message)
	}
	if !bytes.Equal(presp.Result.Subject, hdr.entityHash) {
		log.Error("proof valid")
		return errors.New("proof valid but for a different entity")
	}
	return nil
}

func (wc *WaveCredentials) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	//Generate TLS certificate
	cert, cert2 := genCert()
	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}
	conn := tls.Server(rawConn, &tlsConfig)

	err := conn.Handshake()
	if err != nil {
		rawConn.Close()
		return nil, nil, err
	}
	namespace := make([]byte, 34)
	_, err = io.ReadFull(conn, namespace)
	log.Debug("namespace ", namespace)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not generate header: %v", err)
	}

	hdr, err := wc.ReadPeerHeader(conn)
	if err != nil {
		rawConn.Close()
		return nil, nil, errors.Wrap(err, "Could not read server header")
	}

	err = wc.VerifyClientHandshake(wc.namespace, hdr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not verify server handshake")
	}

	//	if true {
	//		log.Debug("server read proof size")
	//		proofSizeBA := make([]byte, 4)
	//		_, err = io.ReadFull(conn, proofSizeBA)
	//		if err != nil {
	//			rawConn.Close()
	//			return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	//		}
	//		proofSize := binary.LittleEndian.Uint32(proofSizeBA)
	//		if proofSize > 10*1024*1024 {
	//			rawConn.Close()
	//			return nil, nil, fmt.Errorf("bad proof")
	//		}
	//		log.Debug("server read proof")
	//		proof := make([]byte, proofSize)
	//		_, err = io.ReadFull(conn, proof)
	//		if err != nil {
	//			rawConn.Close()
	//			return nil, nil, fmt.Errorf("could not read proof: %v\n", err)
	//		}
	//
	//		// verify proof
	//
	//	}

	header, err := wc.GeneratePeerHeader(namespace, cert2.Signature)
	if err != nil {
		rawConn.Close()
		return nil, nil, fmt.Errorf("could not generate header: %v", err)
	}
	_, err = conn.Write(header)
	if err != nil {
		rawConn.Close()
		return nil, nil, err
	}

	return conn, nil, nil
}

func (wc *WaveCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	}
}
func (wc *WaveCredentials) Clone() credentials.TransportCredentials {
	return &WaveCredentials{
		perspective:     wc.perspective,
		perspectiveHash: wc.perspectiveHash,
		namespace:       wc.namespace,
		wave:            wc.wave,
	}
}

func (wc *WaveCredentials) OverrideServerName(name string) error {
	return nil
}

//A 34 byte multihash
func (wc *WaveCredentials) GeneratePeerHeader(ns []byte, cert []byte) ([]byte, error) {
	hdr := bytes.Buffer{}
	if len(wc.perspectiveHash) != 34 {
		panic(wc.perspectiveHash)
	}
	//First: 34 byte entity hash
	hdr.Write(wc.perspectiveHash)
	//Second: signature of cert
	sigresp, err := wc.wave.Sign(context.Background(), &pb.SignParams{
		Perspective: wc.perspective,
		Content:     cert,
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
	binary.LittleEndian.PutUint32(prooflen, uint32(len(wc.proof)))
	hdr.Write(prooflen)
	hdr.Write(wc.proof)
	return hdr.Bytes(), nil
}

func genCert() (tls.Certificate, *x509.Certificate) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
		panic(err)
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "wavemq-dr",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
		panic(err)
	}
	x509cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	keybytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	certbytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	cert, err := tls.X509KeyPair(certbytes, keybytes)
	if err != nil {
		panic(err)
	}
	return cert, x509cert
}

func (wc *WaveCredentials) ReadPeerHeader(conn io.Reader) (clientHeader, error) {
	var (
		err error
		hdr clientHeader
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

type clientHeader struct {
	entityHash []byte
	signature  []byte
	proof      []byte
}
