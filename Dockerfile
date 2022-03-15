FROM archlinux:base-devel-20220313.0.50300

RUN sudo pacman -Syy --noconfirm bcc bcc-tools python-bcc bc wget

WORKDIR /usr/sbin/

COPY fetch-linux-headers.sh .

RUN fetch-linux-headers.sh
COPY . .

WORKDIR /

