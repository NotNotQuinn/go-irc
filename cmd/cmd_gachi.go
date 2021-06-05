package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/NotNotQuinn/go-irc/config"
)

type APIresponse struct {
	Statuscode int      `json:"statusCode"`
	Timestamp  int64    `json:"timestamp"`
	Data       Data     `json:"data"`
	Error      APIError `json:"error"`
}
type Data struct {
	Statuscode int       `json:"statuscode"`
	Timestamp  int64     `json:"timestamp"`
	Data       InnerData `json:"data"`
	Error      APIError  `json:"error"`
}
type InnerData struct {
	ID            int             `json:"ID"`
	Link          string          `json:"link"`
	Name          string          `json:"name"`
	Videotype     int             `json:"videotype"`
	Tracktype     string          `json:"tracktype"`
	Duration      float64         `json:"duration"`
	Available     bool            `json:"available"`
	Published     time.Time       `json:"published"`
	Notes         string          `json:"notes"`
	Addedby       string          `json:"addedby"`
	Addedon       time.Time       `json:"addedon"`
	Lastedit      string          `json:"lastedit"`
	Parsedlink    string          `json:"parsedlink"`
	Aliases       []string        `json:"aliases"`
	Authors       []Author        `json:"authors"`
	Tags          []string        `json:"tags"`
	Relatedtracks []Relatedtracks `json:"relatedtracks"`
	Legacyid      int             `json:"legacyID"`
	Favourites    int             `json:"favourites"`
}
type Relatedtracks struct {
	Notes        string `json:"notes"`
	Relationship string `json:"relationship"`
	Fromid       int    `json:"fromID"`
	Toid         int    `json:"toID"`
	Name         string `json:"name"`
}

type APIError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Author struct {
	Role string `json:"role"`
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

var (
	gachiCommand *Command = &Command{
		Name: "gachi",
		Execution: func(c *Context) (*Return, error) {
			errReturn := &Return{
				Success: false,
				Reply:   "There was an error in attempting to get a random gachi HandsUp",
			}
			body := bytes.NewBuffer(nil)
			req, err := http.NewRequest("GET", "https://supinic.com/api/track/gachi/random/", body)
			if err != nil {
				return errReturn, err
			}
			req.Header.Set("User-Agent", config.Public.Global.UserAgent)
			client := &http.Client{}
			if err != nil {
				return errReturn, err
			}
			var Response APIresponse
			resp, err := client.Do(req)
			if err != nil {
				return errReturn, err
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errReturn, err
			}
			err = json.Unmarshal(b, &Response)
			if err != nil {
				return errReturn, err
			}
			return &Return{
				Reply:   fmt.Sprintf("gachiGASM %s", Response.Data.Data.Parsedlink),
				Success: true,
			}, nil
		},
		Description: "Fetches a random gachi song from the internet.",
	}
)