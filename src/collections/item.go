package collections

import (
	"fmt"
	"github.com/ohohleo/classify/imports"
	"math/rand"
	"strings"
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
	Status     int     `json:"status"`
	Type       string  `json:"type"`
	IsMatching float32 `json:"probability"`
}

type Item struct {
	ItemGeneric
	Id          string                    `json:"id"`
	Name        string                    `json:"name"`
	CleanedName string                    `json:"cleanedName"`
	CreatedAt   time.Time                 `json:"createAt"`
	Banned      []string                  `json:"banned"`
	Separators  []string                  `json:"separators"`
	Imports     map[string][]imports.Data `json:"imports"`
	WebQuery    string                    `json:"webQuery"`
	Websites    map[string][]Data         `json:"websites"`
	BestMatch   Data                      `json:"bestMatch"`
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

func (i *Item) SetCleanedName(bannedList []string, separators []string) bool {

	previousName := i.CleanedName

	i.CleanedName = i.Name
	i.Banned = make([]string, 0)
	i.Separators = make([]string, 0)

	// No real name set : nothing todo
	if i.Name == "" {
		return false
	}

	// Removed banned strings
	for _, banned := range bannedList {

		if strings.Contains(i.CleanedName, banned) == false {
			continue
		}

		i.CleanedName = strings.Replace(i.CleanedName, banned, "", -1)
		i.Banned = append(i.Banned, banned)
	}

	// Separate elements with separators
	for _, separator := range separators {

		if strings.Contains(i.CleanedName, separator) == false {
			continue
		}

		i.CleanedName = strings.Replace(i.CleanedName, separator, " ", -1)
		i.Separators = append(i.Separators, separator)
	}

	i.CleanedName = strings.TrimSpace(i.CleanedName)

	return previousName != i.CleanedName
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
