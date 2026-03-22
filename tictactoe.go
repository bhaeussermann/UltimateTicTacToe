package main

import (
	"fmt"
	"os"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/alphabeta"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/montecarlo"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/keyboard"
	"github.com/gen2brain/beeep"
	"golang.org/x/term"
)

func main() {
  fmt.Println("═══ Ultimate Tic Tac Toe ═══")

  var startPlayer game.Player = game.Cell_X
  var aiDifficulty ai.Difficulty = ai.Difficulty_Easy

  for true {
    didSelect := getGameOptions(&startPlayer, &aiDifficulty)
    if !didSelect {
      return
    }

    playerX, playerO := getPlayers(startPlayer, aiDifficulty)

    printInstructions()
    state := game.CreateState()
    undoStates := []*game.State { }
    redoStates := []*game.State { }

    var gameContinuation gameContinuation
    for gameContinuation = gameContinuation_Continue; gameContinuation == gameContinuation_Continue; {
      gameContinuation = step(&state, startPlayer, &playerX, &playerO, &undoStates, &redoStates)
    }
    if gameContinuation == gameContinuation_Stop {
      return;
    }
  }
}

type gameContinuation byte

const (
  gameContinuation_Stop = iota
  gameContinuation_Continue
  gameContinuation_Restart
)

func step(
  state **game.State,
  playerSelection game.Player,
  playerX *player.Player,
  playerO *player.Player,
  undoStates *[]*game.State,
  redoStates *[]*game.State) gameContinuation {
  isGameDone, _ := (*state).GetWinState()
  if !isGameDone {
    fmt.Println()
    fmt.Println((*state).GetSuperBoard().GetHorizontalLine())
    fmt.Println()
    fmt.Println((*state).GetSuperBoard().ToString((*state).GetActiveBoard()))
  }

  var currentPlayer player.Player
  if (*state).GetCurrentPlayer() == game.Cell_X {
    currentPlayer = *playerX
  } else {
    currentPlayer = *playerO
  }
  action, move := currentPlayer.GetMove(*state)
  switch action {
  case player.Action_Restart:
    return gameContinuation_Restart
  case player.Action_Terminate:
    return gameContinuation_Stop
  case player.Action_Undo:
    if len(*undoStates) == 0 {
      beep()
    } else {
      *redoStates = append(*redoStates, *state)
      *state = (*undoStates)[len(*undoStates) - 1]
      *undoStates = (*undoStates)[:len(*undoStates) - 1]
    }
  case player.Action_Redo:
    if len(*redoStates) == 0 {
      beep()
    } else {
      *undoStates = append(*undoStates, *state)
      *state = (*redoStates)[len(*redoStates) - 1]
      *redoStates = (*redoStates)[:len(*redoStates) - 1]
    }
  case player.Action_Move:
    isKeyboardPlayer := (*state).GetCurrentPlayer() == playerSelection
    if isKeyboardPlayer {
      *undoStates = append(*undoStates, (*state).Copy())
      *redoStates = (*redoStates)[0:0]
    }
    (*state).Place(move)
  case player.Action_None:
    (*state).CycleCurrentPlayer()
  }

  if action != player.Action_None {
    var winner game.Player
    isGameDone, winner = (*state).GetWinState()
    if isGameDone {
      fmt.Println()
      fmt.Println((*state).GetSuperBoard().ToString(nil))
      fmt.Println()
      switch (winner) {
      case game.Cell_X: fmt.Println("Cross is the winner!")
      case game.Cell_O: fmt.Println("Naughts is the winner!")
      default: fmt.Println("It's a tie.")
      }
      fmt.Println()
    }
  }

  return gameContinuation_Continue
}

func printInstructions() {
  fmt.Println()
  fmt.Println("Type one of the following to place at the corresponding position:")
  fmt.Println()
  fmt.Println(" 7 │ 8 │ 9")
  fmt.Println("───┼───┼──")
  fmt.Println(" 4 │ 5 │ 6")
  fmt.Println("───┼───┼──")
  fmt.Println(" 1 │ 2 │ 3")
  fmt.Println()
  fmt.Println("• Esc to quit")
  fmt.Println("• 'R' to reset")
}

func getGameOptions(playerSelection *game.Player, aiDifficulty *ai.Difficulty) bool {
  fmt.Println()
  printGameSelection(*playerSelection, *aiDifficulty)
  fmt.Println("ENTER to start. Esc to quit.")

  for true {
    key, error := readKey()
    if error != nil {
      fmt.Println(error)
      return false
    }

    if key == '1' {
      if *playerSelection == game.Cell_X { *playerSelection = game.Cell_O } else { *playerSelection = game.Cell_X }
    } else if key == '2' {
      if *aiDifficulty == ai.Difficulty_Easy {
        *aiDifficulty = ai.Difficulty_Medium
      } else if *aiDifficulty == ai.Difficulty_Medium {
        *aiDifficulty = ai.Difficulty_Hard
      } else {
        *aiDifficulty = ai.Difficulty_Easy
      }
    } else if key == 27 { // Escape
      return false
    } else if key == 13 { // Enter
      return true
    } else {
      beep()
    }

    fmt.Println()
    printGameSelection(*playerSelection, *aiDifficulty)
  }
  return false
}

func printGameSelection(selectedPlayer game.Player, aiDifficulty ai.Difficulty) {
  fmt.Print("(1) Player selection: ")
  if selectedPlayer == game.Cell_X { fmt.Println("X") } else { fmt.Println("O") }

  fmt.Print("(2) AI difficulty: ")
  if aiDifficulty == ai.Difficulty_Easy {
    fmt.Println("Easy")
  } else if aiDifficulty == ai.Difficulty_Medium {
    fmt.Println("Medium")
  } else {
    fmt.Println("Hard")
  }
}

func getPlayers(startPlayer game.Player, aiDifficulty ai.Difficulty) (player.Player, player.Player) {
	humanPlayer := &keyboard.Player{}
	var aiPlayer player.Player
	if aiDifficulty == ai.Difficulty_Easy {
		aiPlayer = &alphabeta.Player{Difficulty: ai.Difficulty_Easy}
  } else if aiDifficulty == ai.Difficulty_Medium {
		aiPlayer = &alphabeta.Player{Difficulty: ai.Difficulty_Hard}
  } else {
		aiPlayer = &montecarlo.Player{Difficulty: ai.Difficulty_Hard}
	}
  
	var playerX, playerO player.Player
	if startPlayer == game.Cell_X {
		playerX = humanPlayer
		playerO = aiPlayer
	} else {
		playerX = aiPlayer
		playerO = humanPlayer
	}
	return playerX, playerO
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
