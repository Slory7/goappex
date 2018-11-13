package test

import (
	"fmt"
	"framework/net/httpclient"
	"testing"

	_ "github.com/crgimenes/goconfig/json"
	"github.com/satori/go.uuid"
)

func TestHttpPost(t *testing.T) {
	reqID := uuid.Must(uuid.NewV4()).String()
	result, err := httpclient.HttpSend("http://api.xxx.cn", "/token", nil, nil, httpclient.FORM, "POST", "", httpclient.TokenEmpty, reqID, true, 5)
	fmt.Printf("%v", err)
	fmt.Printf("%v", result)
}
