up:
	docker-compose -f docker-compose.yml up --build
test:
	go test -v ./...
generate proto:
	protoc -I proto ./proto/app/*.proto --go_out=protos/gen/go --go_opt=paths=source_relative --go-grpc_out=protos/gen/go
#go to sqlc:
#	cd internal/repository/sqlc
#generate sqlcode:
#	sqlc generate
#sqlc: go to sqlc, generate sqlcode
#
