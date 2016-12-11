package collections

import (
	"fmt"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"math/rand"
	"time"
)

type Item struct {
	id        string
	createdAt time.Time
	imports   map[string]imports.Data
	websites  map[string]websites.Data
}

func NewItem() *Item {
	return &Item{
		id:        GetRandomId(),
		createdAt: time.Now(),
	}
}

func (i *Item) GetKeywords() string {
	return "Star+Wars"
}

func (i *Item) AddImportData(data imports.Data) {
	if i.imports == nil {
		i.imports = make(map[string]imports.Data)
	}

	i.imports[data.GetUniqKey()] = data
}

func (i *Item) RemoveImportData(data imports.Data) error {
	return nil
}

func (i *Item) AddWebsiteData(data websites.Data) error {
	return nil
}

func (i *Item) RemoveWebsiteData(data websites.Data) error {
	return nil
}

func (i *Item) GetId() string {
	return i.id
}

func (i *Item) String() string {
	return fmt.Sprintf("id: %s", i.id)
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
