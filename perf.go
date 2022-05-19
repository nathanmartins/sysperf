package main

import (
	"os"
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
