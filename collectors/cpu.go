package collectors

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type CPUCollector struct{}

var cpuInfoDesc = prometheus.NewDesc(
	prometheus.BuildFQName("sysperf", "cpu", "info"),
	"CPU information from /proc/cpuinfo.",
	[]string{"package", "core", "cpu", "vendor", "family", "model", "model_name", "microcode", "stepping", "cachesize"}, nil,
)

var cpuPackageThrottleDesc = prometheus.NewDesc(
	prometheus.BuildFQName("sysperf", "cpu", "package_throttles_total"),
	"Number of times this CPU package has been throttled.",
	[]string{"package"}, nil,
)

var cpuCoreThrottleDesc = prometheus.NewDesc(
	prometheus.BuildFQName("sysperf", "cpu", "core_throttles_total"),
	"Number of times this CPU core has been throttled.",
	[]string{"package", "core"}, nil,
)

func (c CPUCollector) Collect(ch chan<- prometheus.Metric) {
	fs, err := procfs.NewDefaultFS()

	if err != nil {
		log.Error().Err(err)
		return
	}
	info, err := fs.CPUInfo()
	if err != nil {
		log.Error().Err(err)
		return
	}

	for _, cpu := range info {
		ch <- prometheus.MustNewConstMetric(cpuInfoDesc,
			prometheus.GaugeValue,
			1,
			cpu.PhysicalID,
			cpu.CoreID,
			strconv.Itoa(int(cpu.Processor)),
			cpu.VendorID,
			cpu.CPUFamily,
			cpu.Model,
			cpu.ModelName,
			cpu.Microcode,
			cpu.Stepping,
			cpu.CacheSize)
	}

	cpus, err := filepath.Glob("/sys/devices/system/cpu/cpu[0-9]*")
	if err != nil {
		log.Error().Err(err)
		return
	}

	packageThrottles := make(map[uint64]uint64)
	packageCoreThrottles := make(map[uint64]map[uint64]uint64)

	for _, cpu := range cpus {

		var physicalPackageID, coreID, coreThrottleCount, packageThrottleCount uint64

		if physicalPackageID, err = readUintFromFile(filepath.Join(cpu, "topology", "physical_package_id")); err != nil {
			log.Debug().Msg(fmt.Sprintf("CPU is missing physical_package_id %s", cpu))
			continue
		}
		if coreID, err = readUintFromFile(filepath.Join(cpu, "topology", "core_id")); err != nil {
			log.Debug().Msg(fmt.Sprintf("CPU is missing core_id %s", cpu))
			continue
		}

		if _, present := packageCoreThrottles[physicalPackageID]; !present {
			packageCoreThrottles[physicalPackageID] = make(map[uint64]uint64)
		}

		if _, present := packageCoreThrottles[physicalPackageID][coreID]; !present {
			if coreThrottleCount, err = readUintFromFile(filepath.Join(cpu, "thermal_throttle", "core_throttle_count")); err == nil {
				packageCoreThrottles[physicalPackageID][coreID] = coreThrottleCount
			} else {
				log.Debug().Msg(fmt.Sprintf("CPU is missing core_throttle_count %s", cpu))
			}
		}

		if _, present := packageThrottles[physicalPackageID]; !present {
			if packageThrottleCount, err = readUintFromFile(filepath.Join(cpu, "thermal_throttle", "package_throttle_count")); err == nil {
				packageThrottles[physicalPackageID] = packageThrottleCount
			} else {
				log.Debug().Msg(fmt.Sprintf("CPU is missing package_throttle_count %s", cpu))
			}
		}
	}

	for physicalPackageID, packageThrottleCount := range packageThrottles {
		ch <- prometheus.MustNewConstMetric(cpuPackageThrottleDesc,
			prometheus.CounterValue,
			float64(packageThrottleCount),
			strconv.FormatUint(physicalPackageID, 10))
	}

	for physicalPackageID, coreMap := range packageCoreThrottles {
		for coreID, coreThrottleCount := range coreMap {
			ch <- prometheus.MustNewConstMetric(cpuCoreThrottleDesc,
				prometheus.CounterValue,
				float64(coreThrottleCount),
				strconv.FormatUint(physicalPackageID, 10),
				strconv.FormatUint(coreID, 10))
		}
	}

}

func (c CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cpuInfoDesc
	ch <- cpuCoreThrottleDesc
	ch <- cpuPackageThrottleDesc
}

func readUintFromFile(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
