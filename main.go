package main

import "github.com/nathanmartins/sysperf/rpc"

func main() {
	a := &rpc.Rpc{}
	a.SendMetric("Hello RPC World")
}
