package player

import (
	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type Player interface {
  GetMove(state *game.State) (Action, *game.Move)
}

type Action byte

const (
  Action_Terminate = iota
  Action_Undo
  Action_Redo
  Action_Move
)
