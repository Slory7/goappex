package security

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
)

func ComputeMD5(input string) string {
	md5 := md5.New()
	io.WriteString(md5, input)
	by := md5.Sum(nil)
	s := fmt.Sprintf("%x", by)
	return s
}

func ComputeSHA256(input string) string {
	sha_256 := sha256.New()
	sha_256.Write([]byte(input))
	by := sha_256.Sum(nil)
	s := fmt.Sprintf("%x", by)
	return s
}

func ComputeSHA512_256(input string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(input))
	by := sha512.Sum512_256([]byte(input))
	s := fmt.Sprintf("%x", by)
	return s
}
