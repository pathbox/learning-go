package parseform

import (
	"net/url"
	"strings"
)

func parseFormData(re url.Values) map[string]string {
	m := make(map[string]string)
	for k, v := range re {
		if strings.HasPrefix(k, "data[") {
			l := len(k)
			key := k[5 : l-1]
			m[key] = v[0]
		}
	}
	return m
}
