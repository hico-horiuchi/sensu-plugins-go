package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"../sensu-plugin/check"
)

func main() {
	var (
		path string
		warn int
		crit int
	)

	c := check.New("CheckPostfix")
	c.Option.StringVarP(&path, "path", "p", "/usr/bin/mailq", "PATH")
	c.Option.IntVarP(&warn, "warn", "w", 5, "WARN")
	c.Option.IntVarP(&crit, "crit", "c", 10, "CRIT")
	c.Init()

	queue, err := mailQueue(path)
	if err != nil {
		c.Error(err)
	}

	switch {
	case queue > crit:
		c.Critical(fmt.Sprintf("%d messages in the postfix mail queue\n", queue))
	case queue > warn:
		c.Warning(fmt.Sprintf("%d messages in the postfix mail queue\n", queue))
	default:
		c.Ok(fmt.Sprintf("%d messages in the postfix mail queue\n", queue))
	}
}

func mailQueue(path string) (int, error) {
	var queue int

	cmd := path + " | egrep '[0-9]+ Kbytes in [0-9]+ Request|Mail queue is empty'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return queue, err
	}

	result := string(out)
	if result == "Mail queue is empty\n" {
		queue = 0
	} else {
		queue, _ = strconv.Atoi(strings.Split(result, " ")[4])
	}

	return queue, nil
}
