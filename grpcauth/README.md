# GRPC Authorization Layer

For a complete example, see the `example/` directory and `setup.sh` script.

### Overview

This package provides WAVE-based authentication + authorization for GRPC services.
All GRPC connections happen over TLS (currently using self-signed certs)
Services and clients possess WAVE entities.
Services provide a proof during the handshake that they are authorized to provide a GRPC service on a given namespace.
Clients provide a proof during the handshake that they are authorized to call that GRPC service on that namespace.

This is indicated with the serve_grpc` and `call_grpc` permissions on the XBOS `GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA==` permission set.

```
# grant to server to serve all methods
wv rtgrant --attester namespace.ent \
           --subject service.ent \
           GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA==:serve_grpc@namespace.ent/\
           <package name>/<service name>/*

# grant to client to call all methods
wv rtgrant --attester namespace.ent \
           --subject client.ent \
           GyC5wUUGKON6uC4gxuH6TpzU9vvuKHGeJa1jUr4G-j_NbA==:call_grpc@namespace.ent/\
           <package name>/<service name>/*
```

The implementation is adapted from [https://github.com/immesys/wavemq](https://github.com/immesys/wavemq)

### TODO

- [ ] use hash of WAVE entity to bootstrap TLS connection
    - [ ] avoids self-signed TLS, which we have now
- [ ] additional app-specific metadata in WAVE proofs
    - time bound on historical timeseries retrieval
    - limiting which streams can be returned
- [ ] build proof automatically (easy)
- [X] client provides proof of authorization
    - am I allowed to call this method?


### Usage

Server (simplified):
```go
import (
    "github.com/gtfierro/xboswave/grpcauth"
)

func main() {

    // setup WAVE perspective and create credentials object
    server_perspective := loadPerspective("service.ent")
    serverwavecreds, _ := grpcauth.NewServerCredentials(server_perspective, "localhost:410")

    // register generic GRPC server with service
    xbospb.RegisterTestServer(grpcServer, testserver{})

    // inject a description of the service into the credentials object
    serverwavecreds.AddServiceInfo(grpcServer)

    // add the service authorization proof (see below)
    serverwavecreds.AddGRPCProofFile("serviceproof.pem")

    // serve GRPC
    grpcServer.Serve(l)
}
```

Client (simplified):

```go
import (
    "github.com/gtfierro/xboswave/grpcauth"
)

func main() {

    // setup WAVE perspective and create CLIENT credentials object
    client_perspective := loadPerspective("client.ent")

    // namespace, GRPC service descriptor
    clientcred, err := grpcauth.NewClientCredentials(client_perspective, "localhost:410", "GyBHxjkpzmGxXk9qgJW6AJHCXleNifvhgusCs0v1MLFWJg==", "xbospb/Test/*")
    if err != nil {
        log.Fatal(err)
    }

    // add credentials object to normal GRPC dial
    clientconn, err := grpc.Dial("localhost:7373", grpc.WithTransportCredentials(clientcred), grpc.FailOnNonTempDialError(true), grpc.WithBlock(), grpc.WithTimeout(30*time.Second))
    if err != nil {
        log.Fatal(err)
    }
    testclient := xbospb.NewTestClient(clientconn)
}
```
