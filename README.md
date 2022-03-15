# sysperf - eBPF based metric-collector/analyzer 

## Supports:

- [x] Networking: IPv4 and IPv6 tcp connections
- [ ] CPU Usage
- [ ] Memory usage 


## How to build: 

This is an [Alpine Linux](https://alpinelinux.org/) based Docker image, in which we install the [BPF Compiler Collection (BCC)](https://github.com/iovisor/bcc) 
and the necessary kernel headers.

```shell
docker build -t nathanmartins/sysperf . 
```

## How to run: 

In the main shell
```shell
docker run --rm --name sysperf --workdir /code --privileged -it -v $(pwd):/code nathanmartins/sysperf bash 
./main.py
```

In a separate shell:
```shell
docker exec -it sysperf wget google.com -O /dev/null
```
