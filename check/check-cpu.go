package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var (
		warn  int
		crit  int
		sleep int
	)

	pflag.IntVarP(&warn, "warn", "w", 80, "WARN")
	pflag.IntVarP(&crit, "crit", "c", 90, "CRIT")
	pflag.IntVarP(&sleep, "sleep", "s", 1, "SLEEP")
	pflag.Parse()

	usage := cpuUsage(sleep)

	switch {
	case usage >= float64(crit):
		fmt.Printf("CheckCPU CRITICAL: %.0f%%\n", usage)
		os.Exit(2)
	case usage >= float64(warn):
		fmt.Printf("CheckCPU WARNING: %.0f%%\n", usage)
		os.Exit(1)
	default:
		fmt.Printf("CheckCPU OK: %.0f%%\n", usage)
		os.Exit(0)
	}
}

func cpuUsage(sleep int) float64 {
	beforeStats := getStats()
	time.Sleep(time.Duration(sleep) * time.Second)
	afterStats := getStats()

	diffStats := make([]float64, len(beforeStats))
	var totalDiff float64 = 0.0

	for i := range beforeStats {
		diffStats[i] = afterStats[i] - beforeStats[i]
		totalDiff += diffStats[i]
	}

	return 100.0 * (totalDiff - diffStats[3]) / totalDiff
}

func getStats() []float64 {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Println("CheckCPU CRITICAL:", err)
		os.Exit(2)
	}

	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Fields(line)[1:]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], _ = strconv.ParseFloat(stats[i], 64)
	}

	return result
}
