package main

import (
	"net/http"
	"sort"
	"strconv"

	"./sensu/plugin"
)

func main() {
	handler := plugin.NewHandler("/etc/sensu/conf.d/handler-delete.json")
	client := handler.Event.Client
	check := handler.Event.Check
	config := handler.Config

	status := config.GetPath("delete", "status").MustInt()
	contain := contains(client.Subscriptions, config.GetPath("delete", "subscriptions").MustArray())
	if check.Name != "keepalive" || check.Status != status || !contain {
		return
	}

	host := config.GetPath("delete", "host").MustString()
	port := config.GetPath("delete", "port").MustInt()
	url := "http://" + host + ":" + strconv.Itoa(port) + "/clients/" + client.Name

	user := config.GetPath("delete", "user").MustString()
	password := config.GetPath("delete", "password").MustString()
	request, _ := http.NewRequest("DELETE", url, nil)
	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}

	http.DefaultClient.Do(request)
}

func contains(list []string, keys []interface{}) bool {
	sort.Strings(list)

	for _, v := range keys {
		key := v.(string)
		i := sort.SearchStrings(list, key)

		if i < len(list) && list[i] == key {
			return true
		}
	}

	return false
}
