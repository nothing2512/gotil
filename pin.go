package gotil

import "math/rand"

// generate pin with given length
func PIN(length int) string {
	chars := "0123456789"
	pinBytes := make([]byte, length)

	for i := range pinBytes {
		pinBytes[i] = chars[rand.Intn(len(chars))]
	}

	return string(pinBytes)
}
