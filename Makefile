proto:
	protoc -I account/ account/account.proto --go_out=plugins=grpc:account

test:
	go test -v ./...

benchmark:
	go test -bench=./... -benchmem -benchtime 10s
