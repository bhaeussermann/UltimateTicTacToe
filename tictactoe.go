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
  fmt.Println("═══ Super Tic Tac Toe ═══")

  for true {
    fmt.Println()
    fmt.Println("Would you like to play as X or O?")
    playerSelection, didSelect := getPlayerSelection()
    if !didSelect {
      return
    }
    var playerX, playerO player.Player
    if playerSelection == game.Cell_X {
      playerX = &player.Keyboard{}
      playerO = &player.AI{}
    } else {
      playerX = &player.AI{}
      playerO = &player.Keyboard{}
    }

    printInstructions()
    state := game.CreateState()
    undoStates := []*game.State { }
    redoStates := []*game.State { }

    var gameContinuation gameContinuation
    for gameContinuation = gameContinuation_Continue; gameContinuation == gameContinuation_Continue; {
      gameContinuation = step(&state, playerSelection, &playerX, &playerO, &undoStates, &redoStates)
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
      beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    } else {
      *redoStates = append(*redoStates, *state)
      *state = (*undoStates)[len(*undoStates) - 1]
      *undoStates = (*undoStates)[:len(*undoStates) - 1]
    }
  case player.Action_Redo: {
    if len(*redoStates) == 0 {
      beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    } else {
      *undoStates = append(*undoStates, *state)
      *state = (*redoStates)[len(*redoStates) - 1]
      *redoStates = (*redoStates)[:len(*redoStates) - 1]
    }
  }
  case player.Action_Move:
    isKeyboardPlayer := (*state).GetCurrentPlayer() == playerSelection
    if isKeyboardPlayer {
      *undoStates = append(*undoStates, (*state).Copy())
      *redoStates = (*redoStates)[0:0]
    }
    (*state).Place(move)
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
  fmt.Println("(Esc to quit; 'R' to reset)")
  fmt.Println()
  fmt.Println(" 1 │ 2 │ 3")
  fmt.Println("───┼───┼──")
  fmt.Println(" 4 │ 5 │ 6")
  fmt.Println("───┼───┼──")
  fmt.Println(" 7 │ 8 │ 9")
}

func getPlayerSelection() (game.Player, bool) {
  oldState, error := term.MakeRaw(int(os.Stdin.Fd()))
  if error != nil {
    fmt.Println(error)
    return game.Cell_None, false
  }
  defer term.Restore(int(os.Stdin.Fd()), oldState)

  readBuffer := make([]byte, 1)
  for true {
    os.Stdin.Read(readBuffer)
    switch readBuffer[0] {
    case 27: return game.Cell_None, false
    case byte('x'), byte('X'): return game.Cell_X, true
    case byte('o'), byte('O'): return game.Cell_O, true
    default: beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
    }
  }
  return game.Cell_None, false
}
