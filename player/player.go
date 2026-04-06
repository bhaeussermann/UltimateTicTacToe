package player

import (
	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type Player interface {
  GetMove(state *game.State, log Log) (Action, *game.Move)
}

type Action byte

const (
  Action_None = iota
  Action_Restart
  Action_Terminate
  Action_Undo
  Action_Redo
  Action_Move
)
