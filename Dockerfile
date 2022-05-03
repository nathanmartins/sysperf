FROM ubuntu:jammy

RUN apt-get update && apt-get -y install linux-tools-common linux-tools-generic linux-tools-`uname -r` golang-go git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build

CMD [ "/app/sysperf" ]
