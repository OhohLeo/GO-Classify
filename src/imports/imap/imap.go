package imap

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jhillyerd/go.enmime"
	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/imports"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

const (
	SEARCH Request = iota
	ALL
)

type Request int

type ImapConfig struct {
	Store struct {
		Path string `json:"path"`
	} `json:"store"`
}

func DefaultConfig() *ImapConfig {

	config := &ImapConfig{}

	// Default store/imap config
	config.Store.Path = "/tmp/classify/imports/imap"

	return config
}

type Imap struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`

	Request      Request `json:"request"`
	MailBox      string  `json:"mailbox"`
	OnlyAttached bool    `json:"onlyAttached"`
	Search       Search  `json:"search"`

	dataChannel chan data.Data
	cnx         *client.Client
	config      *ImapConfig
}

type Search struct {
	Since   time.Time `json:"since"`
	Before  time.Time `json:"before"`
	Larger  uint32    `json:"larger"`
	Smaller uint32    `json:"smaller"`
	Text    []string  `json:"text"`
}

func (s *Search) IsValid() bool {
	return true
}

type ImapOutputParams struct {
	MailBoxes []string `json:"mailboxes"`
}

func ToBuild() imports.Build {
	return imports.Build{
		Create: Create,
	}
}

func Create(input json.RawMessage, config json.RawMessage,
	collections []string) (i imports.Import, params interface{}, err error) {

	var imap Imap
	err = json.Unmarshal(input, &imap)

	// MailBox is required
	if imap.MailBox == "" {

		// Check connection
		err = imap.Connect()
		if err != nil {
			return
		}

		// Returns mailbox
		var mailboxes []string
		mailboxes, err = imap.GetMailBoxes()
		if err != nil {
			return
		}

		params = &ImapOutputParams{
			MailBoxes: mailboxes,
		}

		imap.Stop()
		err = fmt.Errorf("import 'imap' needs more params")
		return
	}

	switch imap.Request {
	case SEARCH:
		if imap.Search.IsValid() == false {
			err = fmt.Errorf("import 'imap' invalid search params")
			return
		}
	case ALL:
	}

	if imap.config == nil {
		imap.config = DefaultConfig()
	}

	i = &imap

	return
}

func (i *Imap) GetRef() imports.Ref {
	return imports.IMAP
}

func (i *Imap) CheckConfig(config json.RawMessage) error {
	return nil
}

func (i *Imap) GetDataList() []data.Data {
	return []data.Data{
		new(data.Email),
		new(data.Attachment),
	}
}

func (i *Imap) GetParam(string, json.RawMessage) (result interface{}, err error) {
	return
}

func (i *Imap) Start() (c chan data.Data, err error) {

	// Check if the analysis is not already going on
	if i.cnx != nil {
		err = fmt.Errorf("import 'imap' already started")
		return
	}

	// Establish connection
	if err = i.Connect(); err != nil {
		return
	}

	switch i.Request {
	case SEARCH:
		go i.GetSearch()
	case ALL:
		go i.GetAllMessages()
	}

	c = make(chan data.Data)

	i.dataChannel = c

	return
}

func (i *Imap) Stop() error {

	// No need to close unitialised connection
	if i.cnx == nil {
		return fmt.Errorf("import 'imap' already stopped")
	}

	i.cnx.Logout()
	i.cnx = nil
	return nil
}

func (i *Imap) Eq(new imports.Import) bool {
	newImap, _ := new.(*Imap)
	return i.Host == newImap.Host &&
		i.Port == newImap.Port &&
		i.Login == newImap.Login
}

func (i *Imap) Connect() error {

	address := fmt.Sprintf("%s:%d", i.Host, i.Port)

	log.Printf("Connecting to '%s'...\n", address)

	// Connect to IMAP server
	cnx, err := client.DialTLS(address, nil)
	if err != nil {
		return fmt.Errorf("import 'imap' connection: %s", err.Error())
	}

	// Login
	if err := cnx.Login(i.Login, i.Password); err != nil {
		return fmt.Errorf("import 'imap' login: %s", err.Error())
	}

	log.Printf("'%s' Connected!\n", address)

	// Store connection
	i.cnx = cnx

	return nil
}

// Permet la gestion des commandes asynchrones
func (i *Imap) GetMailBoxes() (mailboxes []string, err error) {

	if i.cnx == nil {
		err = fmt.Errorf("imap uninitialised")
		return
	}

	// List mailboxes
	infos := make(chan *imap.MailboxInfo, 10)
	done := make(chan error)

	go func() {
		done <- i.cnx.List("", "*", infos)
	}()

	if err = <-done; err != nil {
		i.Stop()
		return
	}

	for m := range infos {
		mailboxes = append(mailboxes, m.Name)
	}

	return
}

func (i *Imap) GetAllMessages() error {

	mailbox, err := i.cnx.Select(i.MailBox, false)
	if err != nil {
		return err
	}

	// Get all messages
	from := uint32(1)
	to := mailbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	return i.Proceed(seqset)
}

func (i *Imap) GetSearch() error {

	_, err := i.cnx.Select(i.MailBox, false)
	if err != nil {
		i.Stop()
		return err
	}

	criteria := &imap.SearchCriteria{
		Since:   i.Search.Since,
		Before:  i.Search.Before,
		Larger:  i.Search.Larger,
		Smaller: i.Search.Smaller,
		Text:    i.Search.Text,
	}

	// Launch research
	seqNums, err := i.cnx.Search(criteria)
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
		i.Stop()
		return err
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	return i.Proceed(seqset)
}

func (i *Imap) Proceed(seqset *imap.SeqSet) error {

	messages := make(chan *imap.Message, 10)
	done := make(chan error)

	go func() {
		done <- i.cnx.Fetch(seqset, []string{"BODY[]"}, messages)
	}()

	if err := <-done; err != nil {
		i.Stop()
		return err
	}

	for msg := range messages {

		rsp := msg.GetBody("BODY[]")
		if rsp == nil {
			fmt.Println("Server didn't returned message body")
			continue
		}

		m, err := mail.ReadMessage(rsp)
		if err != nil {
			i.Stop()
			return err
		}

		header := m.Header

		date, err := header.Date()
		if err != nil {
			i.Stop()
			return err
		}

		email := &data.Email{
			Subject:     header.Get("Subject"),
			From:        header.Get("From"),
			Date:        date,
			ContentType: header.Get("Content-Type"),
		}

		mediaType, params, err := mime.ParseMediaType(
			m.Header.Get("Content-Type"))
		if err != nil {
			i.Stop()
			return err
		}

		if strings.HasPrefix(mediaType, "multipart/") {

			mr := multipart.NewReader(m.Body, params["boundary"])

			email.Attachments = make([]*data.Attachment, 0)

			idx := 0
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					i.Stop()
					return err
				}

				name := Convert2Utf8()
				if name == "" {
					name = fmt.Sprintf("Part%d", idx)
				}

				attachment := &data.Attachment{
					Part:               p,
					Name:               name,
					ContentDisposition: p.FormName(),
					Date:               date,
				}

				if err := attachment.StoreToFile(i.config.Store.Path); err != nil {
					log.Printf("Issue storing %s to '%s': %s\n", name,
						i.config.Store.Path, err.Error())
				}

				email.Attachments = append(email.Attachments, attachment)
				i.dataChannel <- attachment
				idx++
			}
		}

		i.dataChannel <- email
	}

	close(i.dataChannel)
	i.Stop()
	return nil
}

var REMOVE_ISO = regexp.MustCompile(`=\?([a-zA-Z0-9-_]+)\?[a-zA-Z]\?([^\?]+)\?=`)

func Convert2Utf8(name string) string {

	submatches := REMOVE_ISO.FindAllStringSubmatch(name, -1)
	if len(submatches) == 1 && len(submatches[0]) == 3 {
		name = submatches[0][2]
	}

	converted, err := ioutil.ReadAll(quotedprintable.NewReader(
		strings.NewReader(name)))
	if err == nil {
		name = string(converted)
	}

	return name
}
