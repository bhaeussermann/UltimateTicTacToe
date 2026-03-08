package alphabeta

import (
	"math"
	"slices"
	"sync"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
)

type Player struct {
  Difficulty ai.Difficulty
}

func (p *Player) GetMove(state *game.State) (player.Action, *game.Move) {
  done, _ := state.GetWinState()
  if done {
    return player.Action_None, nil
  }
  move, _ := getBestMove(state, p.getDepth(), math.MinInt, math.MaxInt)
  return player.Action_Move, move
}

func (p *Player) getDepth() int {
  if p.Difficulty == ai.Difficulty_Easy {
    return 1
  } else if p.Difficulty == ai.Difficulty_Medium {
    return 4
  } else {
    return 6
  }
}

func getBestMove(state *game.State, depth int, alpha int, beta int) (*game.Move, int) {
  done, _ := state.GetWinState()
  if depth <= 0 || done {
    return nil, getScore(state)
  }

  potentialMoves := getPotentialMoves(state)
  
  var bestMove *game.Move
  bestMoveScore := math.MinInt
  for _, potentialMove := range potentialMoves {
    nextState := state.Copy()
    nextState.Place(&potentialMove)

    _, opponentScore := getBestMove(nextState, depth - 1, -beta, -alpha)
    score := -opponentScore
    if score > bestMoveScore {
      bestMove = &potentialMove
      bestMoveScore = score
      if bestMoveScore >= beta {
        return bestMove, bestMoveScore
      }
      if bestMoveScore > alpha {
        alpha = bestMoveScore
      }
    }
  }

  return bestMove, bestMoveScore
}

func getPotentialMoves(state *game.State) []game.Move {
	var potentialMoves []game.Move
  activeBoardReference := state.GetActiveBoard()
	if activeBoardReference != nil {
		for _, location := range getPotentialMoveLocations(state.GetBoard(activeBoardReference).Cells, state.GetCurrentPlayer()) {
			potentialMoves = append(potentialMoves, game.Move{Board: activeBoardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
		}
	} else {
    for _, potentialMoveBoard := range getPotentialMoveLocations(state.GetSuperBoard(), state.GetCurrentPlayer()) {
      boardReference := game.BoardReference{RowNumber: potentialMoveBoard.RowNumber, ColumnNumber: potentialMoveBoard.ColumnNumber}
      board := state.GetBoard(&boardReference)
      for _, location := range getPotentialMoveLocations(board.Cells, state.GetCurrentPlayer()) {
        potentialMoves = append(potentialMoves, game.Move{Board: &boardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
      }
		}
	}

  return potentialMoves
}

func getPotentialMoveLocations(cellGrid game.CellGrid, player game.Player) []ai.Location {
  potentialMoveLocations := getMoveLocations(cellGrid, player)

  for _, rowNumber := range ai.SideNumbers() {
    for _, columnNumber := range ai.SideNumbers() {
      if cellGrid.IsEmpty(rowNumber, columnNumber) {
        location := ai.Location{RowNumber: rowNumber, ColumnNumber: columnNumber}
        if !slices.Contains(potentialMoveLocations, location) {
          potentialMoveLocations = append(potentialMoveLocations, location)
        }
      }
    }
  }
  return potentialMoveLocations
}

const ForkLocationScore = 1
const WinLocationScore = 2
const BoardScore = 5
const SuperBoardMultiplier = game.Size * game.Size
const GameScore = BoardScore * SuperBoardMultiplier

func getScore(state *game.State) int {
  currentPlayer := state.GetCurrentPlayer()
  done, winner := state.GetWinState()
  if (done) {
    if (winner == currentPlayer) { return GameScore } else { return -GameScore }
  }

  var opponent game.Player
  if currentPlayer == game.Cell_X { opponent = game.Cell_O } else { opponent = game.Cell_X }

  opponentMoveBoardReference := getMoveBoard(state)
  opponentCellGrid := state.GetBoard(opponentMoveBoardReference).Cells
  opponentMove := getMoveLocations(opponentCellGrid, opponent)[0]
  state.Place(&game.Move{Board: opponentMoveBoardReference, RowNumber: opponentMove.RowNumber, ColumnNumber: opponentMove.ColumnNumber})

  done, winner = state.GetWinState()
  if (done) {
    if (winner == currentPlayer) { return GameScore - 1 } else { return -(GameScore - 1) }
  }

  score := (getCellGridScore(state.GetSuperBoard(), currentPlayer) - getCellGridScore(state.GetSuperBoard(), opponent)) * SuperBoardMultiplier
  for _, rowNumber := range ai.SideNumbers() {
    for _, columnNumber := range ai.SideNumbers() {
      board := state.GetBoard(&game.BoardReference{RowNumber: rowNumber, ColumnNumber: columnNumber})
      if board.Done {
        if int(board.Owner) == int(currentPlayer) { score += BoardScore } else { score -= BoardScore }
      } else {
        score += getCellGridScore(board.Cells, currentPlayer) - getCellGridScore(board.Cells, opponent)
      }
    }
  }
  return score
}

func getCellGridScore(cellGrid game.CellGrid, player game.Player) int {
 if len(getWinLocations(cellGrid, player)) != 0 {
  return WinLocationScore
 }
 if len(getForkLocations(cellGrid, player)) != 0 {
  return ForkLocationScore
 }
 return 0
}

func getMoveBoard(state *game.State) *game.BoardReference {
  if state.GetActiveBoard() != nil {
    return state.GetActiveBoard()
  }

  boardLocation := getMoveLocations(state.GetSuperBoard(), state.GetCurrentPlayer())[0]
  return &game.BoardReference{RowNumber: boardLocation.RowNumber, ColumnNumber: boardLocation.ColumnNumber}
}

func getMoveLocations(cellGrid game.CellGrid, me game.Player) []ai.Location {
  var opponent game.Player
  if me == game.Cell_X { opponent = game.Cell_O } else { opponent = game.Cell_X }

  winLocations := getWinLocations(cellGrid, me)
  if len(winLocations) != 0 {
    return winLocations
  }

  opponentWinLocations := getWinLocations(cellGrid, opponent)
  if len(opponentWinLocations) != 0 {
    return opponentWinLocations
  }

  forkLocations := getForkLocations(cellGrid, me)
  if len(forkLocations) != 0 {
    return forkLocations
  }

  opponentForkLocations := getForkLocations(cellGrid, opponent)
  if len(opponentForkLocations) == 1 {
    return opponentForkLocations
  }

  locationThatAvoidsForkLocations := getLocationsOfLineExcludingLocations(cellGrid, me, opponentForkLocations)
  if len(locationThatAvoidsForkLocations) != 0 {
    return locationThatAvoidsForkLocations
  }
  
  if cellGrid.IsEmpty(1, 1) {
    return []ai.Location { ai.Location { RowNumber: 1, ColumnNumber: 1 } }
  }

  return getLocations(cellGrid)
}

func getWinLocations(cellGrid game.CellGrid, player game.Player) []ai.Location {
  winLines := getOpenLines(cellGrid, player, 2)
  if len(winLines) == 0 {
    return []ai.Location {}
  }
  return getCombinedEmptyLocations(winLines, cellGrid)
}

func getForkLocations(cellGrid game.CellGrid, player game.Player) []ai.Location {
  openLines := getOpenLines(cellGrid, player, 1)
  var forkLocations []ai.Location
  var noneLocations []ai.Location
  for _, line := range openLines {
    for _, noneLocationInLine := range line.getEmptyLocations(cellGrid) {
      if slices.Contains(noneLocations, noneLocationInLine) {
        forkLocations = append(forkLocations, noneLocationInLine)
      }
      noneLocations = append(noneLocations, noneLocationInLine)
    }
  }
  return forkLocations
}

func getLocationsOfLineExcludingLocations(cellGrid game.CellGrid, player game.Player, excludeLocations []ai.Location) []ai.Location {
  var locations []ai.Location
  openLines := getOpenLines(cellGrid, player, 1)
  for _, line := range openLines {
    if !line.containsAny(excludeLocations) {
      for _, location := range line.getEmptyLocations(cellGrid) {
        if !slices.Contains(locations, location) {
          locations = append(locations, location)
        }
      }
    }
  }
  return locations
}

func getLocations(cellGrid game.CellGrid) []ai.Location {
  var locations []ai.Location
  for _, rowNumber := range ai.SideNumbers() {
    for _, columnNumber := range ai.SideNumbers() {
      if cellGrid.IsEmpty(rowNumber, columnNumber) {
        locations = append(locations, ai.Location { RowNumber: rowNumber, ColumnNumber: columnNumber })
      }
    }
  }
  return locations
}

func getCombinedEmptyLocations(lines []line, cellGrid game.CellGrid) []ai.Location {
  var emptyLocations []ai.Location
  for _, line := range lines {
    for _, location := range line.getEmptyLocations(cellGrid) {
      if (!slices.Contains(emptyLocations, location)) {
        emptyLocations = append(emptyLocations, location)
      }
    }
  }
  return emptyLocations
}

func (line *line) getEmptyLocations(cellGrid game.CellGrid) []ai.Location {
  var emptyLocations []ai.Location
  for _, location := range line.locations {
    if cellGrid.IsEmpty(location.RowNumber, location.ColumnNumber) {
      emptyLocations = append(emptyLocations, location)
    }
  }
  return emptyLocations
}

func (line *line) containsAny(locations []ai.Location) bool {
  for _, location := range locations {
    if line.contains(location) {
      return true
    }
  }
  return false
}

func (line *line) contains(location ai.Location) bool {
  for _, lineLocation := range line.locations {
    if lineLocation == location {
      return true
    }
  }
  return false
}

func getOpenLines(cellGrid game.CellGrid, player game.Player, targetPlayerCount byte) []line {
  var lines []line
  for _, line := range allLines() {
    var playerCount byte
    var isBlocked bool
    for _, location := range line.locations {
      switch game.Player(cellGrid.GetCell(location.RowNumber, location.ColumnNumber)) {
      case player: playerCount++
      case game.Cell_None: isBlocked = isBlocked || !cellGrid.IsEmpty(location.RowNumber, location.ColumnNumber)
      default: isBlocked = true
      }
    }
    if !isBlocked && (playerCount >= targetPlayerCount) {
      lines = append(lines, line)
    }
  }
  return lines
}

type line struct {
  locations [game.Size]ai.Location
}

func getAllLines() []line {
  allLines := make([]line, game.Size * 2 + 2)
  var nextLine line
  lineIndex := 0
  for _, rowNumber := range ai.SideNumbers() {
    nextLine = line {}
    for index, columnNumber := range ai.SideNumbers() {
      nextLine.locations[index] = ai.Location { RowNumber: rowNumber, ColumnNumber: columnNumber }
    }
    allLines[lineIndex] = nextLine
    lineIndex++
  }

  for _, columnNumber := range ai.SideNumbers() {
    nextLine = line {}
    for index, rowNumber := range ai.SideNumbers() {
      nextLine.locations[index] = ai.Location { RowNumber: rowNumber, ColumnNumber: columnNumber }
    }
    allLines[lineIndex] = nextLine
    lineIndex++
  }

  nextLine = line {}
  for index, rowNumber := range ai.SideNumbers() {
    nextLine.locations[index] = ai.Location { RowNumber: rowNumber, ColumnNumber: rowNumber }
  }
  allLines[lineIndex] = nextLine
  lineIndex++

  nextLine = line {}
  for index, rowNumber := range ai.SideNumbers() {
    nextLine.locations[index] = ai.Location { RowNumber: rowNumber, ColumnNumber: game.Size - rowNumber - 1 }
  }
  allLines[lineIndex] = nextLine

  return allLines
}

var allLines = sync.OnceValue(getAllLines)
