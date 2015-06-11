package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/pflag"
)

func main() {
	var (
		host  string
		port  int
		key   string
		value string
	)

	pflag.StringVarP(&host, "host", "h", "localhost", "HOST")
	pflag.IntVarP(&port, "port", "p", 6379, "PORT")
	pflag.StringVarP(&key, "key", "k", "role", "KEY")
	pflag.StringVarP(&value, "value", "v", "master", "VALUE")
	pflag.Parse()

	info := getInfo(host, port, key)
	if info == value {
		fmt.Printf("CheckRedis OK: Redis %s is %s\n", key, info)
		os.Exit(0)
	} else {
		fmt.Printf("CheckRedis CRITICAL: Redis %s is %s\n", key, info)
		os.Exit(2)
	}
}

func getInfo(host string, port int, key string) string {
	var info string

	client, err := redis.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("CheckRedis CRITICAL:", err)
		os.Exit(2)
	}
	defer client.Close()

	result, err := redis.String(client.Do("INFO"))
	if err != nil {
		fmt.Println("CheckRedis CRITICAL:", err)
		os.Exit(2)
	}

	re := regexp.MustCompile(key + ":(.+)")
	match := re.FindAllStringSubmatch(result, -1)
	if len(match) > 0 {
		info = strings.TrimRight(match[0][1], "\r")
	}

	return info
}
