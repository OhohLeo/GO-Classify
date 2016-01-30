package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
)

type Protocol struct {
	Type string
	Data interface{}
}

// Send a request through websocket
func Send(ws *websocket.Conn, t string, data interface{}) error {

	log.Printf(" <-- %s '%+v'", t, data)

	send, err := json.Marshal(Protocol{
		Type: t,
		Data: data,
	})

	if err != nil {
		log.Error("Can't convert msg: " + err.Error())
		return err
	}

	if err := websocket.Message.Send(ws, string(send)); err != nil {
		log.Error("Can't send msg: " + err.Error())
		return err
	}

	return nil
}

// Send an error request through websocket
func SendError(ws *websocket.Conn, err error) error {

	return Send(ws, "error", err.Error())
}

var interfaces = map[string]Interface{
	"newDirectory": OnNewDirectory,
}

// ServerStart launches web server
func ServerStart() {

	log.Println("Serving at localhost:8080...")

	http.Handle("/socket", websocket.Handler(func(ws *websocket.Conn) {

		log.Info("Connection OK")

		// Handling websocket
		for {

			var input string

			// Get received message
			if err := websocket.Message.Receive(ws, &input); err != nil {
				log.Error("Can't receive: " + err.Error())
				return
			}

			// Handle message
			var rcv Protocol

			if err := json.Unmarshal([]byte(input), &rcv); err != nil {
				log.Error("Unexpected message received")
				continue
			}

			log.Printf(" --> '%+v'", rcv)

			method, ok := interfaces[rcv.Type]
			if ok == false {
				log.Error("Unknown method")
				continue
			}

			if err := method(ws, rcv.Data); err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}))

	http.Handle("/", http.FileServer(http.Dir("www")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// ServerStop stop web server
func ServerStop() {
}
