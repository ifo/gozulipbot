package gozulipbot

import (
	"testing"
)

func TestMakeBot(t *testing.T) {
	cases := map[string][]string{
		"1": {"email@example.com", "secretkey", "one-stream"},
	}

	for k, c := range cases {
		b := MakeBot(c[0], c[1], []string{c[2]})
		if b.EmailAddress != c[0] {
			t.Errorf("case %s, actual %s, expected %s", k, b.EmailAddress, c[0])
		}
		if b.ApiKey != c[1] {
			t.Errorf("case %s, actual %s, expected %s", k, b.ApiKey, c[1])
		}
		if b.Streams[0] != c[2] {
			t.Errorf("case %s, actual %s, expected %s", k, b.Streams[0], c[2])
		}
	}
}
