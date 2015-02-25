package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var (
		warn    int
		crit    int
		warnMnt []string
		critMnt []string
	)

	pflag.IntVarP(&warn, "warn", "w", 80, "WARN")
	pflag.IntVarP(&crit, "crit", "c", 100, "CRIT")
	pflag.Parse()

	usage := diskUsage()

	for i := range usage {
		use, _ := strconv.ParseInt(strings.TrimRight(usage[i][1], "%"), 10, 0)
		switch {
		case use > int64(crit):
			critMnt = append(critMnt, usage[i][0]+" "+usage[i][1])
		case use > int64(warn):
			warnMnt = append(warnMnt, usage[i][0]+" "+usage[i][1])
		}
	}

	switch {
	case len(critMnt) > 0:
		fmt.Printf("CheckDisk CRITICAL: %s\n", strings.Join(critMnt, ", "))
		os.Exit(2)
	case len(warnMnt) > 0:
		fmt.Printf("CheckDisk WARNING: %s\n", strings.Join(warnMnt, ", "))
		os.Exit(1)
	default:
		fmt.Printf("CheckDisk OK\n")
		os.Exit(0)
	}
}

func diskUsage() [][]string {
	out, _ := exec.Command("df", "-l").Output()
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")[1:]
	result := make([][]string, len(lines))

	for i := range lines {
		stats := strings.Fields(lines[i])
		result[i] = []string{stats[5], stats[4]}
	}

	return result
}
