package core

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/foize/go.fifo"
	"github.com/hydrogen18/stoppableListener"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	api         *rest.Api
	config      *ServerConfig
	stoppable   *stoppableListener.StoppableListener
	events      *fifo.Queue
	eventStatus bool
	onEvents    chan Event
	onEvent     chan bool
}

type ServerConfig struct {
	Url string `json:"url"`
}

type ProtocolReq struct {
	Collection string
	Type       string
	Data       interface{}
}

type Event struct {
	Event  string      `json:"event"`
	Status string      `json:"status"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}

// ServerStart launches web server
func (c *Classify) CreateServer(config ServerConfig) (server *Server, err error) {

	if config.Url == "" {
		err = fmt.Errorf("No server configuration found!")
		return
	}

	server = new(Server)

	// Stockage de la configuration
	server.config = &config

	// Init events channel
	server.events = fifo.NewQueue()
	server.onEvents = make(chan Event)
	server.onEvent = make(chan bool)

	listener, err := net.Listen("tcp", config.Url)
	if err != nil {
		return
	}

	server.stoppable, err = stoppableListener.New(listener)
	if err != nil {
		return
	}

	// Store api
	server.api, err = c.GetApi(server)
	if err != nil {
		return
	}

	return
}

func (c *Classify) GetApi(server *Server) (*rest.Api, error) {

	api := rest.NewApi()

	api.Use(rest.DefaultDevStack...)

	// Enable CORS
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			fmt.Printf("REQUEST %s %+v\n", origin, request.URL)
			return true
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	var streamRoute *rest.Route
	if server != nil {
		streamRoute = rest.Get("/stream", server.HandleStream)
	} else {
		streamRoute = rest.Get("/stream", nil)
	}

	router, err := rest.MakeRouter(

		// Handle references
		rest.Get("/references", c.ApiGetReferences),

		// Handle imports
		rest.Post("/imports", c.ApiAddImport),
		rest.Get("/imports", c.ApiGetImports),
		rest.Delete("/imports", c.ApiDeleteImport),
		rest.Put("/imports/start", c.ApiStartImport),
		rest.Put("/imports/stop", c.ApiStopImport),
		rest.Get("/imports/:name/references", c.ApiGetImportReferences),
		rest.Get("/imports/:name/config", c.ApiGetImportConfig),
		rest.Patch("/imports/:name/config", c.ApiPatchImportConfig),
		rest.Put("/imports/:name/:param", c.ApiPutImportParam),
		rest.Put("/imports/:name/tweak", c.ApiPutImportTweaks),
		rest.Get("/imports/:name/tweak", c.ApiGetImportTweaks),

		// Handle exports
		rest.Post("/exports", c.ApiAddExport),
		rest.Get("/exports", c.ApiGetExports),
		rest.Delete("/exports", c.ApiDeleteExport),
		rest.Put("/exports/force", c.ApiForceExport),
		rest.Put("/exports/stop", c.ApiStopExport),
		rest.Get("/exports/:name/config", c.ApiGetExportConfig),
		rest.Patch("/exports/:name/config", c.ApiPatchExportConfig),
		rest.Put("/exports/:name/:param", c.ApiPutExportParam),

		// Handle collections
		rest.Post("/collections", c.ApiPostCollection),
		rest.Get("/collections", c.ApiGetCollections),
		rest.Get("/collections/:name", c.ApiGetCollectionByName),
		rest.Patch("/collections/:name", c.ApiPatchCollection),
		rest.Delete("/collections/:name", c.ApiDeleteCollectionByName),
		rest.Get("/collections/:name/config", c.ApiGetCollectionConfig),
		rest.Patch("/collections/:name/config", c.ApiPatchCollectionConfig),
		rest.Put("/collections/:name/config/:path/:param", c.ApiPutCollectionConfigParam),

		// Handle collection buffer
		rest.Get("/collections/:name/buffers", c.ApiGetCollectionBuffers),
		rest.Delete("/collections/:name/buffers", c.ApiDeleteCollectionBuffers),
		rest.Get("/collections/:name/buffers/:id", c.ApiGetCollectionSingleBuffer),
		rest.Patch("/collections/:name/buffers/:id", c.ApiPatchCollectionSingleBuffer),
		rest.Delete("/collections/:name/buffers/:id", c.ApiDeleteCollectionSingleBuffer),
		rest.Post("/collections/:name/buffers/:id/validate", c.ApiValidateCollectionSingleBuffer),

		// Handle collection items
		rest.Get("/collections/:name/items", c.ApiGetCollectionItems),
		rest.Delete("/collections/:name/items", c.ApiDeleteCollectionItems),
		rest.Get("/collections/:name/items/:id", c.ApiGetCollectionSingleItem),
		rest.Patch("/collections/:name/items/:id", c.ApiPatchCollectionSingleItem),
		rest.Delete("/collections/:name/items/:id", c.ApiDeleteCollectionSingleItem),

		// Establish connection to the web-services
		streamRoute,
	)

	if err != nil {
		return nil, err
	}

	api.SetApp(router)

	return api, nil
}

func (c *Classify) SendEvent(event string, status string, name string, data interface{}) {
	c.Server.SendEvent(event, status, name, data)
}

func (s *Server) Start() {

	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Println("Serving at " + s.config.Url)
	http.Serve(s.stoppable, s.api.MakeHandler())
}

// ServerStop stop web server
func (s *Server) Stop() {
	log.Println("Stop server")
	s.onEvent <- false
	s.stoppable.Stop()
}

// SendEvent add new event on the event channel
func (s *Server) SendEvent(eventType string, status string, name string, data interface{}) {

	fmt.Printf("SEND EVENT [%s] %s name:%s data:%+v\n", status, eventType, name, data)

	event := Event{
		Event:  eventType,
		Status: status,
		Name:   name,
		Data:   data,
	}

	s.events.Add(event)

	// If event streamer is not emptying fifo
	if s.eventStatus == false {

		// Kick event streamer to send this new event
		s.onEvent <- true
		s.eventStatus = true
	}
}

var idx int = 0

func (s *Server) HandleStream(w rest.ResponseWriter, r *rest.Request) {

	// Get flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		rest.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Prepare write response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Close notification
	notify := w.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-notify:
			continue
		case ok := <-s.onEvent:

			if ok == false {
				continue
			}

			for {
				event, ok := s.events.Next().(Event)
				if ok == false {
					break
				}

				eventJson, err := json.Marshal(event)
				if err != nil {
					rest.Error(w, "Encoding event error", http.StatusInternalServerError)
					return
				}

				fmt.Fprintf(w.(http.ResponseWriter), "data: %s\n\n", eventJson)

				// Send data immediately
				flusher.Flush()

				// Wait for 500 ms
				time.Sleep(500 * time.Millisecond)
			}

			s.onEvent <- false
			s.eventStatus = false
			continue
		}
	}
}
