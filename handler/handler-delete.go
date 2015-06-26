package main

import (
	"sort"

	"../lib/handler"
	"github.com/hico-horiuchi/ohgi/sensu"
)

func main() {
	h := handler.New("/etc/sensu/conf.d/handler-delete.json")

	status := h.Config.GetPath("delete", "status").MustInt()
	contain := contains(h.Event.Client.Subscriptions, h.Config.GetPath("delete", "subscriptions").MustArray())
	if h.Event.Check.Name != "keepalive" || h.Event.Check.Status != status || !contain {
		return
	}

	sensu.DefaultAPI = &sensu.API{
		Host:     h.Config.GetPath("delete", "host").MustString(),
		Port:     h.Config.GetPath("delete", "port").MustInt(),
		User:     h.Config.GetPath("delete", "user").MustString(),
		Password: h.Config.GetPath("delete", "password").MustString(),
	}
	sensu.DefaultAPI.DeleteClientsClient(h.Event.Client.Name)
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
