package data

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
)

type Attachment struct {
	Name               string          `json:"name"`
	Date               time.Time       `json:"date"`
	ContentDisposition string          `json:"contentDisposition"`
	File               *File           `json:"file"`
	Part               *multipart.Part `json:"-"`
}

func (s *Attachment) GetName() string {
	return s.Name
}

func (s *Attachment) GetRef() Ref {
	return ATTACHMENT
}

func (s *Attachment) GetDependencies() []Data {

	if s.File == nil {
		s.File = new(File)
	}

	return []Data{
		s.File,
	}
}

func (s *Attachment) GetContents() (contents map[string]string) {

	if s.File == nil {
		return
	}

	return s.File.GetContents()
}

func (s *Attachment) StoreToFile(path string) error {

	if s.Part == nil {
		return fmt.Errorf("No attachment '%s' to store at '%s'",
			s.Name, path)
	}

	// Get filename
	name := s.Date.Format("20060102_150405") + "_" + s.Name
	fullname := path + "/" + name

	// Check if file already exist
	if _, err := os.Stat(fullname); os.IsNotExist(err) == false {
		s.File, err = NewFileFromPath(path, name)
		return err
	}

	// Check if directory already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {

		// Otherwise create it
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	// Read all file
	slurp, err := ioutil.ReadAll(s.Part)
	if err != nil {
		return err
	}

	if len(slurp) <= 0 {
		return nil
	}

	fmt.Printf("Create '%s' with size %d\n", fullname, len(slurp))

	var toWrite []byte
	switch s.Part.Header.Get("Content-Transfer-Encoding") {

	case "base64":
		reader := bytes.NewReader(slurp)
		data := base64.NewDecoder(base64.StdEncoding, reader)
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(data)
		toWrite = buffer.Bytes()
	default:
		toWrite = slurp
	}

	// Write data into destination
	f, err := os.Create(fullname)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(toWrite)
	if err != nil {
		return err
	}

	s.File, err = NewFileFromOsFile(f)
	if err != nil {
		return err
	}

	return nil
}
