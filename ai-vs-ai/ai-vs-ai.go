package main

import (
	"fmt"
	"time"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/montecarlo"
	"github.com/inancgumus/screen"
)

func main() {
  go play()
  fmt.Scanln()
}

func play() {
  playerX := montecarlo.Player{ Difficulty: ai.Difficulty_Easy }
  playerO := montecarlo.Player{ Difficulty: ai.Difficulty_Easy }
  playerXLog := player.CreateLog()
  playerOLog := player.CreateLog()
  playerXScore := 0
  playerOScore := 0
  for true {
    state := game.CreateState()
    var isDone bool
    for isDone, _ = state.GetWinState(); !isDone; isDone, _ = state.GetWinState() {
      printState(state, playerXLog, playerOLog, playerXScore, playerOScore)
      var currentPlayer player.Player
      var currentPlayerLog *player.MessageLog
      if state.GetCurrentPlayer() == game.Cell_X {
        currentPlayer = &playerX
        currentPlayerLog = playerXLog
      } else {
        currentPlayer = &playerO
        currentPlayerLog = playerOLog
      }
      currentPlayerLog.Clear()
      _, move := currentPlayer.GetMove(state, currentPlayerLog)
      state.Place(move)
    }

    _, winner := state.GetWinState()
    switch winner {
    case game.Cell_X: playerXScore++
    case game.Cell_O: playerOScore++
    default:
      playerXScore++
      playerOScore++
    }
    printState(state, playerXLog, playerOLog, playerXScore, playerOScore)
    time.Sleep(time.Second)
  }
}

func printState(state *game.State, playerXLog *player.MessageLog, playerOLog *player.MessageLog, playerXScore int, playerOScore int) {
	clear()
	fmt.Printf("X\t%v\r\nO\t%v\r\n\r\n", playerXScore, playerOScore)
	fmt.Println(state.GetSuperBoard().ToString((*state).GetActiveBoard()))
	fmt.Println("Player X log:")
	for _, message := range playerXLog.GetMessages() {
		fmt.Printf("\t%s", message)
	}
  fmt.Println()
	fmt.Println("Player O log:")
	for _, message := range playerOLog.GetMessages() {
		fmt.Printf("\t%s", message)
	}
  fmt.Println()
  fmt.Println("Press ENTER to stop.")
}

func clear() {
  screen.Clear()
  screen.MoveTopLeft()
}
