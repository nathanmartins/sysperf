ip = 147.182.138.244

build:
	CGO_ENABLED=0 go build -ldflags="-s -w"

transfer:
	ssh root@$(ip) pkill -9 sysperf || true
	scp sysperf root@$(ip):/usr/local/bin/
	ssh root@$(ip) sysperf

.PHONY: all
all:
	$(MAKE) build
	$(MAKE) transfer