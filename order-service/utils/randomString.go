package utils

import "math/rand"

func RandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}
