package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

type alivenessStruct struct {
	Status string
}

func main() {
	var (
		host     string
		port     int
		vhost    string
		user     string
		password string
		timeout  int
	)

	pflag.StringVarP(&host, "host", "h", "localhost", "HOST")
	pflag.IntVarP(&port, "port", "p", 15672, "PORT")
	pflag.StringVarP(&vhost, "vhost", "v", "%2F", "VHOST")
	pflag.StringVarP(&user, "user", "u", "guest", "USER")
	pflag.StringVarP(&password, "password", "w", "guest", "PASSWORD")
	pflag.IntVarP(&timeout, "timeout", "t", 10, "TIMEOUT")
	pflag.Parse()

	status := alivenessTest(host, port, vhost, user, password, timeout)

	switch status {
	case "ok":
		fmt.Println("CheckRabbitMQ OK: RabbitMQ server is alive")
		os.Exit(0)
	case "":
		fmt.Println("CheckRabbitMQ WARNING: Object Not Found")
		os.Exit(1)
	}
}

func alivenessTest(host string, port int, vhost string, user string, password string, timeout int) string {
	var aliveness alivenessStruct
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second

	request := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Host:   host + ":" + strconv.Itoa(port),
			Scheme: "http",
			Opaque: "/api/aliveness-test/" + vhost,
		},
		Header: http.Header{
			"User-Agent": {"godoc-example/0.1"},
		},
	}
	request.SetBasicAuth(user, password)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Println("CheckRabbitMQ CRITICAL:", err)
		os.Exit(2)
	}
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(body, &aliveness)
	return aliveness.Status
}
