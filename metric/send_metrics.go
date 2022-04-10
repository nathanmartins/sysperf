package metric

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func SendMetric(event interface{}, eventName string) error {
	log.Printf("event received: %+v\n", event)
	body, err := json.Marshal(event)
	if err != nil {
		log.Println("error marshaling event:", err)
		return err
	}

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/metric", payload)
	if err != nil {
		log.Println("new request error:", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Event-Name", eventName)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println("sysperf api request error:", err)
		return err
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("sysperf api body read error:", err)
		return err
	}

	log.Printf("event succesfully saved: %s\n", string(response))

	return nil
}
