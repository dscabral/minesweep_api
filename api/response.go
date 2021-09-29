package api

import "net/http"

type gameRes struct {
	ID      string `json:"id"`
	created bool
}

func (s gameRes) Code() int {
	if s.created {
		return http.StatusCreated
	}

	return http.StatusOK
}

func (s gameRes) Headers() map[string]string {
	return map[string]string{}
}

func (s gameRes) Empty() bool {
	return false
}

type clickRes struct {
	State  string `json:"state"`
	Bomb   bool   `json:"bomb"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Value  int    `json:"value"`
}

func (s clickRes) Code() int {
	return http.StatusOK
}

func (s clickRes) Headers() map[string]string {
	return map[string]string{}
}

func (s clickRes) Empty() bool {
	return false
}
