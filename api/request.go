package api

import mine "minesweeper_api"

type addReq struct {
	Line   string `json:"num_line"`
	Column string `json:"num_column"`
	Mines  string `json:"num_mine"`
}

func (req addReq) validate() error {
	if req.Line == "" || req.Column == "" || req.Mines == "" {
		return mine.ErrMalformedEntity
	}
	return nil
}

type clickReq struct {
	Line   string `json:"line"`
	Column string `json:"column"`
	Flag   string `json:"flag"`
}

func (req clickReq) validate() error {
	if req.Line == "" || req.Column == "" {
		return mine.ErrMalformedEntity
	}
	return nil
}
