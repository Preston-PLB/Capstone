#Build Go stuff
FROM golang:1.21-alpine AS builder

RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /build/ui

#Setup libs
COPY ui/go.mod ui/go.sum ./
COPY ./service ../service
COPY ./libs ../libs

RUN go mod download

COPY ui/ .

RUN rm **/*_templ.go; templ generate -path ./templates
RUN go build -o main

#Build NPM stuff
FROM node:18-alpine AS node-builder

WORKDIR /build

COPY ui/ .

RUN npm install
RUN npx tailwindcss -i static/index.css -o dist/output.css

#Final Contianer
FROM alpine:latest

COPY --from=node-builder /build/dist /bin/dist
COPY --from=builder /build/ui/main /bin/main
RUN mkdir -p /etc/capstone
COPY secrets/config.yaml /etc/capstone

EXPOSE "8080"

ENTRYPOINT ["/bin/main"]
