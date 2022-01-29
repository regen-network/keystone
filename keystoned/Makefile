INSTALL_DIR="$(shell go env GOPATH)/bin"

# build is first target so that lone 'make' calls 'make build'
build: proto *.go
	go build .

install: ./keystoned
	cp keystoned $(INSTALL_DIR)

proto: proto/keystone.pb.go proto/keystone_grpc.pb.go

proto/keystone.pb.go: proto/keystone.proto
	protoc --go_out=. --go-grpc_out=. ./proto/*.proto

clean:
	rm -rf proto/*.pb.go
	rm keystoned
