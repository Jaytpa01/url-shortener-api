package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// mockJSONDecoder is a mock implementation of a utils.JSONDecoder
type mockJSONDecoder struct {
	mock.Mock
}

// NewMockJSONDecoder returns our mock implementation of utils.JSONDecoder
func NewMockJSONDecoder() *mockJSONDecoder {
	return &mockJSONDecoder{}
}

// DecodeJSON is a mock implementation of JSONDecoder.DecodeJSON
func (m *mockJSONDecoder) DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	ret := m.Called(w, r, dst)
	return ret.Error(0)
}
