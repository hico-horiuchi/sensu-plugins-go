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

	for _, u := range usage {
		cap, _ := strconv.ParseInt(strings.TrimRight(u[1], "%"), 10, 64)
		switch {
		case cap > int64(crit):
			critMnt = append(critMnt, u[0]+" "+u[1])
		case cap > int64(warn):
			warnMnt = append(warnMnt, u[0]+" "+u[1])
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
	out, err := exec.Command("df", "-lP").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")[1:]
	result := make([][]string, len(lines))

	for i := range lines {
		stats := strings.Fields(lines[i])
		result[i] = []string{stats[5], stats[4]}
	}

	return result
}
