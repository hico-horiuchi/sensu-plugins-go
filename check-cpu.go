package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		warn  int
		crit  int
		sleep int
	)

	pflag.IntVarP(&warn, "warn", "w", 80, "WARN")
	pflag.IntVarP(&crit, "crit", "c", 100, "CRIT")
	pflag.IntVarP(&sleep, "sleep", "s", 1, "SLEEP")
	pflag.Parse()

	usage := cpuUsage(sleep)

	switch {
	case usage > float64(crit):
		fmt.Printf("CheckCPU CRITICAL: %f\n", usage)
		os.Exit(2)
	case usage > float64(warn):
		fmt.Printf("CheckCPU WARNING: %f\n", usage)
		os.Exit(1)
	default:
		fmt.Printf("CheckCPU OK: %f\n", usage)
		os.Exit(0)
	}
}

func cpuUsage(sleep int) float64 {
	before_stats := getStats()
	time.Sleep(time.Duration(sleep) * time.Second)
	after_stats := getStats()

	diff_stats := make([]float64, len(before_stats))
	var total_diff float64 = 0.0

	for i := range before_stats {
		diff_stats[i] = after_stats[i] - before_stats[i]
		total_diff += diff_stats[i]
	}

	return float64(100.0) * (total_diff - diff_stats[3]) / total_diff
}

func getStats() []float64 {
	contents, _ := ioutil.ReadFile("/proc/stat")
	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Split(line, " ")[2:12]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], _ = strconv.ParseFloat(stats[i], 64)
	}

	return result
}
