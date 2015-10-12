package main

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/metrics"
)

func main() {
	m := metrics.New("disk.usage")

	usage, err := diskUsage()
	if err == nil {
		m.Print(usage)
	}
}

func diskUsage() (float64, error) {
	out, err := exec.Command("df", "-lP").Output()
	if err != nil {
		return 0.0, nil
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")[1:]

	var used, totalUsed, available, totalAvailable float64
	for _, line := range lines {
		stats := strings.Fields(line)

		used, err = strconv.ParseFloat(stats[2], 64)
		if err != nil {
			return 0.0, err
		}
		totalUsed += used

		available, err = strconv.ParseFloat(stats[3], 64)
		if err != nil {
			return 0.0, err
		}
		totalAvailable += available
	}

	return 100.0 * totalUsed / (totalUsed + totalAvailable), nil
}
