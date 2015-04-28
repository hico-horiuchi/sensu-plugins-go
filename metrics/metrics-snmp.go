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

		go func(host string) {
			beforeTraffics := snmpWalk(host, community)
			time.Sleep(time.Second)
			afterTraffics := snmpWalk(host, community)

			fmt.Printf("%s.snmp.rx_bytes %d %d\n", host, afterTraffics[0]-beforeTraffics[0], now)
			fmt.Printf("%s.snmp.tx_bytes %d %d\n", host, afterTraffics[1]-beforeTraffics[1], now)
			wg.Done()
		}(host)
	}

	wg.Wait()
}

func snmpWalk(host string, community string) []int64 {
	traffics := make([]int64, 2)

	for i, oid := range []string{"1.3.6.1.2.1.2.2.1.10", "1.3.6.1.2.1.2.2.1.16"} {
		out, _ := exec.Command("snmpwalk", "-v", "2c", "-c", community, host, oid).Output()
		lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")

		for _, line := range lines {
			traffics[i] += ParseInt(strings.Fields(line)[3])
		}
	}

	return traffics
}

func ParseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
