run:
	go run main.go

proto:
	protoc -I account/ account/account.proto --go_out=plugins=grpc:account

test:
	go test -v ./...

benchmark:
	go test -bench=./... -benchmem -benchtime 10s

docker:
	GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o account_service .
	docker build . -t lileio/account_service:latest

