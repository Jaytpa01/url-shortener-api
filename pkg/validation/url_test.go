package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{
			"dfsdfsd",
			false,
		},
		{
			"google.com",
			false,
		},
		{
			"https://google.com",
			true,
		},
		{
			"fgj://sdf.co",
			false,
		},
		{
			"https://",
			false,
		},
		{
			"https://example.com",
			true,
		},
	}

	for _, test := range tests {
		assert.Equalf(t, test.valid, IsValidUrl(test.url), "expect url (%s) to be %t", test.url, test.valid)
	}
}
