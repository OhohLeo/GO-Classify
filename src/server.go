package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/hydrogen18/stoppableListener"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
)

var stoppable *stoppableListener.StoppableListener

type ProtocolReq struct {
	Type string
	Data interface{}
}

type ProtocolRcv struct {
	Type string
	Data json.RawMessage
}

// Send a request through websocket
func Send(ws *websocket.Conn, t string, data interface{}) error {

	log.Printf(" <-- %s '%+v'", t, data)

	send, err := json.Marshal(ProtocolReq{
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

var websockets = map[string]Websocket{
	"new-directory": new(NewDirectory),
}

// ServerStart launches web server
func ServerStart() error {

	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		return err
	}

	stoppable, err = stoppableListener.New(listener)
	if err != nil {
		return err
	}

	api := rest.NewApi()

	api.Use(rest.DefaultDevStack...)

	// Enable CORS
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods:                []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:                []string{"Accept", "Content-Type", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	// Init websocket
	wsHandler := websocket.Handler(handleWebSocket)

	router, err := rest.MakeRouter(

		// Establish connection to the web-services
		rest.Get("/ws", func(w rest.ResponseWriter, r *rest.Request) {
			wsHandler.ServeHTTP(w.(http.ResponseWriter), r.Request)
		}),

		// Handle references
		rest.Get("/references", ApiGetReferences),

		// Handle collections
		rest.Post("/collections", ApiPostCollection),
		rest.Get("/collections", ApiGetCollections),
		rest.Get("/collections/:name", ApiGetCollectionByName),
		rest.Patch("/collections/:name", ApiPatchCollection),
		rest.Delete("/collections/:name", ApiDeleteCollectionByName),

		rest.Put("/collections/:name/start",
			ApiStartCollection),
		rest.Put("/collections/:name/stop",
			ApiStopCollection),

		// Handle collection's imports
		rest.Post("/collections/:name/imports", ApiPostCollectionImport),
		rest.Get("/collections/:name/imports", ApiGetCollectionImports),
		rest.Delete("/collections/:name/imports/:import",
			ApiDeleteCollectionImport),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Println("Serving at localhost:3333...")
	http.Serve(stoppable, api.MakeHandler())

	return nil
}

// ServerStop stop web server
func ServerStop() {
	log.Println("Stop server")
	stoppable.Stop()
}

// handleWebSocket handles websockets requests & responses
func handleWebSocket(ws *websocket.Conn) {

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
		var rcv ProtocolRcv

		if err := json.Unmarshal([]byte(input), &rcv); err != nil {
			log.Error("Unexpected message received")
			continue
		}

		rsp, ok := websockets[rcv.Type]
		if ok == false {
			log.Error("Unknown method")
			continue
		}

		if err := json.Unmarshal([]byte(rcv.Data), &rsp); err != nil {
			log.Error("Unexpected message received")
			continue
		}

		log.Printf(" --> '%s' %+v", rcv.Type, rsp)

		if err := rsp.Handle(ws); err != nil {
			log.Error(err.Error())
			continue
		}
	}
}
