package player

import (
	"fmt"
	"os"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/gen2brain/beeep"
	"golang.org/x/term"
)

type Keyboard struct {}

func (*Keyboard) GetMove(state *game.State) (Action, *game.Move) {
  fmt.Println()
  if state.GetCurrentPlayer() == game.Cell_X {
    fmt.Print("Cross' turn to move: ")
  } else {
    fmt.Print("Naught's turn to move: ")
  }

  boardReference := state.GetActiveBoard()
  if boardReference == nil {
    fmt.Print("\r\nSelect the board to play on: ")
    action, cellReference := getCellReference(
      func(c *cellReference) bool { return state.CanPlaceIn(&game.BoardReference { RowNumber: c.rowNumber, ColumnNumber: c.columnNumber })})
    if action != Action_Move {
      return action, nil
    }
    boardReference = &game.BoardReference { RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }
    fmt.Println(state.GetSuperBoard().ToString(boardReference))
    fmt.Print("Your move: ")
  }

  action, cellReference := getCellReference(
    func(c *cellReference) bool { return state.CanPlace(&game.Move { Board: boardReference, RowNumber: c.rowNumber, ColumnNumber: c.columnNumber }) })
  if action == Action_Move {
    return action, &game.Move { Board: boardReference, RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }
  }
  return action, nil
}

func getCellReference(canPlace func(*cellReference) bool) (Action, *cellReference) {  
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    return Action_Terminate, nil
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)
  
  readBuffer := make([]byte, 1)
  for true {
    os.Stdin.Read(readBuffer)
    if readBuffer[0] == 25 { // Ctrl + Y
      return Action_Redo, nil
    }
    if readBuffer[0] == 26 { // Ctrl + Z
      return Action_Undo, nil
    }
    if readBuffer[0] == 27 { // Escape
      return Action_Terminate, nil
    }
    
    if (readBuffer[0] >= byte('1')) && (readBuffer[0] <= byte('9')) {
      blockNumber := byte(readBuffer[0]) - byte('1')
      cellReference := cellReference { rowNumber: blockNumber / 3, columnNumber: blockNumber % 3 }
      if canPlace(&cellReference) {
        fmt.Print(blockNumber + 1)
        fmt.Print("\r\n")
        return Action_Move, &cellReference
      } else {
        beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
      }
    } else {
      beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    }
  }
  return Action_Terminate, nil
}

type cellReference struct {
  rowNumber byte
  columnNumber byte
}
