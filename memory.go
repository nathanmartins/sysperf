package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MemoryUsage struct {
	Command string  `json:"command"`
	SizeKb  float64 `json:"size-kb"`
}

type MemorySaturation struct {
	SwapUsage   float64 `json:"swap-usage"`
	OOMKillings float64 `json:"oom_killings"`
}

func SampleMemoryUsage() ([]MemoryUsage, error) {

	var samples []MemoryUsage

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

					c := MemoryUsage{
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

func SampleMemorySaturation() (MemorySaturation, error) {
	var ms MemorySaturation
	var found int

	memBytes, err := os.ReadFile("/proc/meminfo")

	if err != nil {
		return ms, err
	}

	temp := strings.Split(string(memBytes), "\n")

	for _, item := range temp {
		if strings.Contains(item, "SwapTotal") {
			re := regexp.MustCompile("[0-9]+")
			found, err = strconv.Atoi(strings.Join(re.FindAllString(item, -1), ""))
			if err != nil {
				return ms, err
			}
			ms.SwapUsage = float64(found)
		}

	}

	return ms, nil
}

// dmesg -T | egrep -i 'killed process'
