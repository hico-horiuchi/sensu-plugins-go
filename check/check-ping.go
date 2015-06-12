package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var (
		host    string
		port    int
		timeout int64
	)

	pflag.StringVarP(&host, "host", "h", "localhost", "HOST")
	pflag.IntVarP(&port, "port", "p", 80, "PORT")
	pflag.Int64VarP(&timeout, "timeout", "T", 1, "TIMEOUT")
	pflag.Parse()

	address := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)

	if err != nil {
		fmt.Println("CheckPing CRITICAL:", err)
		os.Exit(2)
	}
	defer conn.Close()

	fmt.Printf("CheckPing OK: %s\n", address)
}
