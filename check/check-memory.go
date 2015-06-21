package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"../lib/check"
)

func main() {
	var (
		warn int
		crit int
	)

	c := check.New("CheckMemory")
	c.Option.IntVarP(&warn, "warn", "w", 80, "WARN")
	c.Option.IntVarP(&crit, "crit", "c", 90, "CRIT")
	c.Init()

	usage, err := memoryUsage()
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
