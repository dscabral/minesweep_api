package minesweeper_api

import (
	"go.uber.org/zap"
	"sync"
)

var _ MineSweepService = (*mineSweepService)(nil)

type mineSweepService struct {
	logger *zap.Logger
	board map[string]Board
	mu sync.Mutex
}

func NewService(logger *zap.Logger) MineSweepService {
	return &mineSweepService{
		logger: logger,
		board: make(map[string]Board),
	}
}

