package collections

import (
	"fmt"
	"github.com/ohohleo/classify/imports"
	"strings"
	"time"
)

const (
	CREATED = iota
	BUFFER_WAITING
	BUFFER_SENDING
	BUFFER_SENT
	BUFFER_VALIDATE
)

type BufferItem struct {
	ItemGeneric
	Status      int                       `json:"status"`
	CleanedName string                    `json:"cleanedName"`
	CreatedAt   time.Time                 `json:"createAt"`
	Banned      []string                  `json:"banned"`
	Separators  []string                  `json:"separators"`
	Imports     map[string][]imports.Data `json:"imports"`
	WebQuery    string                    `json:"webQuery"`
	Websites    map[string][]Data         `json:"websites"`
	MatchId     string                    `json:"matchId"`
	IsMatching  float32                   `json:"probability"`
	Data        Data                      `json:"data"`
}

func NewBufferItem() *BufferItem {
	item := &BufferItem{
		CreatedAt: time.Now(),
	}

	item.ItemGeneric.Id = GetRandomId()

	// Status init
	item.Status = CREATED

	return item
}

func (i *BufferItem) SetCleanedName(bannedList []string, separators []string) bool {

	previousName := i.CleanedName

	i.CleanedName = i.Name
	i.Banned = make([]string, 0)
	i.Separators = make([]string, 0)

	// No real name set : nothing todo
	if i.Name == "" {
		return false
	}

	// Separate elements with separators
	for _, separator := range separators {

		if strings.Contains(i.CleanedName, separator) == false {
			continue
		}

		i.CleanedName = strings.Replace(i.CleanedName, separator, " ", -1)
		i.Separators = append(i.Separators, separator)
	}

	// Removed banned strings
	for _, banned := range bannedList {

		if strings.Contains(i.CleanedName, banned) == false {
			continue
		}

		i.CleanedName = strings.Replace(i.CleanedName, banned, "", -1)
		i.Banned = append(i.Banned, banned)

	}

	// Cleaned spaces
	i.CleanedName = strings.Replace(i.CleanedName, "  ", " ", -1)
	i.CleanedName = strings.TrimSpace(i.CleanedName)

	// Set web query
	i.WebQuery = strings.Replace(i.CleanedName, " ", "+", -1)

	return previousName != i.CleanedName
}

func (i *BufferItem) AddImportData(data imports.Data) {

	if i.Imports == nil {
		i.Imports = make(map[string][]imports.Data)
	}

	if i.Name == "" {
		i.Name = data.String()
	}

	i.Imports[data.GetUniqKey()] = append(i.Imports[data.GetUniqKey()], data)
}

func (i *BufferItem) RemoveImportData(data imports.Data) {
}

func (i *BufferItem) AddWebsiteData(name string, data Data) {

	if i.Websites == nil {
		i.Websites = make(map[string][]Data)
	}

	i.Websites[name] = append(i.Websites[name], data)
}

func (i *BufferItem) RemoveWebsiteData(data Data) {
}

func (i *BufferItem) GetType() string {
	return "buffer"
}

func (i *BufferItem) GetId() string {
	return i.Id
}

func (i *BufferItem) String() string {
	return fmt.Sprintf("id: %s", i.Id)
}
