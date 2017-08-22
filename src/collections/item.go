package collections

import (
	"math/rand"
	"time"
)

type Data interface {
	GetType() string
	GetId() string
}

type ItemGeneric struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Specific interface{} `json:"specific"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GetRandomId() string {

	size := 32
	src := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, size)

	for i, cache, remain := size-1, src.Int63(), letterIdxMax; i >= 0; {

		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits

		remain--
	}

	return string(b)
}
