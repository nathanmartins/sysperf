FROM alpine:3.15
RUN apk add --update \
    bash \
    bc \
    build-base \
    bison \
    flex \
    curl \
    elfutils-dev \
    linux-headers \
    make \
    wget \
    openssl-dev \
    tar \
    gzip \
    python3

# Setting up some python things
ENV PYTHONUNBUFFERED=1
RUN  ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools

WORKDIR /usr/sbin/

# Cache long step
COPY fetch-linux-headers.sh .
RUN fetch-linux-headers.sh

COPY . .
RUN apk add --update bcc-tools bcc-doc

WORKDIR /
