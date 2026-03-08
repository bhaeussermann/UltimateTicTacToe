package ai

import (
	"sync"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type Difficulty byte

const (
  Difficulty_Easy = iota
  Difficulty_Medium
  Difficulty_Hard
)

type Location struct {
  RowNumber byte
  ColumnNumber byte
}

func getSideNumbers() []byte {
  sideNumbers := make([]byte, game.Size)
  var number byte
  for number = 0; number < game.Size; number++ {
    sideNumbers[number] = number
  }
  return sideNumbers
}

var SideNumbers = sync.OnceValue(getSideNumbers)
