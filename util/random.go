package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// TODO: テストパッケージに移す
func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	/* #nosec */
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(letters)

	for i := 0; i < n; i++ {
		/* #nosec */
		c := letters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomURL() string {
	return fmt.Sprintf("http://%s.com", RandomString(6))
}

func RandomThumbnailURL() string {
	return fmt.Sprintf("http://%s.com/thumbnail.jpg", RandomString(6))
}

func RandomHTML() string {
	return RandomString(100)
}
