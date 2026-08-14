package handlers

import "testing"

func TestParseAmount(t *testing.T) {
	cases := []struct {
		in      string
		want    float64
		wantErr bool
	}{
		{"12,50", 12.5, false},
		{"87.40", 87.4, false},
		{"1 142,86", 1142.86, false},
		{"1 142,86", 1142.86, false}, // non-breaking-space thousands separator
		{"-12,50", -12.5, false},     // returns are recorded as negative amounts
		{" -12,50 ", -12.5, false},
		{"-1 000", -1000, false},
		{"0", 0, true},
		{"0,00", 0, true},
		{"", 0, true},
		{"abc", 0, true},
		{"12,50 zł", 0, true},
		{"NaN", 0, true},
		{"Inf", 0, true},
		{"-Inf", 0, true},
	}
	for _, c := range cases {
		got, err := parseAmount(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("parseAmount(%q) = %v, want error", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseAmount(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("parseAmount(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
