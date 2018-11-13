package httpclient

import (
	"net/url"
)

func BuildRequestParameters(m map[string]string) string {
	s := ""
	for k, v := range m {
		if len(s) > 0 {
			s += "&"
		}
		s += k + "=" + url.QueryEscape(v)
	}
	return s
}
