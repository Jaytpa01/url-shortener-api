package api

// CreateURLRequest represents the expected request body
// when either shortening or lengthening a url
type CreateURLRequest struct {
	URL string `json:"url"`
}

// URLResponse is the expected response
// from the api when creating or reading a URL
type URLResponse struct {
	Token          string `json:"token"`
	DestinationURL string `json:"url"`
	QRCode         string `json:"qr_code"`
	Visits         int    `json:"visits"`
}
