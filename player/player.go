package player

import (
	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type Player interface {
  GetMove(state *game.State) (*game.Move, bool)
}
