FROM ubuntu:jammy

RUN apt-get update && apt-get -y install linux-tools-common linux-tools-generic linux-tools-`uname -r`
