package main

import (
	"./sensu"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type metricsStruct struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"@timestamp"`
}

func main() {
	handler := sensu.NewHandler("/etc/sensu/conf.d/handler-elasticsearch.json")
	lines := strings.Split(strings.TrimRight(handler.Event.Check.Output, "\n"), "\n")

	for i := range lines {
		metrics := newMetrics(lines[i])
		body, _ := json.Marshal(metrics)
		payload := strings.NewReader(string(body))

		url := createURL(handler.Event, handler.Config)
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

func createURL(event sensu.EventStruct, config simplejson.Json) string {
	host := config.GetPath("elasticsearch", "host").MustString()
	port := config.GetPath("elasticsearch", "port").MustInt()
	now := time.Now().Format("2006.01.02")

	es_index := config.GetPath("elasticsearch", "index").MustString()
	es_type := event.Check.Name
	es_id := createID()

	return fmt.Sprintf("http://%s:%d/%s-%s/%s/%x", host, port, es_index, now, es_type, es_id)
}

func createID() []byte {
	now := strconv.FormatInt(int64(time.Now().Nanosecond()), 10)

	hash := md5.New()
	hash.Write([]byte(now))

	return hash.Sum(nil)
}
