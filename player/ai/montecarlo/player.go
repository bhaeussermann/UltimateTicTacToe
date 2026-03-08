package montecarlo

import (
	"math/rand"
	"time"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
)

type Player struct {}

func (*Player) GetMove(state *game.State) (player.Action, *game.Move) {
  done, _ := state.GetWinState()
  if done {
    return player.Action_None, nil
  }

  me := state.GetCurrentPlayer()
  potentialMoves := getPotentialMoves(state)
  totalGames := 0
  winCounts := make([]int, len(potentialMoves))
  for deadline := time.Now().Add(time.Second * 2); time.Now().Before(deadline); {
    for moveIndex, move := range potentialMoves {
      stateCopy := state.Copy()
      stateCopy.Place(&move)
      winner := play(stateCopy)
      if (winner == me) || ((winner == game.Cell_None) && (rand.Intn(2) == 0)) {
        winCounts[moveIndex]++
      }
      totalGames++
    }
  }

  maximumWins := -1
  var bestMove game.Move
  for moveIndex, move := range potentialMoves {
    if winCounts[moveIndex] > maximumWins {
      maximumWins = winCounts[moveIndex]
      bestMove = move
    }
  }
  return player.Action_Move, &bestMove
}

func play(state *game.State) game.Player {
  done, winner := state.GetWinState()
  for ; !done; done, winner = state.GetWinState() {
    potentialMoves := getPotentialMoves(state)
    move := potentialMoves[rand.Intn(len(potentialMoves))]
    state.Place(&move)
  }
  return winner
}

func getPotentialMoves(state *game.State) []game.Move {
	var potentialMoves []game.Move
  activeBoardReference := state.GetActiveBoard()
	if activeBoardReference != nil {
		for _, location := range getPotentialMoveLocations(state.GetBoard(activeBoardReference).Cells) {
			potentialMoves = append(potentialMoves, game.Move{Board: activeBoardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
		}
	} else {
    for _, potentialMoveBoard := range getPotentialMoveLocations(state.GetSuperBoard()) {
      boardReference := game.BoardReference{RowNumber: potentialMoveBoard.RowNumber, ColumnNumber: potentialMoveBoard.ColumnNumber}
      board := state.GetBoard(&boardReference)
      for _, location := range getPotentialMoveLocations(board.Cells) {
        potentialMoves = append(potentialMoves, game.Move{Board: &boardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
      }
		}
	}

  return potentialMoves
}

func getPotentialMoveLocations(cellGrid game.CellGrid) []ai.Location {
  var potentialMoveLocations []ai.Location

  for _, rowNumber := range ai.SideNumbers() {
    for _, columnNumber := range ai.SideNumbers() {
      if cellGrid.IsEmpty(rowNumber, columnNumber) {
        location := ai.Location{RowNumber: rowNumber, ColumnNumber: columnNumber}
        potentialMoveLocations = append(potentialMoveLocations, location)
      }
    }
  }
  return potentialMoveLocations
}
