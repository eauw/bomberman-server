package main

import (
	"math/rand"
	"time"
)

// returns a random string of values given in the letters array
func randomString(length int) string {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	lettersLength := len(letters)

	randomString := ""

	for i := 0; i < length; i++ {
		randomLetter := letters[randomNumber(0, lettersLength)]
		randomString += randomLetter
	}

	return randomString
}

// returns a random integer between min-value and max-value
func randomNumber(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Intn(max-min) + min
}
