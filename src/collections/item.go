package collections

import (
	"fmt"
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
	"math/rand"
	"time"
)

type Item struct {
	Id        string                     `json:"id"`
	CreatedAt time.Time                  `json:"createAt"`
	Imports   map[string]imports.Data    `json:"imports"`
	Websites  map[string][]websites.Data `json:"websites"`
}

func NewItem() *Item {
	return &Item{
		Id:        GetRandomId(),
		CreatedAt: time.Now(),
	}
}

func (i *Item) GetKeywords() string {
	return "Star+Wars"
}

func (i *Item) AddImportData(data imports.Data) {
	if i.Imports == nil {
		i.Imports = make(map[string]imports.Data)
	}

	i.Imports[data.GetUniqKey()] = data
}

func (i *Item) RemoveImportData(data imports.Data) {
}

func (i *Item) AddWebsiteData(name string, data websites.Data) {

	if i.Websites == nil {
		i.Websites = make(map[string][]websites.Data)
	}

	i.Websites[name] = append(i.Websites[name], data)
}

func (i *Item) RemoveWebsiteData(data websites.Data) {
}

func (i *Item) GetId() string {
	return i.Id
}

func (i *Item) String() string {
	return fmt.Sprintf("id: %s", i.Id)
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
