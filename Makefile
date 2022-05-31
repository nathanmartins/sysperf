build:
	CGO_ENABLED=0 go build -ldflags="-s -w"

transfer:
	ssh root@161.35.0.6 pkill -9 sysperf || true
	scp sysperf root@161.35.0.6:/usr/local/bin/
	ssh root@161.35.0.6 sysperf

.PHONY: all
all:
	$(MAKE) build
	$(MAKE) transfer