# sysperf - eBPF based metric-collector/analyzer 

## Supports:

- [ ] Networking: IPv4 and IPv6 tcp connections
- [ ] CPU Usage
- [ ] Memory usage

## Dependencies

You'll need a recent LLVM installation that can target bpfeb and bpfel.

Aside from LLVM, you'll need at least go 1.17 and your PATH variable including the default location that `go install` writes to.

After installing go, install Cilium's bpf2go utility:

```bash
go install github.com/cilium/ebpf/cmd/bpf2go@latest
```

## Building

Once all dependencies are installed, run `make`. 
You can test the output binary on any of the Vagrant machines by running `sudo ./sysperf` in the `/vagrant` directory.


## Docker

## How to build:

This is an [Alpine Linux](https://alpinelinux.org/) based Docker image.

```shell
docker build -t nathanmartins/sysperf . 
```

## How to run:

In the main shell
```shell
docker run --rm --name sysperf --privileged -it nathanmartins/sysperf main.py
```
