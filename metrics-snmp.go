package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var (
		hosts     string
		community string
		wg        sync.WaitGroup
	)

	pflag.StringVarP(&hosts, "hosts", "h", "127.0.0.1", "HOSTS")
	pflag.StringVarP(&community, "community", "c", "public", "COMMUNITY")
	pflag.Parse()

	now := time.Now().Unix()

	for _, host := range strings.Split(hosts, ",") {
		wg.Add(1)

		go func() {
			before_traffic := snmpWalk(host, community)
			time.Sleep(time.Second)
			after_traffic := snmpWalk(host, community)

			fmt.Printf("%s.snmp.traffic %d %d\n", host, after_traffic-before_traffic, now)
			wg.Done()
		}()
	}

	wg.Wait()
}

func snmpWalk(host string, community string) int64 {
	out, _ := exec.Command("snmpwalk", "-v", "2c", "-c", community, host, "1.3.6.1.2.1.2.2.1.10").Output()
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")

	var traffic int64 = 0

	for _, line := range lines {
		traffic += ParseInt(strings.Fields(line)[3])
	}

	return traffic
}

func ParseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
