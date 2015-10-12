package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hico-horiuchi/sensu-plugins-go/lib/handler"
)

type metricsStruct struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"@timestamp"`
}

func main() {
	h := handler.New("/etc/sensu/conf.d/handler-elasticsearch.json")
	lines := strings.Split(strings.TrimRight(h.Event.Check.Output, "\n"), "\n")

	for _, line := range lines {
		request, err := http.NewRequest("POST", url(&h.Event, &h.Config), strings.NewReader(payload(line)))
		if err != nil {
			continue
		}
		request.Header.Set("Content-Type", "application/json")

		http.DefaultClient.Do(request)
	}
}

func payload(line string) string {
	arr := strings.Fields(line)

	value, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		return ""
	}

	unix, err := strconv.ParseInt(arr[2], 10, 64)
	if err != nil {
		return ""
	}

	body, err := json.Marshal(metricsStruct{
		Key:       arr[0],
		Value:     value,
		Timestamp: time.Unix(unix, 0).Format(time.RFC3339),
	})
	if err != nil {
		return ""
	}

	return string(body)
}

func url(event *handler.EventStruct, config *handler.ConfigStruct) string {
	return fmt.Sprintf("http://%s:%d/%s-%s/%s/%x",
		config.GetPath("elasticsearch", "host").MustString(),
		config.GetPath("elasticsearch", "port").MustInt(),
		config.GetPath("elasticsearch", "index").MustString(),
		time.Now().Format("2006.01.02"),
		event.Check.Name,
		esID(),
	)
}

func esID() []byte {
	now := strconv.FormatInt(int64(time.Now().Nanosecond()), 10)
	hash := md5.New()
	hash.Write([]byte(now))

	return hash.Sum(nil)
}
