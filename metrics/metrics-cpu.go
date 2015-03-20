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
	var scheme string
	hostname, _ := os.Hostname()

	pflag.StringVarP(&scheme, "scheme", "s", hostname, "SCHEME")
	pflag.Parse()

	fmt.Printf("%s.cpu.usage %f %d\n", scheme, cpuUsage(1), time.Now().Unix())
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

	return 100.0 * (total_diff - diff_stats[3]) / total_diff
}

func getStats() []float64 {
	contents, _ := ioutil.ReadFile("/proc/stat")
	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Fields(line)[1:11]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], _ = strconv.ParseFloat(stats[i], 64)
	}

	return result
}
