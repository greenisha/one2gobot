package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"greenisha/one2gobot/model"
	"net/http"
	"time"
)

type Client struct {
	RestEndpoint string
}

// look for station
func (c *Client) FindStation(search string) ([]model.Station, error) {
	var out []model.Station
	err := getJson(c.RestEndpoint+"/en/typeahead?ajax=1&short=1&pattern="+search, &out)
	if err != nil {
		return []model.Station{}, err
	}
	return out, nil
}

func (c *Client) FindStationBySlug(slug string) (model.Station, error) {
	var out []model.Station
	err := getJson(c.RestEndpoint+"/en/typeahead?ajax=1&short=1&pattern="+slug, &out)
	if err != nil {
		return model.Station{}, err
	}
	for _, v := range out {
		if v.Slug == slug {
			return v, nil
		}
	}
	return model.Station{}, errors.New("not found")
}

// helper get function
func getJson(url string, target interface{}) error {
	myClient := &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

// helper post function
func postJson(url string, target interface{}, data interface{}) error {
	myClient := &http.Client{Timeout: 10 * time.Second}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r, err := myClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
