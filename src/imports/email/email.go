package email

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/ohohleo/classify/imports"
	"log"
)

type Email struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
	MailBox  string `json:"mailbox"`

	dataChannel chan imports.Data
	cnx         *client.Client
	isRunning   bool
}

type SearchCriteria struct {
}

type EmailOutputParams struct {
	MailBoxes []string `json:"mailboxes"`
}

func ToBuild() imports.BuildImport {
	return imports.BuildImport{
		Create: Create,
	}
}

func Create(input json.RawMessage, config map[string][]string, collections []string) (i imports.Import, params interface{}, err error) {

	var email Email
	err = json.Unmarshal(input, &email)

	// if no mail box specified :
	if email.MailBox == "" {

		// Check connection
		err = email.Connect()
		if err != nil {
			return
		}

		// Returns mailbox
		var mailboxes []string
		mailboxes, err = email.GetMailBoxes()
		if err != nil {
			return
		}

		params = &EmailOutputParams{
			MailBoxes: mailboxes,
		}

		email.Stop()
		err = fmt.Errorf("import 'email' needs more params")
		return
	}

	i = &email

	return
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

	e.dataChannel = c

	return c, nil
}

func (e *Email) Stop() {

	// No need to close unitialised connection
	if e.cnx == nil {
		return
	}

	e.cnx.Logout()
	e.cnx = nil
}

func (e *Email) Eq(new imports.Import) bool {
	newEmail, _ := new.(*Email)
	return e.Host == newEmail.Host &&
		e.Port == newEmail.Port &&
		e.Login == newEmail.Login
}

func (e *Email) Connect() error {

	address := fmt.Sprintf("%s:%d", e.Host, e.Port)

	log.Printf("Connecting to '%s'...\n", address)

	// Connect to IMAP server
	cnx, err := client.DialTLS(address, nil)
	if err != nil {
		return fmt.Errorf("import 'email' connection: %s", err.Error())
	}

	// Login
	if err := cnx.Login(e.Login, e.Password); err != nil {
		return fmt.Errorf("import 'email' login: %s", err.Error())
	}

	log.Printf("'%s' Connected!\n", address)

	// Store connection
	e.cnx = cnx

	return nil
}

// Permet la gestion des commandes asynchrones
func (e *Email) GetMailBoxes() (mailboxes []string, err error) {

	if e.cnx == nil {
		err = fmt.Errorf("email uninitialised")
		return
	}

	// List mailboxes
	infos := make(chan *imap.MailboxInfo, 10)

	done := make(chan error, 1)
	go func() {
		done <- e.cnx.List("", "*", infos)
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("infos:")
	for m := range infos {
		mailboxes = append(mailboxes, m.Name)
		log.Println("* " + m.Name)
	}

	return
}

// Permet la gestion des commandes asynchrones
func (e *Email) Search() error {

	if e.cnx == nil {
		return fmt.Errorf("email uninitialised")
	}

	return nil
}

// func main() {

// 	log.Println("Connected")

// 	// Don't forget to logout
// 	defer

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
