package minesweeper_api

import "errors"

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid author or content).
	ErrMalformedEntity = errors.New("malformed entity specification")
	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")
	// ErrUnsupportedContentType indicates a unsupported content-type (should be application/json)
	ErrUnsupportedContentType = errors.New("unsupported content-type")
)

type Board struct {
	Id    string
	Cells map[string]Cell
}

type Cell struct {
	State  string
	Bomb   bool
	Line   int
	Column int
	Value  int
}

type MineSweepService interface {
	Start(numCols int, numLines int, numBombs int) (string, error)
	Stop(boardID string) Board
	Click(line int, column int, flag bool, boardID string) (Cell, error)
}
