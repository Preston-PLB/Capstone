FROM golang:1.21-alpine AS builder

WORKDIR /build/service

COPY service/go.mod service/go.sum ./
COPY ./ui ../ui
COPY ./libs ../libs

RUN go mod download

COPY ./service .

RUN GOEXPERIMENT=loopvar go build -o main

FROM amazonlinux:2023

COPY docker/resolv.conf /etc/resolv.conf
COPY --from=builder /build/service/main /bin/main
RUN mkdir -p /etc/capstone
COPY secrets/config.yaml /etc/capstone

EXPOSE "8080"

ENTRYPOINT ["/bin/main"]
