package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateQRCodeLink(t *testing.T) {
	tests := []struct {
		url               string
		expectedQrCodeUrl string
	}{
		{
			"https://example.com",
			"https://api.qrserver.com/v1/create-qr-code/?data=https%3A%2F%2Fexample.com",
		},
		{
			"localhost:8080",
			"https://api.qrserver.com/v1/create-qr-code/?data=localhost%3A8080",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedQrCodeUrl, GenerateQRCodeLink(test.url))
	}
}
