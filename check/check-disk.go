package main

import (
	"os/exec"
	"strconv"
	"strings"

	"../sensu-plugin/check"
)

func main() {
	var (
		warn    int
		crit    int
		warnMnt []string
		critMnt []string
	)

	c := check.New("CheckDisk")
	c.Option.IntVarP(&warn, "warn", "w", 80, "WARN")
	c.Option.IntVarP(&crit, "crit", "c", 100, "CRIT")
	c.Init()

	usage, err := diskUsage()
	if err != nil {
		c.Error(err)
	}

	for _, u := range usage {
		cap, err := strconv.ParseInt(strings.TrimRight(u[1], "%"), 10, 64)
		if err != nil {
			c.Error(err)
		}

		switch {
		case cap >= int64(crit):
			critMnt = append(critMnt, u[0]+" "+u[1])
		case cap >= int64(warn):
			warnMnt = append(warnMnt, u[0]+" "+u[1])
		}
	}

	switch {
	case len(critMnt) > 0:
		c.Critical(strings.Join(critMnt, ", "))
	case len(warnMnt) > 0:
		c.Warning(strings.Join(warnMnt, ", "))
	default:
		c.Ok("OK")
	}
}

func diskUsage() ([][]string, error) {
	out, err := exec.Command("df", "-lP").Output()
	if err != nil {
		return [][]string{}, err
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")[1:]
	result := make([][]string, len(lines))

	for i := range lines {
		stats := strings.Fields(lines[i])
		result[i] = []string{stats[5], stats[4]}
	}

	return result, nil
}
