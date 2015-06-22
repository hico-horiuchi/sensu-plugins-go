package main

import (
	"os/exec"
	"strconv"
	"strings"

	"../lib/metrics"
)

func main() {
	m := metrics.New("memory.usage")

	usage, err := memoryUsage()
	if err == nil {
		m.Print(usage)
	}
}

func memoryUsage() (float64, error) {
	out, err := exec.Command("free").Output()
	if err != nil {
		return 0.0, err
	}
	lines := strings.Split(string(out), "\n")

	total, err := strconv.ParseFloat(strings.Fields(lines[1])[1], 64)
	if err != nil {
		return 0.0, err
	}

	free, err := strconv.ParseFloat(strings.Fields(lines[2])[3], 64)
	if err != nil {
		return 0.0, err
	}

	return 100.0 - (100.0 * free / total), nil
}
