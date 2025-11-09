package player

import (
	"fmt"
	"os"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/gen2brain/beeep"
	"golang.org/x/term"
)

type Keyboard struct {}

func (*Keyboard) GetMove(state *game.State) (*game.Move, bool) {
  fmt.Println()
  if state.GetCurrentPlayer() == game.X {
    fmt.Print("Cross' turn to move: ")
  } else {
    fmt.Print("Naught's turn to move: ")
  }

  boardReference := state.GetActiveBoard()
  if boardReference == nil {
    fmt.Print("\r\nSelect the board to play on: ")
    cellReference, shouldContinue := getCellReference(
      func(c *cellReference) bool { return state.CanPlaceIn(&game.BoardReference { RowNumber: c.rowNumber, ColumnNumber: c.columnNumber })})
    if !shouldContinue {
      return nil, false
    }
    boardReference = &game.BoardReference { RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }
    fmt.Println(state.GetBoard().ToString(boardReference))
    fmt.Print("Your move: ")
  }

  cellReference, shouldContinue := getCellReference(
    func(c *cellReference) bool { return state.CanPlace(&game.Move { Board: boardReference, RowNumber: c.rowNumber, ColumnNumber: c.columnNumber }) })
  if !shouldContinue {
    return nil, false
  }
  return &game.Move { Board: boardReference, RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }, true
}

func getCellReference(canPlace func(*cellReference) bool) (*cellReference, bool) {  
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    return nil, false
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)
  
  readBuffer := make([]byte, 1)
  for true {
    os.Stdin.Read(readBuffer)
    if readBuffer[0] == 27 { // Escape
      return nil, false
    }
    if (readBuffer[0] >= byte('1')) && (readBuffer[0] <= byte('9')) {
      blockNumber := byte(readBuffer[0]) - byte('1')
      cellReference := cellReference { rowNumber: blockNumber / 3, columnNumber: blockNumber % 3 }
      if canPlace(&cellReference) {
        fmt.Print(blockNumber + 1)
        fmt.Print("\r\n")
        return &cellReference, true
      } else {
        beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
      }
    } else {
      beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    }
  }
  return nil, false
}

type cellReference struct {
  rowNumber byte
  columnNumber byte
}
