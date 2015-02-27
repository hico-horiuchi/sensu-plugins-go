package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	var scheme string
	hostname, _ := os.Hostname()

	pflag.StringVarP(&scheme, "scheme", "s", hostname, "SCHEME")
	pflag.Parse()

	fmt.Printf("%s.traffic.bytes %d %d\n", scheme, trafficBytes(1), time.Now().Unix())
}

func trafficBytes(sleep int) int64 {
	before_traffic := getTraffic()
	time.Sleep(time.Duration(sleep) * time.Second)
	after_traffic := getTraffic()

	return after_traffic - before_traffic
}

func getTraffic() int64 {
	var traffic int64 = 0
	ifpaths, _ := filepath.Glob("/sys/class/net/*")

	for _, ifpath := range ifpaths {
		info, _ := os.Stat(ifpath)
		base := path.Base(ifpath)

		if !info.IsDir() || base == "lo" {
			continue
		}

		txBytes, _ := ioutil.ReadFile(ifpath + "/statistics/tx_bytes")
		traffic += ParseInt(strings.TrimRight(string(txBytes), "\n"))

		rxBytes, _ := ioutil.ReadFile(ifpath + "/statistics/rx_bytes")
		traffic += ParseInt(strings.TrimRight(string(rxBytes), "\n"))
	}

	return traffic
}

func ParseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
