package httpclient

type TokenInfo struct {
	Token_Type   string
	Access_Token string
}

var TokenEmpty TokenInfo = TokenInfo{}

func (t TokenInfo) String() string {
	return t.Token_Type + " " + t.Access_Token
}

func (t TokenInfo) IsValid() bool {
	return len(t.Access_Token) > 0
}
