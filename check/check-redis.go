package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"../sensu-plugin/check"
	"github.com/garyburd/redigo/redis"
)

func main() {
	var (
		host  string
		port  int
		key   string
		value string
	)

	c := check.New("CheckRedis")
	c.Option.StringVarP(&host, "host", "h", "localhost", "HOST")
	c.Option.IntVarP(&port, "port", "P", 6379, "PORT")
	c.Option.StringVarP(&key, "key", "k", "role", "KEY")
	c.Option.StringVarP(&value, "value", "v", "master", "VALUE")
	c.Init()

	info, err := redisInfo(host, port, key)
	if err != nil {
		c.Error(err)
	}

	if info == value {
		c.Ok(fmt.Sprintf("Redis %s is %s", key, info))
	} else {
		c.Warning(fmt.Sprintf("Redis %s is %s", key, info))
	}
}

func redisInfo(host string, port int, key string) (string, error) {
	var info string

	client, err := redis.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return info, err
	}
	defer client.Close()

	result, err := redis.String(client.Do("INFO"))
	if err != nil {
		return info, err
	}

	re := regexp.MustCompile(key + ":(.+)")
	match := re.FindStringSubmatch(result)
	if len(match) > 0 {
		info = strings.TrimRight(match[1], "\r")
	}

	return info, nil
}
