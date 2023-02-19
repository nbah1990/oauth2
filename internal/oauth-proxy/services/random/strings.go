package random

import (
	"math/rand"
)

func MakeString(length int) string {
	var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))] //nolint:gosec
	}
	return string(b)
}
