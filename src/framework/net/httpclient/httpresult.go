package httpclient

type HttpResult struct {
	StatusCode int
	IsSuccess  bool
	Content    string
	Message    string
}
