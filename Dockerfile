FROM alpine:latest

RUN apk add --no-cache ca-certificates
ADD build/account_service /bin
CMD ["account_service", "server"]
