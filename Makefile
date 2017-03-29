proto:
	protoc -I $$GOPATH/src/ -I . account_service.proto --lile-server_out=. --go_out=plugins=grpc:$$GOPATH/src

run:
	go run account_service/main.go

.PHONY: test
test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem -benchtime 10s ./...

docker:
	GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o build/account_service ./account_service
	docker build . -t lileio/account_service:`git rev-parse --short HEAD`
	@echo lileio/account_service:`git rev-parse --short HEAD`

