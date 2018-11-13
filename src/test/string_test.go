package test

import (
	"fmt"
	"framework/utils"
	"testing"
	"time"
)

func TestReverse(t *testing.T) {
	for _, c := range []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	} {
		got := utils.Reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
func TestToString(t *testing.T) {
	m := time.Minute * 3
	fmt.Println(m.String())
}
