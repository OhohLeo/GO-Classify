package core

import (
	"fmt"
	"github.com/ohohleo/classify/data"
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
	Item
	Status      int                    `json:"status"`
	CleanedName string                 `json:"cleanedName"`
	CreatedAt   time.Time              `json:"createAt"`
	Banned      []string               `json:"banned"`
	Separators  []string               `json:"separators"`
	Imports     []data.Data            `json:"imports"`
	WebQuery    string                 `json:"webQuery"`
	Websites    map[string][]data.Data `json:"websites"`
	MatchId     string                 `json:"matchId"`
	IsMatching  float32                `json:"probability"`
}

func NewBufferItem(id Id) *BufferItem {
	item := &BufferItem{
		CreatedAt: time.Now(),
	}

	item.Item.Id = id
	item.Status = CREATED

	return item
}

func (i *BufferItem) SetCleanedName(bannedList []string, separators []string) bool {

	previousName := i.CleanedName

	i.CleanedName = i.Item.Engine.GetName()
	i.Banned = make([]string, 0)
	i.Separators = make([]string, 0)

	// No real name set : nothing todo
	if i.Item.Engine.GetName() == "" {
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

func (i *BufferItem) AddImportData(d data.Data) {
	i.Imports = append(i.Imports, d)
}

func (i *BufferItem) RemoveImportData(d data.Data) {
}

func (i *BufferItem) AddWebsiteData(name string, d data.Data) {

	if i.Websites == nil {
		i.Websites = make(map[string][]data.Data)
	}

	i.Websites[name] = append(i.Websites[name], d)
}

func (i *BufferItem) RemoveWebsiteData(d data.Data) {
}

func (i *BufferItem) String() string {
	return fmt.Sprintf("id: %s", i.Id)
}
