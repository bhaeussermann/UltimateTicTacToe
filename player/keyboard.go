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
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    return nil, false
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)
  
  fmt.Println()
  if state.GetCurrentPlayer() == game.X {
    fmt.Print("Cross' turn to move: ")
  } else {
    fmt.Print("Naught's turn to move: ")
  }
    
  readBuffer := make([]byte, 1)
  for true {
    os.Stdin.Read(readBuffer)
    if readBuffer[0] == 27 { // Escape
      return nil, false
    }
    if (readBuffer[0] >= byte('1')) && (readBuffer[0] <= byte('9')) {
      blockNumber := byte(readBuffer[0]) - byte('1')
      rowNumber := blockNumber / 3
      columnNumber := blockNumber % 3
      move := &game.Move { RowNumber: rowNumber, ColumnNumber: columnNumber }
      if (state.CanPlace(move)) {
        fmt.Println(blockNumber + 1)
        return move, true
      } else {
        beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
      }
    } else {
      beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    }
  }
  return nil, false
}
