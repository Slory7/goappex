package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nuveo/log"
)

func HttpGetShort(baseUrl string, relativeUrl string, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, "", "GET", "", TokenEmpty, "", false, timeout)
}

func HttpGet(baseUrl string, relativeUrl string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, "", "GET", "", token, xRequestID, isSecure, timeout)
}

//HttpPostShort Json
func HttpPostShort(baseUrl string, relativeUrl string, postJsonData string, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, APPLICATION_JSON, "POST", postJsonData, TokenEmpty, "", false, timeout)
}

//HttpPost Json
func HttpPost(baseUrl string, relativeUrl string, postJsonData string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, APPLICATION_JSON, "POST", postJsonData, token, xRequestID, isSecure, timeout)
}

func HttpPostForm(baseUrl string, relativeUrl string, postData map[string]string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, FORM, "POST", BuildRequestParameters(postData), token, xRequestID, isSecure, timeout)
}

func HttpDelete(baseUrl string, relativeUrl string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, "", "DELETE", "", token, xRequestID, isSecure, timeout)
}

func HttpPut(baseUrl string, relativeUrl string, postJsonData string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	return HttpSend(baseUrl, relativeUrl, nil, nil, APPLICATION_JSON, "PUT", postJsonData, token, xRequestID, isSecure, timeout)
}

func HttpSend(baseUrl string, relativeUrl string, urlParams map[string]string, headers map[string]string, contentType string, method string, postData string, token TokenInfo, xRequestID string, isSecure bool, timeout int) (result HttpResult, err error) {
	startTime := time.Now()
	url := baseUrl + relativeUrl
	if urlParams != nil {
		url += "?" + BuildRequestParameters(urlParams)
	}
	logs := fmt.Sprintf("HttpClient:[%s] URL:%s\n", method, url)

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	var reqReader io.Reader
	if len(postData) > 0 {
		var sData string
		if isSecure {
			sData = "[SECURITY]"
		} else {
			sData = postData
		}
		logs += fmt.Sprintf("PostData:%s\n", sData)
		reqReader = strings.NewReader(postData)
	}
	req, err := http.NewRequest(method, url, reqReader)
	if err == nil {
		req.Header.Set("accept", "application/json")
		req.Header.Set("User-Agent", "GoAppEx/1.0")

		if len(contentType) > 0 {
			logs += fmt.Sprintf("ContentType:%s\n", contentType)
			req.Header.Set("Content-Type", contentType)
		}
		if token.IsValid() {
			logs += fmt.Sprintf("Token:%s\n", token.String()[:10]+"...")
			req.Header.Set("authorization", token.String())
		}
		if len(xRequestID) > 0 {
			logs += fmt.Sprintf("X-Request-ID:%s\n", xRequestID)
			req.Header.Set("X-Request-ID", xRequestID)
		}
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()

			result.StatusCode = resp.StatusCode
			result.IsSuccess = isStatusSuccess(result.StatusCode)

			content, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				result.Content = string(content)
				if !result.IsSuccess && len(result.Content) > 0 {
					var omsg msgObj
					if er := json.Unmarshal(content, &omsg); er == nil {
						if len(omsg.Message) > 0 {
							result.Message = omsg.Message
						} else if len(omsg.Msg) > 0 {
							result.Message = omsg.Msg
						} else if len(omsg.Error) > 0 {
							result.Message = omsg.Error
						}
					}
				}
			}
		}
	}
	logs += fmt.Sprintf("Status:%d\n", result.StatusCode)
	if err != nil {
		logs += fmt.Sprintf("Error:%v\n", err)
	}
	if len(result.Message) > 0 {
		logs += fmt.Sprintf("Message:%s\n", result.Message)
	}
	if len(result.Content) > 0 {
		if isSecure {
			logs += fmt.Sprintf("Content:%s\n", "[SECURITY]")
		} else {
			logs += fmt.Sprintf("Content:%s\n", result.Content)
		}
	}
	logs += fmt.Sprintf("Elapsed Time:%v\n", time.Since(startTime))
	log.Println(logs)
	return
}

func isStatusSuccess(status int) bool {
	return status >= 200 && status <= 299
}

type msgObj struct {
	Msg     string
	Message string
	Error   string
}
