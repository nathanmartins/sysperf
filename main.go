package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	var perfFile PerfFile

	content, err := ioutil.ReadFile("perf.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &perfFile)
	if err != nil {
		log.Fatal(err)
	}

	commandMapping := make(map[int]string)

	for _, sample := range perfFile.Samples {

		if _, ok := commandMapping[sample.Pid]; !ok {
			commandMapping[sample.Pid] = sample.Comm
		}
	}

	fmt.Println(commandMapping)

}
