package core

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
func (c *Classify) CreateServer() (*Server, error) {

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
		rest.Get("/references", c.ApiGetReferences),

		// Handle collections
		rest.Post("/collections", c.ApiPostCollection),
		rest.Get("/collections", c.ApiGetCollections),
		rest.Get("/collections/:name", c.ApiGetCollectionByName),
		rest.Patch("/collections/:name", c.ApiPatchCollection),
		rest.Delete("/collections/:name", c.ApiDeleteCollectionByName),

		// Handle imports
		rest.Post("/imports", c.ApiAddImport),
		rest.Get("/imports", c.ApiGetImports),
		rest.Delete("/imports/:import", c.ApiDeleteImport),
		rest.Put("/imports/start", c.ApiStartImport),
		rest.Put("/imports/stop", c.ApiStopImport),
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

	log.Info("WebSocket Connection OK")
	s.websockets = append(s.websockets, ws)

	s.Read(ws)
}

func (s *Server) Read(ws *websocket.Conn) {

	msg := make([]byte, 512)
	_, err := ws.Read(msg)
	if err != nil {
		return
	}
}

// Send a request through websocket
func (s *Server) Send(collectionName string, itemType string, data interface{}) error {

	// No websockets : nothing to do
	if len(s.websockets) == 0 {
		return nil
	}

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
