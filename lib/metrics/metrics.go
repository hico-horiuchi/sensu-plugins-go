package metrics

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type metricsStruct struct {
	hostname string
	scheme   string
	Option   *pflag.FlagSet
}

func New(scheme string) *metricsStruct {
	fqdn, _ := os.Hostname()

	metrics := &metricsStruct{
		hostname: strings.Split(fqdn, ".")[0],
		scheme:   scheme,
		Option:   pflag.NewFlagSet(scheme, 1),
	}

	return metrics
}

func (m *metricsStruct) Hostname(hostname string) *metricsStruct {
	m.hostname = hostname
	return m
}

func (m *metricsStruct) Scheme(scheme string) *metricsStruct {
	m.scheme = scheme
	return m
}

func (m metricsStruct) Init() {
	m.Option.Parse(os.Args[1:])
}

func (m metricsStruct) Print(value float64) {
	fmt.Printf("%s.%s %f %d\n", m.hostname, m.scheme, value, time.Now().Unix())
}
