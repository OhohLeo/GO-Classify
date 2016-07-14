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

type Server struct {
	api        *rest.Api
	stoppable  *stoppableListener.StoppableListener
	websockets []*websocket.Conn
}

type ProtocolReq struct {
	Collection string
	Type       string
	Data       interface{}
}

// ServerStart launches web server
func CreateServer() (*Server, error) {

	server := new(Server)

	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		return nil, err
	}

	server.stoppable, err = stoppableListener.New(listener)
	if err != nil {
		return nil, err
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
	wsHandler := websocket.Handler(server.handleWebSocket)

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
		return nil, err
	}

	api.SetApp(router)

	// Store api
	server.api = api

	return server, nil
}

func (s *Server) Start() {

	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Println("Serving at localhost:3333...")
	http.Serve(s.stoppable, s.api.MakeHandler())
}

// ServerStop stop web server
func (s *Server) Stop() {
	log.Println("Stop server")
	s.stoppable.Stop()
}

// handleWebSocket store all connections established
func (s *Server) handleWebSocket(ws *websocket.Conn) {

	if s.websockets == nil {
		s.websockets = make([]*websocket.Conn, 0)
	}

	log.Info("Connection OK")
	s.websockets = append(s.websockets, ws)
}

// Send a request through websocket
func (s *Server) Send(collectionName string, itemType string, data interface{}) error {

	log.Printf(" <-- [%s] %s '%+v'", collectionName, itemType, data)

	send, err := json.Marshal(ProtocolReq{
		Collection: collectionName,
		Type:       itemType,
		Data:       data,
	})

	if err != nil {
		log.Error("Can't convert msg: " + err.Error())
		return err
	}

	for _, ws := range s.websockets {

		if ws.IsClientConn() == false {
			log.Error("Client not connected!")
			continue
		}

		if err := websocket.Message.Send(ws, string(send)); err != nil {
			log.Error("Can't send msg: " + err.Error())
		}
	}

	return nil
}

// Send an error request through websocket
func (s *Server) SendError(collectionName string, err error) error {

	return s.Send(collectionName, "error", err.Error())
}
