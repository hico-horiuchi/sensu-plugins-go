package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/metrics"
)

func main() {
	var sleep int

	m := metrics.New("")
	m.Option.IntVarP(&sleep, "sleep", "s", 1, "SLEEP")
	m.Init()

	traffics, err := trafficsBytes(1)
	if err == nil {
		m.Scheme("traffic.rx_bytes").Print(float64(traffics[0]))
		m.Scheme("traffic.tx_bytes").Print(float64(traffics[1]))
	}
}

func trafficsBytes(sleep int) ([]int64, error) {
	beforeTraffics, err := getTraffics()
	if err != nil {
		return []int64{}, err
	}

	time.Sleep(time.Duration(sleep) * time.Second)

	afterTraffics, err := getTraffics()
	if err != nil {
		return []int64{}, err
	}

	return []int64{
		afterTraffics[0] - beforeTraffics[0],
		afterTraffics[1] - beforeTraffics[1],
	}, nil
}

func getTraffics() ([]int64, error) {
	traffics := make([]int64, 2)
	ifpaths, err := filepath.Glob("/sys/class/net/*")
	if err != nil {
		return []int64{}, err
	}

	for _, ifpath := range ifpaths {
		info, err := os.Stat(ifpath)
		if err != nil {
			return []int64{}, err
		}
		base := path.Base(ifpath)

		if !info.IsDir() || base == "lo" {
			continue
		}

		rxBytes, err := ioutil.ReadFile(ifpath + "/statistics/rx_bytes")
		if err != nil {
			return []int64{}, err
		}
		traffics[0], err = strconv.ParseInt(strings.TrimRight(string(rxBytes), "\n"), 10, 64)
		if err != nil {
			return []int64{}, err
		}

		txBytes, err := ioutil.ReadFile(ifpath + "/statistics/tx_bytes")
		if err != nil {
			return []int64{}, err
		}
		traffics[1], err = strconv.ParseInt(strings.TrimRight(string(txBytes), "\n"), 10, 64)
		if err != nil {
			return []int64{}, err
		}
	}

	return traffics, nil
}
