# compile the service
FROM golang:1.12 as builder
COPY . /transcribe-service
WORKDIR /transcribe-service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o transcribe-service

# use scratch for minimal image size
FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /transcribe-service .
ENTRYPOINT ["./transcribe-service"]