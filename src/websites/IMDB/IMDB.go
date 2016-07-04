package IMDB

import (
	"github.com/ohohleo/classify/requests"
	"github.com/ohohleo/classify/websites"
	"strings"
)

type IMDB struct {
	Url string
}

func New() *IMDB {
	return &IMDB{
		Url: "http://imdb.wemakesites.net/api/",
	}
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

// Launch a search request through the IMDB Api
func (i *IMDB) search(input string) chan *Results {

	c := make(chan *Results)

	go func() {
		var rsp Response

		// Replace space by '+'
		input = strings.Replace(input, " ", "+", -1)

		// Send the request
		channel, err := requests.Send("GET", i.Url+"search?q="+input, nil, nil, &rsp)
		if err != nil {
			close(c)
		}

		_, ok := <-channel
		if ok == false {
			close(c)
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
		channel, err := requests.Send("GET", i.Url+id, nil, nil, &rsp)
		if err != nil {
			close(c)
		}

		// Wait for the result
		_, ok := <-channel
		if ok == false {
			close(c)
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
		}

		// Search all the movies
		for _, title := range results.Titles {
			data, ok := <-i.getResource(title.Id)
			if ok {

				// TODO: duration & released

				c <- &websites.Movie{
					Name:        title.Title,
					Url:         title.Url,
					Image:       data.Image,
					Description: data.Description,
					Directors:   data.Directors,
					Cast:        data.Cast,
					Genres:      data.Genres,
				}
			}
		}

		close(c)
	}()

	return c
}
