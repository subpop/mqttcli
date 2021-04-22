package mqttcli

import (
	"net/http"
	"strings"
)

type HeadersValue struct {
	Headers http.Header
}

func (v HeadersValue) String() string {
	var s string
	for key, vals := range v.Headers {
		for _, val := range vals {
			s += key + "=" + val + ","
		}
	}
	return strings.Trim(s, ",")
}

func (v HeadersValue) Set(s string) error {
	keyValuePairs := strings.Split(s, ",")
	for _, kvp := range keyValuePairs {
		p := strings.Split(kvp, "=")
		key := p[0]
		val := p[1]
		v.Headers.Add(key, val)
	}
	return nil
}
