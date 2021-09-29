package api

import mine "minesweeper_api"

type addReq struct {
	Line   int `json:"num_line"`
	Column int `json:"num_column"`
	Mines  int `json:"num_mine"`
}

func (req addReq) validate() error {
	if req.Line == 0 || req.Column == 0 || req.Mines == 0 {
		return mine.ErrMalformedEntity
	}
	return nil
}

type clickReq struct {
	id string
	Line   int `json:"line"`
	Column int `json:"column"`
	Flag   bool `json:"flag"`
}

func (req clickReq) validate() error {
	if req.id == "" || req.Line == 0 || req.Column == 0 {
		return mine.ErrMalformedEntity
	}
	return nil
}
