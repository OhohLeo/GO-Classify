package core

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/hydrogen18/stoppableListener"
	"github.com/manucorporat/sse"
	"net"
	"net/http"
)

type Server struct {
	api       *rest.Api
	config    *ServerConfig
	stoppable *stoppableListener.StoppableListener
	events    chan sse.Event
}

type ServerConfig struct {
	Url string `json:"url"`
}

type ProtocolReq struct {
	Collection string
	Type       string
	Data       interface{}
}

// ServerStart launches web server
func (c *Classify) CreateServer(config ServerConfig) (server *Server, err error) {

	server = new(Server)

	// Stockage de la configuration
	server.config = &config

	listener, err := net.Listen("tcp", config.Url)
	if err != nil {
		return
	}

	server.stoppable, err = stoppableListener.New(listener)
	if err != nil {
		return
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

	// Init events channel
	server.events = make(chan sse.Event, 1)

	router, err := rest.MakeRouter(

		// Establish connection to the web-services
		rest.Get("/stream", func(w rest.ResponseWriter, r *rest.Request) {
			for {
				event, ok := <-server.events
				if ok {
					sse.Encode(w.(http.ResponseWriter), event)
				}
			}
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
		rest.Get("/imports/config", c.ApiGetImportsConfig),
	)
	if err != nil {
		return
	}

	api.SetApp(router)

	// Store api
	server.api = api

	return
}

func (s *Server) Start() {

	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Println("Serving at " + s.config.Url)
	http.Serve(s.stoppable, s.api.MakeHandler())
}

// ServerStop stop web server
func (s *Server) Stop() {
	log.Println("Stop server")
	s.stoppable.Stop()
}

// OnEvent add new event on the event channel
func (s *Server) OnEvent(event string, id string, data interface{}) {
	s.events <- sse.Event{
		Event: event,
		Id:    id,
		Data:  data,
	}
}
