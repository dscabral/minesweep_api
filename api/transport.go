package api

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	mine "minesweeper_api"
	"net/http"
	"strings"
)

func MakeHandler(svcName string, svc mine.MineSweepService) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r := bone.New()

	r.Post("/games", kithttp.NewServer(
		addGameEnpoint(svc),
		decodeAddRequest,
		EncodeResponse,
		opts...,
	))
	r.Post("/games/:id/click", kithttp.NewServer(
		clickCellEnpoint(svc),
		decodeClickRequest,
		EncodeResponse,
		opts...,
	))

	return r

}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, mine.ErrUnsupportedContentType
	}
	req := addReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, mine.ErrMalformedEntity
	}
	return req, nil
}

func decodeClickRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, mine.ErrUnsupportedContentType
	}
	req := clickReq{
		id:     bone.GetValue(r, "id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, mine.ErrMalformedEntity
	}
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	switch err {
	case mine.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case mine.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case mine.ErrUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}