package IMDB

import (
	"fmt"
	"github.com/ohohleo/classify/collections"
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"strings"
)

type IMDB struct {
	Url    string
	apiKey string
}

func New() *IMDB {
	return &IMDB{
		Url: "http://imdb.wemakesites.net/api/",
	}
}

func (i *IMDB) SetConfig(config map[string]string) bool {

	fmt.Printf("SetConfig %+v\n", config)

	// Get API key
	apiKey, ok := config["api_key"]
	if ok {
		i.apiKey = apiKey
	}

	// Get alternative url
	url, ok := config["url"]
	if ok {
		i.Url = url
	}

	return true
}

type Response struct {
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Term      string `json:"term,omitempty"`
	SearchUrl string `json:"search_url,omitempty"`
	Data      Data   `json:"data"`
}

type Data struct {
	Id          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	Title       string   `json:"title,omitempty"`
	Year        int      `json:"int,omitempty"`
	Description string   `json:"description,omitempty"`
	Certificate string   `json:"certificate,omitempty"`
	Duration    string   `json:"duration,omitempty"`
	Released    string   `json:"released,omitempty"`
	Filmography []Title  `json:"filmography,omitempty"`
	Cast        []string `json:"cast,omitempty"`
	Genres      []string `json:"genre,omitempty"`
	Occupation  []string `json:"occupation,omitempty"`
	Directors   []string `json:"directors,omitempty"`
	Writers     []string `json:"writers,omitempty"`
	MediaLinks  []string `json:"mediaLinks,omitempty"`
	Image       string   `json:"image,omitempty"`
	Results     Results  `json:"results,omitempty"`
}

type Results struct {
	Titles     []Title
	Names      []Title
	Characters []Title
	Companies  []Title
}

type Title struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Info  int    `json:"info,omitempty"`
	Year  int    `json:"year,omitempty"`
}

func (i *IMDB) GetName() string {
	return "IMDB"
}

// Launch a search request through the IMDB Api
func (i *IMDB) search(input string) chan *Results {

	c := make(chan *Results)

	go func() {
		var rsp Response

		queries := make(map[string]string)

		// Replace space by '+'
		queries["q"] = strings.Replace(input, " ", "+", -1)

		// Send the request
		channel, err := i.send("search", queries, &rsp)
		if err != nil {
			fmt.Printf("Request error: %s\n", err.Error())
			close(c)
			return
		}

		_, ok := <-channel
		if ok == false {
			close(c)
			return
		}

		// Return the results
		data, ok := i.checkResponse(rsp)
		if ok {
			c <- &data.Results
		}

		close(c)
	}()

	return c
}

// Get the data informations with IMDB id
func (i *IMDB) getResource(id string) chan *Data {

	c := make(chan *Data)

	go func() {
		var rsp Response
		channel, err := i.send(id, nil, &rsp)
		if err != nil {
			close(c)
			return
		}

		// Wait for the result
		_, ok := <-channel
		if ok == false {
			close(c)
			return
		}

		// Return the results
		data, ok := i.checkResponse(rsp)
		if ok {
			c <- data
		}

		close(c)
	}()

	return c
}

// Return true if the response is ok
func (i *IMDB) checkResponse(rsp Response) (*Data, bool) {

	if rsp.Code == 200 {
		return &rsp.Data, true
	}

	return nil, false
}

// Generic method used to search matching movies
func (i *IMDB) Search(input string) chan websites.Data {

	c := make(chan websites.Data)

	go func() {

		// Launch the research
		results, ok := <-i.search(input)
		if ok == false {
			close(c)
			return
		}

		// Search all the movies
		for _, title := range results.Titles {
			data, ok := <-i.getResource(title.Id)
			if ok {

				// TODO: duration & released

				movie := &collections.Movie{
					Name:        title.Title,
					Url:         title.Url,
					Image:       data.Image,
					Description: data.Description,
					Directors:   data.Directors,
					Cast:        data.Cast,
					Genres:      data.Genres,
				}

				movie.Init()

				c <- movie
			}
		}

		close(c)
	}()

	return c
}

// Envoie d'une requête
func (i *IMDB) send(id string, queries map[string]string, rsp interface{}) (chan *requests.Response, error) {

	if i.apiKey != "" {

		if queries == nil {
			queries = make(map[string]string)
		}

		queries["api_key"] = i.apiKey
	}

	return requests.Send("GET", i.Url+id, nil, queries, nil, &rsp)
}
