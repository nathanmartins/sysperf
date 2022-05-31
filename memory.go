package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type MemoryLatency struct {
	Command string  `json:"command"`
	SizeKb  float64 `json:"size-kb"`
}

func SampleMemoryLatency() ([]MemoryLatency, error) {

	var samples []MemoryLatency

	files, err := ioutil.ReadDir("/proc/")
	if err != nil {
		return samples, err
	}

	for _, f := range files {

		_, cErr := strconv.Atoi(f.Name())

		if cErr == nil {

			comm, _ := os.ReadFile(fmt.Sprintf("/proc/%s/comm", f.Name()))
			fullB, _ := os.ReadFile(fmt.Sprintf("/proc/%s/status", f.Name()))
			reader := bytes.NewReader(fullB)

			scanner := bufio.NewScanner(reader)
			scanner.Split(bufio.ScanLines)
			var readerLines []string

			for scanner.Scan() {
				readerLines = append(readerLines, scanner.Text())
			}

			for _, eachLine := range readerLines {
				if strings.Contains(eachLine, "VmSize") {
					size, _ := strconv.Atoi(strings.Fields(eachLine)[1])

					command := string(comm)
					command = strings.TrimSuffix(command, "\n")

					c := MemoryLatency{
						Command: command,
						SizeKb:  float64(size),
					}
					samples = append(samples, c)
				}
			}

		}
	}

	return samples, err
}
