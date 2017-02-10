FROM alpine:latest

RUN apk add --no-cache ca-certificates
ADD account_service /
CMD ["/account_service"]
