# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

RUN apk add --no-cache clang make llvm
RUN go install github.com/cilium/ebpf/cmd/bpf2go@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN make

CMD [ "/app/sysperf" ]