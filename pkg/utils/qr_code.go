package utils

import "net/url"

func GenerateQRCodeLink(urlString string) string {
	escapedUrl := url.QueryEscape(urlString)
	return "https://api.qrserver.com/v1/create-qr-code/?data=" + escapedUrl
}
