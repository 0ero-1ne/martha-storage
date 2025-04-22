package utils

import (
	"errors"
	"math/rand/v2"
)

type Randomizer struct{}

func NewRandomizer() Randomizer {
	return Randomizer{}
}

/*
GenerateString creates pseudo random string from lower case english alphabet chars and numbers. Method returns error
if length <= 0
*/
func (randomizer Randomizer) GenerateString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("Invalid length param value")
	}

	alphabet := "abcdefghijklmnopqrstuvwxyz0123456789"
	randomString := make([]byte, length)

	for i := range length {
		randomString[i] = alphabet[rand.N(len(alphabet))]
	}

	return string(randomString[:]), nil
}
