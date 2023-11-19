from golang:1.21-alpine as builder

WORKDIR /build/service

COPY service/go.mod service/go.sum ./
COPY ./ui ../ui
COPY ./libs ../libs

RUN go mod download

COPY ./service .

RUN go build -o main

FROM alpine:latest

COPY --from=builder /build/service/main /bin/main
RUN mkdir -p /etc/capstone
COPY secrets/config.yaml /etc/capstone

EXPOSE "8080"

ENTRYPOINT ["/bin/main"]
