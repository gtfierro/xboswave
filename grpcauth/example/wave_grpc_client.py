import grpc
import threading
import socket
import select
import time
import sys
import struct
import eapi_pb2
import eapi_pb2_grpc
import base64
import grpcserver_pb2
import grpcserver_pb2_grpc
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from tlslite import TLSConnection

XBOS_PERMSET = base64.urlsafe_b64decode("GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA==")

class WAVEGRPCClient:
    def __init__(self, address_tuple, namespace, entityfile, proof_file='clientproof.pem', waved='localhost:410'):
        self.address_tuple = address_tuple
        self.ns = namespace
        self.nsbytes = base64.urlsafe_b64decode(self.ns)
        self.entityfile = open(entityfile, 'rb').read()
        self.perspective = eapi_pb2.Perspective(
            entitySecret=eapi_pb2.EntitySecret(
                DER=self.entityfile
            )
        )
        self._listen_address = None
        self._ready = threading.Event()

        self.wave_channel = grpc.insecure_channel(waved)
        self.wave_client = eapi_pb2_grpc.WAVEStub(self.wave_channel)
        resp = self.wave_client.Inspect(eapi_pb2.InspectParams(
            content=self.entityfile,
        ))
        self.entityhash = resp.entity.hash


        self.proof_file = open('clientproof.pem','rb').read()
        resp = self.wave_client.VerifyProof(eapi_pb2.VerifyProofParams(
            proofDER=self.proof_file,
        ))
        self.sigresp = self.wave_client.Sign(eapi_pb2.SignParams(
            perspective=self.perspective,
            content=self.proof_file,
        ))

        # setup server
        server_thread = threading.Thread(target=self.get_client_connection)
        server_thread.start()

    def setup_connection(self):
        hdr = self.generate_peer_header()

        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        sock.connect(self.address_tuple)
        self.upstream_connection = TLSConnection(sock)
        hs = self.upstream_connection.handshakeClientCert()
        self.upstream_connection.write(self.nsbytes)
        self.upstream_connection.write(hdr)
        invalid = self.read_peer_header(self.upstream_connection)
        if invalid.message != '':
            raise Exception("GRPC Server sent invalid header or proof {0}".format(invalid))

    @property
    def listen_address(self):
        self._ready.wait()
        return "{0}:{1}".format(*self._listen_address)

    def get_client_connection(self):
        listen_port = 5005
        while True:
            listen_address = ('localhost', listen_port)
            server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            try:
                server.bind(listen_address)
                self._listen_address = listen_address
                self._ready.set()
                break
            except Exception as e:
                print("Failed to listen on {0}".format(listen_address), e)
                time.sleep(1)
                listen_port += 1
        print("Listening on {0}".format(listen_address))
        server.listen(10)

        while True:
            client_socket, addr = server.accept()
            # reconnect to the GRPC server on each call
            self.setup_connection()

            # print out the local connection information
            print("new client call")

            # start a thread to talk to the remote host
            proxy_thread = threading.Thread(target=self.handle_client,
                                            args=(client_socket,))

            proxy_thread.start()

    def handle_client(self, client_socket):
        while True:
            try:
                local_buffer = receive_from(client_socket)
                if len(local_buffer):
                    self.upstream_connection.send(local_buffer)
                # receive back the response
                remote_buffer = receive_from(self.upstream_connection)
                if len(remote_buffer):
                    # send the response to the local socket
                    client_socket.send(remote_buffer)
                # if no more data on the either side, close the connections
                if not len(local_buffer) or not len(remote_buffer):
                    print("Done with call")
                    break
            finally:
                client_socket.close()
                self.upstream_connection.close()

    def generate_peer_header(self):
        buf = bytes()
        buf += self.entityhash
        buf += struct.pack('<H', len(self.sigresp.signature))
        buf += self.sigresp.signature
        buf += struct.pack('<I', len(self.proof_file))
        buf += self.proof_file
        return buf

    def read_peer_header(self, conn):
        entityhash = conn.read(max=34,min=34)
        sigsize = struct.unpack('<H', conn.read(max=2,min=2))[0]
        signature = conn.read(max=sigsize, min=sigsize)
        proofsize = struct.unpack('<I', conn.read(max=4, min=4))[0]
        proof = conn.read(max=proofsize, min=proofsize)
        #TODO verify this
        # TODO: need peer certificate
        cert = self.upstream_connection.session.serverCertChain.x509List[0].bytes

        c = x509.load_der_x509_certificate(cert, default_backend())

        vresp = self.wave_client.VerifySignature(eapi_pb2.VerifySignatureParams(
            signer=entityhash,
            signature=signature,
            content=c.signature,
        ))
        if vresp.error.message != "":
            return vresp.error

        proofresp = self.wave_client.VerifyProof(eapi_pb2.VerifyProofParams(
            proofDER=proof,
            subject=entityhash,
            requiredRTreePolicy=eapi_pb2.RTreePolicy(
                namespace=self.nsbytes,
                statements=[
                    eapi_pb2.RTreePolicyStatement(
                        permissionSet=XBOS_PERMSET,
                        permissions=["serve_grpc"],
                        resource="xbospb/Test/*"
                    ),
                ],
            )
        ))
        if proofresp.result == None:
            return "no proof"
        return proofresp.error

def receive_from(connection):
    buffer = b""
    try:
        connection.settimeout(1)
        while True:
            data = connection.recv(4096)
            if not data:
                break
            buffer += data
    except:
        pass
    return buffer

client = WAVEGRPCClient( ('localhost', 7373),  "GyBHxjkpzmGxXk9qgJW6AJHCXleNifvhgusCs0v1MLFWJg==", "client.ent")

channel = grpc.insecure_channel(client.listen_address)
stub = grpcserver_pb2_grpc.TestStub(channel)
while True:
    t = time.time()
    resp = stub.TestUnary(grpcserver_pb2.TestParams(x="hello 123456"))
    t2 = time.time()
    print("Elapsed: {0}. Threads active: {1}. Got Response {2}".format(t2-t, threading.active_count(), resp))
