package models

import "testing"

const (
	script = "<script></script>"
)

func TestRegexWhitelistFilter(t *testing.T) {
	cases := []struct {
		sin, rin     string
		wants, wante bool
	}{
		{"aA.bc_d", RegexFilterPackage, true, false},
		{"aA.bc_d" + script, RegexFilterPackage, false, true},
		{"Hello, 世界.'", RegexFilterText, true, false},
		{"Hello, 世界.'" + script, RegexFilterText, false, true},
	}
	for _, c := range cases {
		got, err := NewUserString(c.sin, 100, RegexBlacklist(c.rin))
		if (got != "") != c.wants {
			t.Errorf("Ouptput doesn't match: got %q want %q. case %q", got, c.wants, c)
		}
		if (err != nil) != c.wante {
			t.Errorf("Error output doesn't match: err %q want %q", err, c.wante)
		}
	}
}
