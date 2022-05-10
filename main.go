package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println("running perf")

	err := exec.Command("perf", "record", "-F", "99", "-a", "-g", "--", "sleep", "5").Run()

	if err != nil {
		return fmt.Errorf("failed to run")
	}

	log.Println("done running perf")

	log.Println("converting")
	err = exec.Command("perf", "data", "convert", "--to-json", "perf.json").Run()

	if err != nil {
		return fmt.Errorf("failed to convert: %s", err)
	}

	log.Println("done converting")

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

	commandMapping := make(map[int]string)

	for _, sample := range perfFile.Samples {

		if _, ok := commandMapping[sample.Pid]; !ok {
			commandMapping[sample.Pid] = sample.Comm
		}
	}

	CleanUp()
	fmt.Println(commandMapping)

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
