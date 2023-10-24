from golang:1.21-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

FROM alpine:latest

COPY --from=builder /build/main /bin/main

EXPOSE "8080"

ENTRYPOINT ["/bin/main"]
