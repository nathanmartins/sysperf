package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

type CPUSaturation struct {
	Usage float64 `json:"usage"`
	Busy  float64 `json:"busy"`
	Total float64 `json:"total"`
}

func SampleCPUSaturation(interval time.Duration) (CPUSaturation, error) {
	idle0, total0 := getCPUSample()
	time.Sleep(interval)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	sample := CPUSaturation{
		Usage: 100 * (totalTicks - idleTicks) / totalTicks,
		Busy:  totalTicks - idleTicks,
		Total: totalTicks,
	}

	return sample, nil
}

func SampleCPULatency(interval time.Duration) error {

	_ = FileCleanUp("cpu-lat*")

	duration := fmt.Sprintf("%d", interval)

	command := []string{"sched", "record", "-o", "cpu-lat.data", "--", "sleep", duration}
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
