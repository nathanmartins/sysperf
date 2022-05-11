package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

func SampleCPUSaturation(interval time.Duration) error {
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

	err := SendMetric(sample, "cpu_saturation")
	if err != nil {
		log.Println("failed to send cpu_saturation metric, API offline?")
	}

	return nil
}
