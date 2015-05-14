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
	var scheme string
	hostname, _ := os.Hostname()

	pflag.StringVarP(&scheme, "scheme", "s", hostname, "SCHEME")
	pflag.Parse()

	fmt.Printf("%s.cpu.usage %f %d\n", scheme, cpuUsage(1), time.Now().Unix())
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
	contents, _ := ioutil.ReadFile("/proc/stat")
	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Fields(line)[1:]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], _ = strconv.ParseFloat(stats[i], 64)
	}

	return result
}
