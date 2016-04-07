package helper

import (
	"math/rand"
	"time"
)

// returns a random string of values given in the letters array
func RandomString(length int) string {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	lettersLength := len(letters)

	randomString := ""

	for i := 0; i < length; i++ {
		randomLetter := letters[RandomNumber(0, lettersLength)]
		randomString += randomLetter
	}

	return randomString
}

// returns a random integer between min-value and max-value
func RandomNumber(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Intn(max-min) + min
}

// Erzeugt eine 8-stellige PlayerID
func GeneratePlayerID() string {
	return RandomString(8)
}
