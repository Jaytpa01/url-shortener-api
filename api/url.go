package api

// CreateUrlRequest represents the expected request body
// when either shortening or lengthening a url
type CreateUrlRequest struct {
	Url string `json:"url"`
}

// UrlResponse is the expected response
// from the api when creating or reading a Url
type UrlResponse struct {
	Token     string `json:"token"`
	TargetUrl string `json:"target_url"`
	QRCode    string `json:"qr_code"`
}

type UrlVisitsResponse struct {
	Visits int `json:"visits"`
}
