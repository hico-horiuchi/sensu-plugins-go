package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../lib/handler"
)

type metricsStruct struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"@timestamp"`
}

func main() {
	h := handler.New("/etc/sensu/conf.d/h-elasticsearch.json")
	lines := strings.Split(strings.TrimRight(h.Event.Check.Output, "\n"), "\n")

	for _, line := range lines {
		metrics := newMetrics(line)
		body, _ := json.Marshal(metrics)
		payload := strings.NewReader(string(body))

		url := createURL(h.Event, h.Config)
		request, _ := http.NewRequest("POST", url, payload)
		request.Header.Set("Content-Type", "application/json")

		http.DefaultClient.Do(request)
	}
}

func newMetrics(line string) metricsStruct {
	metrics := metricsStruct{}

	arr := strings.Fields(line)
	value, _ := strconv.ParseFloat(arr[1], 64)
	unix, _ := strconv.ParseInt(arr[2], 10, 64)

	metrics.Key = arr[0]
	metrics.Value = value
	metrics.Timestamp = time.Unix(unix, 0).Format(time.RFC3339)

	return metrics
}

func createURL(event handler.EventStruct, config handler.ConfigStruct) string {
	host := config.GetPath("elasticsearch", "host").MustString()
	port := config.GetPath("elasticsearch", "port").MustInt()
	now := time.Now().Format("2006.01.02")

	esIndex := config.GetPath("elasticsearch", "index").MustString()
	esType := event.Check.Name
	esId := createID()

	return fmt.Sprintf("http://%s:%d/%s-%s/%s/%x", host, port, esIndex, now, esType, esId)
}

func createID() []byte {
	now := strconv.FormatInt(int64(time.Now().Nanosecond()), 10)

	hash := md5.New()
	hash.Write([]byte(now))

	return hash.Sum(nil)
}
