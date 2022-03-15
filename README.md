# sysperf - eBPF based metric-collector/analyzer 


## How to build: 

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
