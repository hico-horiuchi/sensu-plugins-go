package main

import (
	"net"
	"strconv"
	"time"

	"../lib/check"
)

func main() {
	var (
		host    string
		port    int
		timeout int64
	)

	c := check.New("CheckPing")
	c.Option.StringVarP(&host, "host", "h", "localhost", "HOST")
	c.Option.IntVarP(&port, "port", "P", 22, "PORT")
	c.Option.Int64VarP(&timeout, "timeout", "t", 5, "TIMEOUT")
	c.Init()

	address := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)

	if err != nil {
		c.Error(err)
	}
	defer conn.Close()

	c.Ok(address)
}
