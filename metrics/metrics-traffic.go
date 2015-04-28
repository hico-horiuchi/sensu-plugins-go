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

	traffics := trafficsBytes(1)
	now := time.Now().Unix()

	fmt.Printf("%s.traffic.rx_bytes %d %d\n", scheme, traffics[0], now)
	fmt.Printf("%s.traffic.tx_bytes %d %d\n", scheme, traffics[1], now)
}

func trafficsBytes(sleep int) []int64 {
	beforeTraffics := getTraffics()
	time.Sleep(time.Duration(sleep) * time.Second)
	afterTraffics := getTraffics()

	return []int64{
		afterTraffics[0] - beforeTraffics[0],
		afterTraffics[1] - beforeTraffics[1],
	}
}

func getTraffics() []int64 {
	traffics := make([]int64, 2)
	ifpaths, _ := filepath.Glob("/sys/class/net/*")

	for _, ifpath := range ifpaths {
		info, _ := os.Stat(ifpath)
		base := path.Base(ifpath)

		if !info.IsDir() || base == "lo" {
			continue
		}

		rxBytes, _ := ioutil.ReadFile(ifpath + "/statistics/rx_bytes")
		traffics[0] = ParseInt(strings.TrimRight(string(rxBytes), "\n"))

		txBytes, _ := ioutil.ReadFile(ifpath + "/statistics/tx_bytes")
		traffics[1] = ParseInt(strings.TrimRight(string(txBytes), "\n"))
	}

	return traffics
}

func ParseInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
