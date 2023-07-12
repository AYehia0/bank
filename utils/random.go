package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	chars := make([]byte, length)
	for i := range chars {
		chars[i] = charset[rand.Intn(len(charset))]
	}
	return string(chars)
}

func randomStringFromList(list []string) string {
	n := len(list)
	return list[rand.Intn(n)]
}

func GetRandomOwnerName() string {
	return randomString(8)
}

func GetRandomAmount() int64 {
	return randomInt(0, 1000)
}

func GetRandomCurrency() string {
	return randomStringFromList([]string{"USD", "EUR", "EGP", "CAD"})
}
