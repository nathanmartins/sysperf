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
const USER_ID = "0fa48e83-1b1a-4822-98e5-807bf06b2b63"

func main() {

	log.Printf("starting sysperf agent version: %s\n", VERSION)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	for {

		go func() {
			err := SampleCPUSaturation(3 * time.Second)
			if err != nil {
				cleanup()
				log.Fatal(err)
			}
		}()

		go func() {
			err := SampleCPULatency(1)
			if err != nil {
				cleanup()
				log.Fatal(err)
			}
		}()

		time.Sleep(10 * time.Second) // Agent run interval
	}
}
