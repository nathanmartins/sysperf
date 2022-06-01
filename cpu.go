package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

type CPUUsage struct {
	Usage float64 `json:"usage"`
	Busy  float64 `json:"busy"`
	Total float64 `json:"total"`
}

type CPULatency struct {
	Command         string  `json:"command"`
	TimeSpentOnCPU  float64 `json:"time-spent-on-cpu"`
	RunQueueLatency float64 `json:"run-queue-latency"`
}

func SampleCPUUsage(interval time.Duration) (CPUUsage, error) {
	idle0, total0 := getCPUSample()
	time.Sleep(interval)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	sample := CPUUsage{
		Usage: 100 * (totalTicks - idleTicks) / totalTicks,
		Busy:  totalTicks - idleTicks,
		Total: totalTicks,
	}

	return sample, nil
}

func SampleCPULatency() ([]CPULatency, error) {

	var samples []CPULatency

	files, err := ioutil.ReadDir("/proc/")
	if err != nil {
		return samples, err
	}

	for _, f := range files {

		_, cErr := strconv.Atoi(f.Name())

		if cErr == nil {

			comm, _ := os.ReadFile(fmt.Sprintf("/proc/%s/comm", f.Name()))
			fullB, _ := os.ReadFile(fmt.Sprintf("/proc/%s/schedstat", f.Name()))

			latencies := strings.Split(string(fullB), " ")

			intoLatencies := make([]int, len(latencies))

			for i, s := range latencies {
				intoLatencies[i], _ = strconv.Atoi(s)
			}

			if len(latencies) == 3 {
				c := CPULatency{
					Command:         string(comm),
					TimeSpentOnCPU:  float64(intoLatencies[0]) / float64(time.Millisecond),
					RunQueueLatency: float64(intoLatencies[1]) / float64(time.Millisecond),
				}
				samples = append(samples, c)
			}
		}
	}

	return samples, err
}
