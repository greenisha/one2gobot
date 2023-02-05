package model

type Station struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	C       string `json:"c"`
	Slug    string `json:"slug"`
	V       string `json:"v"`
}

type CallbackData struct {
	Slug string
}
