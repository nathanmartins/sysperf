package collectors

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	reParens = regexp.MustCompile(`\((.*)\)`)
)

type MemInfoCollector struct{}

func (c MemInfoCollector) Collect(ch chan<- prometheus.Metric) {
	memInfo, err := getMemInfo()
	if err != nil {
		log.Printf("error while getting meminfo: %s\n", err)
	}

	var metricType prometheus.ValueType

	for k, v := range memInfo {
		if strings.HasSuffix(k, "_total") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("sysperf", "memory", strings.ToLower(k)),
				fmt.Sprintf("Memory information field %s.", k),
				nil, nil,
			),
			metricType, v,
		)
	}
}

func (c MemInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range getDescriptions() {
		ch <- desc
	}
}

func getMemInfo() (map[string]float64, error) {

	file, err := os.Open(filepath.Join(procfs.DefaultMountPoint, "meminfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMemInfo(file)
}

func getDescriptions() []*prometheus.Desc {
	memInfo, err := getMemInfo()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("error while getting meminfo: %s\n", err))
	}

	descriptions := make([]*prometheus.Desc, 0)

	for k, _ := range memInfo {
		descriptions = append(descriptions, prometheus.NewDesc(
			prometheus.BuildFQName("sysperf", "memory", strings.ToLower(k)),
			fmt.Sprintf("Memory information field %s.", k),
			nil, nil,
		))
	}

	return descriptions
}

func parseMemInfo(r io.Reader) (map[string]float64, error) {
	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		// Workaround for empty lines occasionally occur in CentOS 6.2 kernel 3.10.90.
		if len(parts) == 0 {
			continue
		}
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %w", err)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = reParens.ReplaceAllString(key, "_${1}")
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
			key = key + "_bytes"
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		memInfo[key] = fv
	}

	return memInfo, scanner.Err()
}
