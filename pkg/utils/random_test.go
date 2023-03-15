package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateRandomString(t *testing.T) {
	stringMap := map[string]struct{}{}
	r := NewRandomiser()

	for lengthOfString := 6; lengthOfString < 12; lengthOfString++ {
		for numToGenerate := 1000; numToGenerate >= 0; numToGenerate-- {

			generatedString := r.GenerateRandomString(lengthOfString)
			_, exists := stringMap[generatedString]

			assert.Falsef(t, exists, "%s already exists", generatedString) // string shouldn't alrteady exist
			assert.Equal(t, lengthOfString, len(generatedString))          // generated string should be anticipated length

			stringMap[generatedString] = struct{}{}
		}
	}
}

func Test_GenerateRandomString_LengthSix(t *testing.T) {
	stringMap := map[string]struct{}{}
	length := 6
	r := NewRandomiser()

	for numToGenerate := 100000; numToGenerate >= 0; numToGenerate-- {

		generatedString := r.GenerateRandomString(length)
		_, exists := stringMap[generatedString]

		assert.Falsef(t, exists, "%s already exists", generatedString) // string shouldn't alrteady exist
		assert.Equal(t, length, len(generatedString))                  // generated string should be anticipated length

		stringMap[generatedString] = struct{}{}
	}

}
