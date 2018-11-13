package utils

import (
	"fmt"
)

//Reverse string
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

//JoinString parameters
func JoinString(sep string, params ...interface{}) string {
	var s string
	for index := 0; index < len(params); index++ {
		if index > 0 {
			s += sep
		}
		s += fmt.Sprintf("%v", params[index])
	}
	return s
}
