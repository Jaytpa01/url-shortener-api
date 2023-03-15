package mocks

import "github.com/stretchr/testify/mock"

// mockRandomiser is a mock implementation of utils.Random
type mockRandomiser struct {
	mock.Mock
}

func NewMockRandomiser() mockRandomiser {
	return mockRandomiser{}
}

func (m mockRandomiser) GenerateRandomString(n int) string {
	ret := m.Called(n)
	return ret.String(0)
}
