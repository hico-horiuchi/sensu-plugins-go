package main

import (
	"net/http"
	"strconv"
	"time"

	"../lib/check"
)

func main() {
	var (
		url      string
		redirect bool
		timeout  int
	)

	c := check.New("CheckHTTP")
	c.Option.StringVarP(&url, "url", "u", "http://localhost/", "URL")
	c.Option.BoolVarP(&redirect, "redirect", "r", false, "REDIRECT")
	c.Option.IntVarP(&timeout, "timeout", "t", 15, "TIMEOUT")
	c.Init()

	status, err := statusCode(url, timeout)
	if err != nil {
		c.Error(err)
	}

	switch {
	case status >= 400:
		c.Critical(strconv.Itoa(status))
	case status >= 300 && redirect:
		c.Ok(strconv.Itoa(status))
	case status >= 300:
		c.Warning(strconv.Itoa(status))
	default:
		c.Ok(strconv.Itoa(status))
	}
}

func statusCode(url string, timeout int) (int, error) {
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}
