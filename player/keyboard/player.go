package keyboard

import (
	"fmt"
	"os"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/gen2brain/beeep"
	"golang.org/x/term"
)

type Player struct {}

func (*Player) GetMove(state *game.State) (player.Action, *game.Move) {
  done, _ := state.GetWinState()

  if done {
    fmt.Print("Play again? (Y / N) ")
    for true {
      key, error := readKey()
      if error != nil {
        fmt.Println(error)
        return player.Action_Terminate, nil
      }
      if (key == byte('y')) || (key == byte('Y')) || (key == byte('r') || (key == byte('R'))) {
        return player.Action_Restart, nil
      }
      if (key == byte('n')) || (key == byte('N')) || (key == 27) /* Escape */ {
        return player.Action_Terminate, nil
      }
      if key == 26 { // Ctrl + Z
        return player.Action_Undo, nil
      }
      beep()
    }
  }

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
    if action != player.Action_Move {
      return action, nil
    }
    boardReference = &game.BoardReference { RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }
    fmt.Println(state.GetSuperBoard().ToString(boardReference))
    fmt.Print("Your move: ")
  }

  action, cellReference := getCellReference(
    func(c *cellReference) bool { return state.CanPlace(&game.Move { Board: boardReference, RowNumber: c.rowNumber, ColumnNumber: c.columnNumber }) })
  if action == player.Action_Move {
    return action, &game.Move { Board: boardReference, RowNumber: cellReference.rowNumber, ColumnNumber: cellReference.columnNumber }
  }
  return action, nil
}

func getCellReference(canPlace func(*cellReference) bool) (player.Action, *cellReference) {  
  for true {
    key, error := readKey()
    if error != nil {
      fmt.Println(error)
      return player.Action_Terminate, nil
    }

    if key == 25 { // Ctrl + Y
      return player.Action_Redo, nil
    }
    if key == 26 { // Ctrl + Z
      return player.Action_Undo, nil
    }
    if key == 27 { // Escape
      return player.Action_Terminate, nil
    }
    if (key == byte('r') || (key == byte('R'))) {
      return player.Action_Restart, nil
    }
    
    if (key >= byte('1')) && (key <= byte('9')) {
      blockNumber := byte(key) - byte('1')
      cellReference := cellReference { rowNumber: 2 - blockNumber / 3, columnNumber: blockNumber % 3 }
      if canPlace(&cellReference) {
        fmt.Print(blockNumber + 1)
        fmt.Print("\r\n")
        return player.Action_Move, &cellReference
      } else {
        beep()
      }
    } else {
      beep()
    }
  }
  return player.Action_Terminate, nil
}

func readKey() (byte, error) {
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    return 0, error
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)
  readBuffer := make([]byte, 1)
  os.Stdin.Read(readBuffer)
  return readBuffer[0], nil
}

func beep() {
  beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
}

type cellReference struct {
  rowNumber byte
  columnNumber byte
}
