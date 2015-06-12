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
		path string
		warn int
		crit int
	)

	pflag.StringVarP(&path, "path", "p", "/usr/bin/mailq", "PATH")
	pflag.IntVarP(&warn, "warn", "w", 5, "WARN")
	pflag.IntVarP(&crit, "crit", "c", 10, "CRIT")
	pflag.Parse()

	queue := mailQueue(path)

	switch {
	case queue > crit:
		fmt.Printf("CheckPostfix CRITICAL: %d messages in the postfix mail queue\n", queue)
		os.Exit(2)
	case queue > warn:
		fmt.Printf("CheckPostfix WARNING: %d messages in the postfix mail queue\n", queue)
		os.Exit(1)
	default:
		fmt.Printf("CheckPostfix OK: %d messages in the postfix mail queue\n", queue)
		os.Exit(0)
	}
}

func mailQueue(path string) int {
	var queue int

	cmd := path + " | egrep '[0-9]+ Kbytes in [0-9]+ Request|Mail queue is empty'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	result := string(out)
	if result == "Mail queue is empty\n" {
		queue = 0
	} else {
		queue, _ = strconv.Atoi(strings.Split(result, " ")[4])
	}

	return queue
}
