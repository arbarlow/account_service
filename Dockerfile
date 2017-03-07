FROM alpine:latest

COPY . /src/github.com/lileio/account_service
RUN apk add --no-cache ca-certificates
ADD build/account_service /bin
CMD ["account_service", "server"]
