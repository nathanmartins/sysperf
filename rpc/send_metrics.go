package rpc

import (
	"log"
	"net/rpc"
)

type Rpc struct{}

type Event struct {
	Data interface{}
	Name string
}

type Reply struct {
	FinalMessage string
}

var client *rpc.Client

func getClient() (*rpc.Client, error) {
	if client == nil {
		c, err := rpc.DialHTTPPath("tcp", "127.0.0.1:8080", "/rpc")
		if err != nil {
			log.Println("error dialing:", err)
			return c, err
		}
		client = c
	}

	return client, nil
}

func (r *Rpc) SendMetric(payload interface{}, metricName string) error {
	c, err := getClient()
	if err != nil {
		return err
	}
	e := Event{Data: payload, Name: metricName}
	repl := new(Reply)

	replyCall := c.Call("Rpc.SendMetric", e, repl)
	log.Printf("%+v", replyCall)
	log.Printf("%+v", repl)
	return nil
}
