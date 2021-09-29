package minesweeper_api

import (
	"fmt"
	"github.com/gofrs/uuid"
	"math/rand"
	"time"
)

func (svc *mineSweepService) Start(numCols int, numLines int, numBombs int) (string, error) {
	var mapCells = map[string]Cell{}
	var mapRandom = map[int]bool{}
	for r := 0; r < numBombs; r++ {
		rand.Seed(time.Now().UnixNano())
		min := 0
    	max := numLines * numCols
    	mapRandom[rand.Intn(max - min) + min] = true
	}
	var count = 0
	for i := 0 ; i < numCols; i++ {
		for j := 0; j < numLines; j++ {
    		key := key(i+1, j+1)
			if ok := mapRandom[count]; ok {
    			mapCells[key] = Cell{
					State:  "covered",
					Bomb:   true,
					Line:   j,
					Column: i,
				}
			} else {
				mapCells[key] = Cell{
					State:  "covered",
					Bomb:   false,
					Line:   j,
					Column: i,
				}
			}
			count++
		}
	}
	ID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	svc.board[ID.String()] = Board{Id: ID.String(), Cells: mapCells}
	return ID.String(), nil
}

func (svc *mineSweepService) Stop(boardID string) Board {
	for k, v := range svc.board[boardID].Cells {
		v.State = "uncovered"
		svc.board[boardID].Cells[k] = v
	}
	return svc.board[boardID]
}

func (svc *mineSweepService) Click(line int, column int, flag bool, boardID string) (Cell, error) {
	key := key(column, line)
	if v, ok := svc.board[boardID].Cells[key]; ok {
		if flag || v.State == "flagged" {
			v.State = "flagged"
		} else {
			v.State = "uncovered"
			svc.board[boardID].Cells[key] = v
			if v.Bomb {
				return svc.board[boardID].Cells[key], nil
			}
		}
		return svc.board[boardID].Cells[key], nil
	}
	return Cell{}, ErrNotFound
}

func key(col int, line int) string {
	return fmt.Sprintf("%d-%d", col, line)
}