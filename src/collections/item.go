package collections

import (
	"fmt"
	"github.com/ohohleo/classify/imports"
	"math/rand"
	"time"
)

const (
	CREATED = iota
	BUFFER_WAITING
	BUFFER_SENDING
	BUFFER_SENT
	BUFFER_OUT
)

type Data interface {
	GetType() string
}

type ItemGeneric struct {
	Status     int
	Type       string  `json:"type"`
	IsMatching float32 `json:"probability"`
}

type Item struct {
	ItemGeneric
	Id        string                    `json:"id"`
	Name      string                    `json:"name"`
	CreatedAt time.Time                 `json:"createAt"`
	Imports   map[string][]imports.Data `json:"imports"`
	Websites  map[string][]Data         `json:"websites"`
	BestMatch Data                      `json:"best_match"`
}

func NewItem() *Item {
	item := &Item{
		Id:        GetRandomId(),
		CreatedAt: time.Now(),
	}

	// Status init
	item.Status = CREATED

	return item
}

func (i *Item) AddImportData(data imports.Data) {
	if i.Imports == nil {
		i.Imports = make(map[string][]imports.Data)
	}

	if i.Name == "" {
		i.Name = data.String()
	}

	i.Imports[data.GetUniqKey()] = append(i.Imports[data.GetUniqKey()], data)
}

func (i *Item) RemoveImportData(data imports.Data) {
}

func (i *Item) AddWebsiteData(name string, data Data) {

	if i.Websites == nil {
		i.Websites = make(map[string][]Data)
	}

	i.Websites[name] = append(i.Websites[name], data)
}

func (i *Item) RemoveWebsiteData(data Data) {
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
