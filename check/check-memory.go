package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

func main() {
	var (
		warn int
		crit int
	)

	pflag.IntVarP(&warn, "warn", "w", 80, "WARN")
	pflag.IntVarP(&crit, "crit", "c", 90, "CRIT")
	pflag.Parse()

	usage := memoryUsage()

	switch {
	case usage >= float64(crit):
		fmt.Printf("CheckMemory CRITICAL: %.0f%%\n", usage)
		os.Exit(2)
	case usage >= float64(warn):
		fmt.Printf("CheckMemory WARNING: %.0f%%\n", usage)
		os.Exit(1)
	default:
		fmt.Printf("CheckMemory OK: %.0f%%\n", usage)
		os.Exit(0)
	}
}

func memoryUsage() float64 {
	out, err := exec.Command("free").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	lines := strings.Split(string(out), "\n")
	total, _ := strconv.ParseFloat(strings.Fields(lines[1])[1], 64)
	free, _ := strconv.ParseFloat(strings.Fields(lines[2])[3], 64)

	return 100.0 - (100.0 * free / total)
}
