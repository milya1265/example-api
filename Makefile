up:
	docker-compose -f docker-compose.yml up --build
test:
	go test -v ./...
generate grpc:
	protoc -I proto ./proto/app/app.proto --go_out=protos/gen/go --go_opt=paths=source_relative --go-grpc_out=protos/gen/go

