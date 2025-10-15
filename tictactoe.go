package main

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/gen2brain/beeep"
)

func main() {
  printInstructions()
  
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    os.Exit(1)
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)
  readBuffer := make([]byte, 1)

  state := game.CreateState()
  var done bool
  var winner game.Player
  for ; !done; done, winner = state.GetWinState() {
    fmt.Println(state.GetBoard().ToString() + "\r\n")
    if state.GetCurrentPlayer() == game.X {
      fmt.Print("Cross' turn to move: ")
    } else {
      fmt.Print("Naught's turn to move: ")
    }

    for didMove := false; !didMove; {
      os.Stdin.Read(readBuffer)
      if readBuffer[0] == 27 { // Escape
        return
      }
      if (readBuffer[0] >= byte('1')) && (readBuffer[0] <= byte('9')) {
        blockNumber := byte(readBuffer[0]) - byte('1')
        rowNumber := blockNumber / 3
        columnNumber := blockNumber % 3
        if (state.Place(rowNumber, columnNumber)) {
          fmt.Println(blockNumber + 1)
          didMove = true
        } else {
          beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
        }
      } else {
        beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
      }
    }
    fmt.Println("\r\n")
  }

  fmt.Println(state.GetBoard().ToString() + "\r\n")

  switch (winner) {
  case game.X: fmt.Println("Cross is the winner!")
  case game.O: fmt.Println("Naughts is the winner!")
  default: fmt.Println("It's a tie.")
  }
}

func printInstructions() {
  fmt.Println("=== Super Tic Tac Toe ===")
  fmt.Println()
  fmt.Println("Type one of the following to place at the corresponding position (Esc to quit):")
  fmt.Println()
  fmt.Println(" 1 | 2 | 3")
  fmt.Println("----------")
  fmt.Println(" 4 | 5 | 6")
  fmt.Println("----------")
  fmt.Println(" 7 | 8 | 9")
  fmt.Println()
}
