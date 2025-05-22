package geometry

import "testing"

func TestBrightnessToChar(t *testing.T) {
	cases := []struct {
		brightness float64
		expected   string
	}{
		{0.0, " "},
		{0.12, "."},
		{0.24, ":"},
		{0.29, "-"},
		{0.41, "="},
		{0.5, "+"},
		{0.6, "*"},
		{0.7, "#"},
		{0.8, "%"},
		{0.91, "$"},
		{1.0, "@"},
	}

	for _, c := range cases {
		actual := brightnessToChar(c.brightness)
		if actual != c.expected {
			t.Errorf("brightnessToChar(%f) expected %s, got %s", c.brightness, c.expected, actual)
		}
	}
}
