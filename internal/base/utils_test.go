package base

import "testing"

func TestParsePrice(t *testing.T) {
	cases := map[string]int64{
		"8.70":  870,
		"4.10":  410,
		"13.60": 1360,
		"7":     700,
	}
	for in, want := range cases {
		got, err := ParsePrice(in)
		if err != nil {
			t.Fatalf("ParsePrice(%q) returned error: %v", in, err)
		}
		if got != want {
			t.Errorf("ParsePrice(%q) = %d, want %d", in, got, want)
		}
	}
}
