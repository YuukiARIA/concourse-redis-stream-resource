FROM golang:1.12 AS builder

COPY . /workspace
WORKDIR /workspace

ENV CGO_ENABLED 0

RUN go build -o bin/check cmd/check/main.go \
 && go build -o bin/in    cmd/in/main.go \
 && go build -o bin/out   cmd/out/main.go

RUN go test -v ./...

FROM busybox

COPY --from=builder /workspace/bin /opt/resource
