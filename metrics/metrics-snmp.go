package main

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/metrics"
)

func main() {
	var (
		hosts     string
		community string
		sleep     int
		wg        sync.WaitGroup
	)

	m := metrics.New("")
	m.Option.StringVarP(&hosts, "hosts", "h", "127.0.0.1", "HOSTS")
	m.Option.StringVarP(&community, "community", "c", "public", "COMMUNITY")
	m.Option.IntVarP(&sleep, "sleep", "s", 1, "SLEEP")
	m.Init()

	for _, host := range strings.Split(hosts, ",") {
		wg.Add(1)

		go func(host string) {
			var port string

			beforeTraffics, err := snmpWalk(host, community)
			if err != nil {
				wg.Done()
				return
			}

			time.Sleep(time.Duration(sleep) * time.Second)

			afterTraffics, err := snmpWalk(host, community)
			if err != nil {
				wg.Done()
				return
			}

			tmp := metrics.New("").Hostname(host)
			for i := range beforeTraffics[0] {
				port = strconv.Itoa(i + 1)
				tmp.Scheme("snmp.rx_bytes." + port).Print(float64(afterTraffics[0][i] - beforeTraffics[0][i]))
				tmp.Scheme("snmp.tx_bytes." + port).Print(float64(afterTraffics[1][i] - beforeTraffics[1][i]))
			}

			wg.Done()
		}(host)
	}

	wg.Wait()
}

func snmpWalk(host string, community string) ([2][]int64, error) {
	var traffic int64
	var traffics [2][]int64

	for i, oid := range []string{"1.3.6.1.2.1.2.2.1.10", "1.3.6.1.2.1.2.2.1.16"} {
		out, err := exec.Command("snmpwalk", "-v", "2c", "-c", community, host, oid).Output()
		if err != nil {
			return traffics, err
		}

		lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
		for _, line := range lines {
			traffic, err = strconv.ParseInt(strings.Fields(line)[3], 10, 64)
			if err != nil {
				return traffics, err
			}

			traffics[i] = append(traffics[i], traffic)
		}
	}

	return traffics, nil
}
