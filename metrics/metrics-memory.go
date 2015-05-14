package main

import (
	"fmt"
	"os"
	"os/exec"
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

	fmt.Printf("%s.memory.usage %f %d\n", scheme, memoryUsage(), time.Now().Unix())
}

func memoryUsage() float64 {
	out, _ := exec.Command("free").Output()
	lines := strings.Split(string(out), "\n")

	total, _ := strconv.ParseFloat(strings.Fields(lines[1])[1], 64)
	free, _ := strconv.ParseFloat(strings.Fields(lines[2])[3], 64)

	return 100.0 - (100.0 * free / total)
}
