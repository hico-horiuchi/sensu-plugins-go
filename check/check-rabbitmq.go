package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"../sensu-plugin/check"
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

	c := check.New("CheckRabbitMQ")
	c.Option.StringVarP(&host, "host", "h", "localhost", "HOST")
	c.Option.IntVarP(&port, "port", "P", 15672, "PORT")
	c.Option.StringVarP(&vhost, "vhost", "v", "%2F", "VHOST")
	c.Option.StringVarP(&user, "user", "u", "guest", "USER")
	c.Option.StringVarP(&password, "password", "p", "guest", "PASSWORD")
	c.Option.IntVarP(&timeout, "timeout", "t", 10, "TIMEOUT")
	c.Init()

	status, err := alivenessTest(host, port, vhost, user, password, timeout)
	if err != nil {
		c.Error(err)
	}

	switch status {
	case "ok":
		c.Ok("RabbitMQ server is alive")
	default:
		c.Warning("Object Not Found")
	}
}

func alivenessTest(host string, port int, vhost string, user string, password string, timeout int) (string, error) {
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
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	json.Unmarshal(body, &aliveness)
	return aliveness.Status, nil
}
