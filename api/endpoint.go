package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"minesweeper_api"
)

func addGameEnpoint(svc minesweeper_api.MineSweepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		id, err := svc.Start(req.Column, req.Line, req.Mines)
		if err != nil {
			return nil, err
		}
		return gameRes{
			ID:      id,
			created: true,
		}, nil
	}
}

func clickCellEnpoint(svc minesweeper_api.MineSweepService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(clickReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		cell, err := svc.Click(req.Line, req.Column, req.Flag, req.id)
		if err != nil {
			res := svc.Stop(req.id)
			return res, nil
		}
		return clickRes{
			State:  cell.State,
			Bomb:   cell.Bomb,
			Line:   cell.Line,
			Column: cell.Column,
			Value:  cell.Value,
		}, nil
	}
}