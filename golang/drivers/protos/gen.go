package mqpb

//ignore go:generate protoc -I/usr/local/include -I. -I$GOPATH/src/github.com/immesys/wave/eapi/pb -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. wavemq.proto

//go:generate protoc -I/usr/local/include -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. wavemq.proto eapi.proto
