package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
		return fmt.Errorf("sysperf api error: %s", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("sysperf api error: %s", err)
	}

	log.Printf("API response: %s\n", response)

	return nil
}

func CPULatency() error {

	_ = FileCleanUp("cpu-lat*")

	command := []string{"sched", "record", "-o", "cpu-lat.data", "--", "sleep", "1"}
	err := exec.Command("perf", command...).Run()

	if err != nil {
		_ = FileCleanUp("cpu-lat*")
		return fmt.Errorf("failed to run CPULatency: %s", err)
	}
	err = exec.Command("perf", "data", "-i", "cpu-lat.data", "convert", "--to-json", "cpu-lat.json").Run()

	if err != nil {
		_ = FileCleanUp("cpu-lat*")
		return fmt.Errorf("failed to convert CPULatency: %s", err)
	}

	var perfFile PerfFile

	content, err := ioutil.ReadFile("cpu-lat.json")
	if err != nil {
		_ = FileCleanUp("cpu-lat*")
		return err
	}

	err = json.Unmarshal(content, &perfFile)
	if err != nil {
		_ = FileCleanUp("cpu-lat*")
		return err
	}

	err = SendMetric(perfFile, "cpu_latency")
	if err != nil {
		_ = FileCleanUp("cpu-lat*")
		log.Println("failed to send cpu_latency metric, API offline?")
		//return err
	}

	return nil
}

func FileCleanUp(pattern string) error {

	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			return err
		}
	}

	return err
}
