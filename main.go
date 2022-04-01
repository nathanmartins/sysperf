package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathanmartins/sysperf/custom_ebpf"

	"github.com/cilium/ebpf/rlimit"
)

func init() {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	errors := make(chan error, 1)
	go func() {
		errors <- custom_ebpf.RunCGroup(ctx)
	}()

	<-stop
	cancel()

	if err := <-errors; err != nil {
		log.Fatalf("error linking ebpf program: %v", err)
	}
}
