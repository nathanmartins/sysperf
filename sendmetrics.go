package main

import (
	"log"
	"net/rpc"
)

type Rpc struct{}

type Args struct {
	Message string
	Context string
}

type Reply struct {
	FinalMessage string
}

func main() {
	a := &Rpc{}
	a.SendMetric("Hello RPC World")
}

func (r *Rpc) SendMetric(message string) {
	client, err := rpc.DialHTTPPath("tcp", "127.0.0.1:8080", "/rpc")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	a := Args{Message: message, Context: "sysperf"}
	var repl Reply

	replyCall := client.Call("Rpc.SendMetric", a, &repl)
	log.Printf("%+v", replyCall)
	log.Printf("%+v", repl)
}
