package main

import (
	"errors"
	"golang.org/x/net/websocket"
)

type Interface func(ws *websocket.Conn, data interface{}) error

// OnNewDirectory handle new directory input interface
func OnNewDirectory(ws *websocket.Conn, data interface{}) error {

	path, ok := data.(string)
	if ok == false {
		return errors.New("Unexpected path!")
	}

	c, err := ReadDirectory(path, true)

	if err != nil {
		SendError(ws, err)
		return err
	}

	for {
		newFile, ok := <-c
		if ok == false {
			return nil
		}

		Send(ws, "newFile", newFile)
	}

	return nil
}
