package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"../lib/check"
)

func main() {
	var (
		warn  int
		crit  int
		sleep int
	)

	c := check.New("CheckCPU")
	c.Option.IntVarP(&warn, "warn", "w", 80, "WARN")
	c.Option.IntVarP(&crit, "crit", "c", 90, "CRIT")
	c.Option.IntVarP(&sleep, "sleep", "s", 1, "SLEEP")
	c.Init()

	usage, err := cpuUsage(sleep)
	if err != nil {
		c.Error(err)
	}

	switch {
	case usage >= float64(crit):
		c.Critical(fmt.Sprintf("%.0f%%", usage))
	case usage >= float64(warn):
		c.Warning(fmt.Sprintf("%.0f%%", usage))
	default:
		c.Ok(fmt.Sprintf("%.0f%%", usage))
	}
}

func cpuUsage(sleep int) (float64, error) {
	var usage, totalDiff float64

	beforeStats, err := getStats()
	if err != nil {
		return usage, err
	}

	time.Sleep(time.Duration(sleep) * time.Second)

	afterStats, err := getStats()
	if err != nil {
		return usage, err
	}

	diffStats := make([]float64, len(beforeStats))
	for i := range beforeStats {
		diffStats[i] = afterStats[i] - beforeStats[i]
		totalDiff += diffStats[i]
	}

	usage = 100.0 * (totalDiff - diffStats[3]) / totalDiff
	return usage, nil
}

func getStats() ([]float64, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return []float64{}, err
	}

	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Fields(line)[1:]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], err = strconv.ParseFloat(stats[i], 64)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
