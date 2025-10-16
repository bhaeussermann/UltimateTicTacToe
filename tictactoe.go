package main

import (
	"fmt"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
)

func main() {
  printInstructions()

  state := game.CreateState()
  playerX := &player.Keyboard {}
  playerO := &player.Keyboard {}

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
}
