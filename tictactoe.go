package main

import (
	"fmt"
	"os"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/gen2brain/beeep"
	"golang.org/x/term"
)

func main() {
  fmt.Println("=== Super Tic Tac Toe ===")
  fmt.Println()
  fmt.Println("Would you like to play as X or O?")

  playerSelection, didSelect := getPlayerSelection()
  if !didSelect {
    return
  }
  var playerX, playerO player.Player
  if playerSelection == game.X {
    playerX = &player.Keyboard {}
    playerO = &player.AI {}
  } else {
    playerX = &player.AI {}
    playerO = &player.Keyboard {}
  }

  printInstructions()
  state := game.CreateState()
  var done bool
  var winner game.Player
  for ; !done; done, winner = state.GetWinState() {
    fmt.Println()
    fmt.Println(state.GetBoard().ToString())

    var currentPlayer player.Player
    if state.GetCurrentPlayer() == game.X {
      currentPlayer = playerX
    } else {
      currentPlayer = playerO
    }
    move, shouldContinue := currentPlayer.GetMove(state)
    if !shouldContinue {
      return
    }
    state.Place(move)
  }

  fmt.Println()
  fmt.Println(state.GetBoard().ToString())
  fmt.Println()

  switch (winner) {
  case game.X: fmt.Println("Cross is the winner!")
  case game.O: fmt.Println("Naughts is the winner!")
  default: fmt.Println("It's a tie.")
  }
}

func printInstructions() {
  fmt.Println("Type one of the following to place at the corresponding position (Esc to quit):")
  fmt.Println()
  fmt.Println(" 1 | 2 | 3")
  fmt.Println("----------")
  fmt.Println(" 4 | 5 | 6")
  fmt.Println("----------")
  fmt.Println(" 7 | 8 | 9")
}

func getPlayerSelection() (game.Player, bool) {
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    return game.None, false
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)

  readBuffer := make([]byte, 1)
  for true {
    os.Stdin.Read(readBuffer)
    switch readBuffer[0] {
    case 27: return game.None, false
    case byte('x'), byte('X'): return game.X, true
    case byte('o'), byte('O'): return game.O, true
    default: beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    }
  }
  return game.None, false
}
