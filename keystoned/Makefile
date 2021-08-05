proto: proto/keystone.proto proto/keystone.pb.go

proto/keystone.pb.go:
	protoc --go_out=plugins=grpc:. ./proto/*.proto

clean:
	rm -rf proto/*.pb.go
