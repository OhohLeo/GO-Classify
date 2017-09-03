package data

import (
	"time"
)

type Attachment struct {
	Name               string    `json:"name"`
	Date               time.Time `json:"date"`
	ContentDisposition string    `json:"contentDisposition"`
	Data               []byte    `json:"data"`
}

func (s *Attachment) GetName() string {
	return s.Name
}

func (s *Attachment) GetRef() Ref {
	return ATTACHMENT
}
