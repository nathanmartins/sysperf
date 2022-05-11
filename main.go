package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type PerfFile struct {
	LinuxPerfJSONVersion int `json:"linux-perf-json-version"`
	Headers              struct {
		HeaderVersion int       `json:"header-version"`
		CapturedOn    time.Time `json:"captured-on"`
		DataOffset    int       `json:"data-offset"`
		DataSize      int       `json:"data-size"`
		FeatOffset    int       `json:"feat-offset"`
		Hostname      string    `json:"hostname"`
		OsRelease     string    `json:"os-release"`
		Arch          string    `json:"arch"`
		CPUDesc       string    `json:"cpu-desc"`
		Cpuid         string    `json:"cpuid"`
		NrcpusOnline  int       `json:"nrcpus-online"`
		NrcpusAvail   int       `json:"nrcpus-avail"`
		PerfVersion   string    `json:"perf-version"`
		Cmdline       []string  `json:"cmdline"`
	} `json:"headers"`
	Samples []struct {
		Timestamp int64  `json:"timestamp"`
		Pid       int    `json:"pid"`
		Tid       int    `json:"tid"`
		Comm      string `json:"comm"`
		Callchain []struct {
			IP     string `json:"ip"`
			Symbol string `json:"symbol,omitempty"`
			Dso    string `json:"dso,omitempty"`
		} `json:"callchain"`
	} `json:"samples"`
}

func main() {

	CleanUp()

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			err := RunPerf()
			if err != nil {
				CleanUp()
				log.Fatal(err)
			}

			err = ProcessPerf()
			if err != nil {
				CleanUp()
				log.Fatal(err)
			}
		case <-quit:
			CleanUp()
			ticker.Stop()
			return
		}
	}

}

func RunPerf() error {

	err := exec.Command("perf", "record", "-F", "99", "-a", "--", "sleep", "5").Run()

	if err != nil {
		return fmt.Errorf("failed to run")
	}
	err = exec.Command("perf", "data", "convert", "--to-json", "perf.json").Run()

	if err != nil {
		return fmt.Errorf("failed to convert: %s", err)
	}

	return nil
}

func ProcessPerf() error {
	var perfFile PerfFile

	content, err := ioutil.ReadFile("perf.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &perfFile)
	if err != nil {
		return err
	}

	err = SendMetric(perfFile, "all_cpus")
	if err != nil {
		return err
	}

	CleanUp()

	return nil
}

func CleanUp() {
	p := []string{
		"perf.json",
		"perf.data",
		"perf.data.old",
	}

	for _, pp := range p {
		_ = os.Remove(pp)
	}
}

func SendMetric(event interface{}, eventName string) error {
	body, err := json.Marshal(event)
	if err != nil {
		log.Println("error marshaling event:", err)
		return err
	}

	sysPerfUrl, found := os.LookupEnv("SERVER_URL")

	if !found {
		log.Println("WARNING: SERVER_URL env var not found, using default: localhost:8080")
		sysPerfUrl = "http://localhost:8080/metric"
	}

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, sysPerfUrl, payload)
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
