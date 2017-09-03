package data

import (
	"time"
)

type Email struct {
	Subject     string        `json:"subject"`
	From        string        `json:"from"`
	To          []string      `json:"to"`
	CC          []string      `json:"cc"`
	CCI         []string      `json:"cci"`
	Date        time.Time     `json:"date"`
	ContentType string        `json:"contentType"`
	Attachments []*Attachment `json:"attachments"`
}

func (s *Email) GetName() string {
	return s.From + ":" + s.Subject
}

func (s *Email) GetRef() Ref {
	return EMAIL
}
