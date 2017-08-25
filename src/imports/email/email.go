package email

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/ohohleo/classify/imports"
	"log"
)

type SearchCriteria struct {
}

type Email struct {
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Login    string         `json:"login"`
	Password string         `json:"password"`
	Search   SearchCriteria `json:"search"`

	cnx       *client.Client
	isRunning bool
}

func ToBuild() imports.BuildImport {
	return imports.BuildImport{
		Create: func(input json.RawMessage, config map[string][]string, collections []string) (i imports.Import, err error) {
			var email Email
			err = json.Unmarshal(input, &email)
			i = &email
			return
		},
	}
}

func (e *Email) GetRef() imports.Ref {
	return imports.EMAIL
}

func (e *Email) Check(config map[string][]string, collections []string) error {
	return nil
}

func (e *Email) Start() (chan imports.Data, error) {

	c := make(chan imports.Data)

	// Check if the analysis is not already going on
	if e.isRunning {
		return c, fmt.Errorf("import 'email' already started")
	}

	address := fmt.Sprintf("%s:%d", e.Host, e.Port)

	log.Printf("Connecting to '%s'...\n", address)

	// Connect to IMAP server
	cnx, err := client.DialTLS(address, nil)
	if err != nil {
		return c, fmt.Errorf("import 'email' connection: %s", err.Error())
	}

	// Login
	if err := cnx.Login(e.Login, e.Password); err != nil {
		return c, fmt.Errorf("import 'email' login: %s", err.Error())
	}

	// Store connection
	e.cnx = cnx

	return c, nil
}

func (e *Email) Stop() {

	// No need to close unitialised connection
	if e.cnx == nil {
		return
	}

	e.cnx.Logout()
}

func (e *Email) Eq(new imports.Import) bool {
	newEmail, _ := new.(*Email)
	return e.Host == newEmail.Host &&
		e.Port == newEmail.Port &&
		e.Login == newEmail.Login
}

// func main() {

// 	log.Println("Connected")

// 	// Don't forget to logout
// 	defer

// 	// List mailboxes
// 	mailboxes := make(chan *imap.MailboxInfo, 10)
// 	done := make(chan error, 1)
// 	go func() {
// 		done <- c.List("", "*", mailboxes)
// 	}()

// 	log.Println("Mailboxes:")
// 	for m := range mailboxes {
// 		log.Println("* " + m.Name)
// 	}

// 	if err := <-done; err != nil {
// 		log.Fatal(err)
// 	}

// 	// Select INBOX
// 	mbox, err := c.Select("INBOX", false)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Flags for INBOX:", mbox.Flags)

// 	// Search
// 	criteria := &imap.SearchCriteria{
// 		Larger: 1000,
// 		Body:   []string{"salaire"},
// 	}

// 	seqNums, err := c.Search(criteria)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("SeqNums:%+v\n", seqNums)

// 	// // Get the last 4 messages
// 	// from := uint32(1)
// 	// to := mbox.Messages
// 	// if mbox.Messages > 3 {
// 	// 	// We're using unsigned integers here, only substract if the result is > 0
// 	// 	from = mbox.Messages - 3
// 	// }

// 	seqset := new(imap.SeqSet)
// 	seqset.AddNum(seqNums...)

// 	messages := make(chan *imap.Message, len(seqNums))

// 	done = make(chan error, 1)

// 	go func() {
// 		done <- c.Fetch(seqset, []string{imap.EnvelopeMsgAttr}, messages)
// 	}()

// 	log.Println("Messages:")
// 	for msg := range messages {
// 		log.Println("* " + msg.Envelope.Subject)
// 	}

// 	if err := <-done; err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Done!")
// }
