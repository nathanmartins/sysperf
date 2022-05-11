package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time" // or "runtime"
)

func cleanup() {
	fmt.Println("cleanup")
}

const VERSION = "0.0.1"

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	for {
		log.Printf("starting sysperf agent version: %s\n", VERSION)
		example()
		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}
}
