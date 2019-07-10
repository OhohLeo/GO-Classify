package api

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
	url         string
	stoppable   *stoppableListener.StoppableListener
	events      *fifo.Queue
	eventStatus bool
	onEvent     chan interface{}
	onStop      chan bool
}

type ProtocolReq struct {
	Collection string
	Type       string
	Data       interface{}
}

// ServerStart launches web server
func (a *API) NewServer(url string) (server *Server, err error) {
	if url == "" {
		err = fmt.Errorf("no server URL found")
		return
	}

	server = new(Server)

	// Init events channel
	server.events = fifo.NewQueue()
	server.onEvent = make(chan interface{})
	server.onStop = make(chan bool)

	listener, err := net.Listen("tcp", url)
	if err != nil {
		return
	}

	server.stoppable, err = stoppableListener.New(listener)
	if err != nil {
		return
	}

	// Store api
	server.api, err = a.GetAPI(server)
	if err != nil {
		return
	}

	server.url = url
	return
}

func (a *API) GetAPI(server *Server) (*rest.Api, error) {

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
		rest.Get("/references", a.GetReferences),

		// Handle imports
		rest.Post("/imports", a.AddImport),
		rest.Get("/imports", a.GetImports),
		rest.Delete("/imports", a.DeleteImport),
		rest.Put("/imports/start", a.StartImport),
		rest.Put("/imports/stop", a.StopImport),
		rest.Get("/imports/:name/references", a.GetImportReferences),
		rest.Get("/imports/:name/config", a.GetImportConfig),
		rest.Patch("/imports/:name/config", a.PatchImportConfig),
		rest.Put("/imports/:name/params/:param", a.PutImportParams),

		// Handle exports
		rest.Post("/exports", a.AddExport),
		rest.Get("/exports", a.GetExports),
		rest.Delete("/exports", a.DeleteExport),
		rest.Put("/exports/force", a.ForceExport),
		rest.Put("/exports/stop", a.StopExport),
		rest.Get("/exports/:name/config", a.GetExportConfig),
		rest.Patch("/exports/:name/config", a.PatchExportConfig),
		rest.Put("/exports/:name/params/:param", a.PutExportParams),

		// Handle collections
		rest.Post("/collections", a.PostCollection),
		rest.Get("/collections", a.GetCollections),
		rest.Get("/collections/:name", a.GetCollectionByName),
		rest.Patch("/collections/:name", a.PatchCollection),
		rest.Delete("/collections/:name", a.DeleteCollectionByName),
		rest.Get("/collections/:name/references", a.GetCollectionReferences),
		rest.Get("/collections/:name/config", a.GetCollectionConfig),
		rest.Patch("/collections/:name/config", a.PatchCollectionConfig),
		rest.Put("/collections/:name/config/:path/:param", a.PutCollectionConfigParam),

		// Handle collection buffer
		rest.Get("/collections/:name/buffers", a.GetCollectionBuffers),
		rest.Delete("/collections/:name/buffers", a.DeleteCollectionBuffers),
		rest.Get("/collections/:name/buffers/:id", a.GetCollectionSingleBuffer),
		rest.Patch("/collections/:name/buffers/:id", a.PatchCollectionSingleBuffer),
		rest.Delete("/collections/:name/buffers/:id", a.DeleteCollectionSingleBuffer),
		rest.Post("/collections/:name/buffers/:id/validate", a.ValidateCollectionSingleBuffer),

		// Handle collection items
		rest.Get("/collections/:name/items", a.GetCollectionItems),
		rest.Delete("/collections/:name/items", a.DeleteCollectionItems),
		rest.Get("/collections/:name/items/:id", a.GetCollectionSingleItem),
		rest.Patch("/collections/:name/items/:id", a.PatchCollectionSingleItem),
		rest.Delete("/collections/:name/items/:id", a.DeleteCollectionSingleItem),

		// Establish connection to the web-services
		streamRoute,
	)

	if err != nil {
		return nil, err
	}

	api.SetApp(router)

	return api, nil
}

func (s *Server) Start() error {
	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Println("Serving at " + s.url)
	return http.Serve(s.stoppable, s.api.MakeHandler())
}

// ServerStop stop web server
func (s *Server) Stop() {
	log.Println("Stop server")
	s.onStop <- false
	s.stoppable.Stop()
}

// SendEvent add new event on the event channel
func (s *Server) SendEvent(event interface{}) {
	fmt.Printf("SEND EVENT %+v\n", event)

	s.events.Add(event)

	// If event streamer is not emptying fifo
	if s.eventStatus == false {

		// Kick event streamer to send this new event
		s.onStop <- true
		s.eventStatus = true
	}
}

var idx int = 0

func (s *Server) HandleStream(w rest.ResponseWriter, r *rest.Request) {
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
		case ok := <-s.onStop:

			if ok == false {
				continue
			}

			for {
				eventJson, err := json.Marshal(s.events.Next())
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

			s.onStop <- false
			s.eventStatus = false
			continue
		}
	}
}
