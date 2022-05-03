# sysperf
```shell
sudo docker run -it --cap-add SYS_ADMIN --privileged sysperf bash

perf record -F 99 -a -g -- sleep 30

perf data convert --to-json
```
