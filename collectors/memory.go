package collectors

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	reParens = regexp.MustCompile(`\((.*)\)`)
)

type MeminfoCollector struct {
	logger log.Logger
}

func (m *MeminfoCollector) Describe(descs chan<- *prometheus.Desc) {
	variableLabels := []string{"a_variable_label"}
	labels := prometheus.Labels{
		"label_key": "label_value",
	}

	descs <- prometheus.NewDesc("meminfo", "a stupid help text", variableLabels, labels)
}

func (m *MeminfoCollector) Collect(metrics chan<- prometheus.Metric) {
	logger := log.New(os.Stdout, "memcollector", log.LstdFlags)
	collector, err := NewMeminfoCollector(*logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = collector.Update(metrics)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
	}()

}

func (m *MeminfoCollector) Update(ch chan<- prometheus.Metric) error {

	memInfo, err := GetMemInfo()
	if err != nil {
		panic(err)
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
				prometheus.BuildFQName("sysperf", "memory", k),
				fmt.Sprintf("Memory information field %s.", k),
				nil, nil,
			),
			metricType, v,
		)
	}

	return nil
}

func NewMeminfoCollector(logger log.Logger) (Collector, error) {
	return &MeminfoCollector{logger}, nil
}

func GetMemInfo() (map[string]float64, error) {

	file, err := os.Open(filepath.Join(procfs.DefaultMountPoint, "meminfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMemInfo(file)
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
