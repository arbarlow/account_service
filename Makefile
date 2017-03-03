proto:
	protoc -I ../image_service/image_service -I . account_service.proto --go_out=plugins=grpc:$$GOPATH/src

test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem -benchtime 10s ./...

docker:
	GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o build/account_service ./account_service
	docker build . -t lileio/account_service:latest

