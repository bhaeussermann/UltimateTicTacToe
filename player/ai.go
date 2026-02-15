package player

import (
	"math"
	"slices"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type AI struct {}

func (*AI) GetMove(state *game.State) (*game.Move, bool) {
  potentialMoves := getPotentialMoves(state)
  
  var bestMove *game.Move
  bestMoveScore := math.MinInt
  for _, potentialMove := range potentialMoves {
    nextState := state.Copy()
    nextState.Place(&potentialMove)

    score := getScore(nextState, state.GetCurrentPlayer())
    if (score > bestMoveScore) {
      bestMove = &potentialMove
      bestMoveScore = score
    }
  }

  return bestMove, true
}

func getPotentialMoves(state *game.State) []game.Move {
	var potentialMoves []game.Move
  activeBoardReference := state.GetActiveBoard()
	if activeBoardReference != nil {
		for _, location := range getPotentialMoveLocations(state.GetBoard(activeBoardReference).Cells, state.GetCurrentPlayer()) {
			potentialMoves = append(potentialMoves, game.Move{Board: activeBoardReference, RowNumber: location.rowNumber, ColumnNumber: location.columnNumber})
		}
	} else {
    for _, potentialMoveBoard := range getPotentialMoveLocations(state.GetSuperBoard(), state.GetCurrentPlayer()) {
      boardReference := game.BoardReference{RowNumber: potentialMoveBoard.rowNumber, ColumnNumber: potentialMoveBoard.columnNumber}
      board := state.GetBoard(&boardReference)
      for _, location := range getPotentialMoveLocations(board.Cells, state.GetCurrentPlayer()) {
        potentialMoves = append(potentialMoves, game.Move{Board: &boardReference, RowNumber: location.rowNumber, ColumnNumber: location.columnNumber})
      }
		}
	}

  return potentialMoves
}

func getPotentialMoveLocations(cellGrid game.CellGrid, player game.Player) []location {
  potentialMoveLocations := getMoveLocations(cellGrid, player)

  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      if cellGrid.GetCell(rowNumber, columnNumber) == game.None {
        location := location{rowNumber: rowNumber, columnNumber: columnNumber}
        if !slices.Contains(potentialMoveLocations, location) {
          potentialMoveLocations = append(potentialMoveLocations, location)
        }
      }
    }
  }
  return potentialMoveLocations
}

func getScore(state *game.State, me game.Player) int {
  done, winner := state.GetWinState()
  if (done) {
    if (winner == me) { return game.Size * game.Size * 9 + 1 } else { return game.Size * game.Size * -9 - 1 }
  }

  var opponent game.Player
  if me == game.X { opponent = game.O } else { opponent = game.X }

  opponentMoveBoardReference := getMoveBoard(state)
  opponentCellGrid := state.GetBoard(opponentMoveBoardReference).Cells
  opponentMove := getMoveLocations(opponentCellGrid, opponent)[0]
  state.Place(&game.Move{Board: opponentMoveBoardReference, RowNumber: opponentMove.rowNumber, ColumnNumber: opponentMove.columnNumber})

  done, winner = state.GetWinState()
  if (done) {
    if (winner == me) { return game.Size * game.Size * 9 } else { return game.Size * game.Size * -9 }
  }

  score := (getCellGridScore(state.GetSuperBoard(), me) - getCellGridScore(state.GetSuperBoard(), opponent)) * game.Size * game.Size * 3
  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      board := state.GetBoard(&game.BoardReference{RowNumber: rowNumber, ColumnNumber: columnNumber})
      if board.Done {
        if int(board.Owner) == int(me) { score += 3 } else { score -= 3 }
      } else {
        score += getCellGridScore(board.Cells, me) - getCellGridScore(board.Cells, opponent)
      }
    }
  }
  return score
}

func getCellGridScore(cellGrid game.CellGrid, player game.Player) int {
 if len(getWinLocations(cellGrid, player)) != 0 {
  return 2
 }
 if len(getForkLocations(cellGrid, player)) != 0 {
  return 1
 }
 return 0
}

func getMoveBoard(state *game.State) *game.BoardReference {
  if state.GetActiveBoard() != nil {
    return state.GetActiveBoard()
  }

  boardLocation := getMoveLocations(state.GetSuperBoard(), state.GetCurrentPlayer())[0]
  return &game.BoardReference{RowNumber: boardLocation.rowNumber, ColumnNumber: boardLocation.columnNumber}
}

func getMoveLocations(cellGrid game.CellGrid, me game.Player) []location {
  var opponent game.Player
  if me == game.X { opponent = game.O } else { opponent = game.X }
  
  if cellGrid.GetCell(1, 1) == game.None {
    return []location { location { rowNumber: 1, columnNumber: 1 } }
  }

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

  return getLocations(cellGrid)
}

func getWinLocations(cellGrid game.CellGrid, player game.Player) []location {
  winLines := getOpenLines(cellGrid, player, 2)
  if len(winLines) == 0 {
    return []location {}
  }
  return getCombinedNoneLocations(winLines, cellGrid)
}

func getForkLocations(cellGrid game.CellGrid, player game.Player) []location {
  openLines := getOpenLines(cellGrid, player, 1)
  var forkLocations []location
  var noneLocations []location
  for _, line := range openLines {
    for _, noneLocationInLine := range line.getNoneLocations(cellGrid) {
      if slices.Contains(noneLocations, noneLocationInLine) {
        forkLocations = append(forkLocations, noneLocationInLine)
      }
      noneLocations = append(noneLocations, noneLocationInLine)
    }
  }
  return forkLocations
}

func getLocationsOfLineExcludingLocations(cellGrid game.CellGrid, player game.Player, excludeLocations []location) []location {
  var locations []location
  openLines := getOpenLines(cellGrid, player, 1)
  for _, line := range openLines {
    if !line.containsAny(excludeLocations) {
      for _, location := range line.getNoneLocations(cellGrid) {
        if !slices.Contains(locations, location) {
          locations = append(locations, location)
        }
      }
    }
  }
  return locations
}

func getLocations(cellGrid game.CellGrid) []location {
  var locations []location
  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      if cellGrid.GetCell(rowNumber, columnNumber) == game.None {
        locations = append(locations, location { rowNumber: rowNumber, columnNumber: columnNumber })
      }
    }
  }
  return locations
}

func getCombinedNoneLocations(lines []line, cellGrid game.CellGrid) []location {
  var noneLocations []location
  for _, line := range lines {
    for _, location := range line.getNoneLocations(cellGrid) {
      if (!slices.Contains(noneLocations, location)) {
        noneLocations = append(noneLocations, location)
      }
    }
  }
  return noneLocations
}

func (line *line) getNoneLocations(cellGrid game.CellGrid) []location {
  var noneLocations []location
  for _, location := range line.locations {
    if cellGrid.GetCell(location.rowNumber, location.columnNumber) == game.None {
      noneLocations = append(noneLocations, location)
    }
  }
  return noneLocations
}

func (line *line) containsAny(locations []location) bool {
  for _, location := range locations {
    if line.contains(location) {
      return true
    }
  }
  return false
}

func (line *line) contains(location location) bool {
  for _, lineLocation := range line.locations {
    if lineLocation == location {
      return true
    }
  }
  return false
}

func getOpenLines(cellGrid game.CellGrid, player game.Player, targetPlayerCount byte) []line {
  var lines []line
  for _, line := range allLines {
    var playerCount byte
    var hasOpponent bool
    for _, location := range line.locations {
      switch game.Player(cellGrid.GetCell(location.rowNumber, location.columnNumber)) {
      case player: playerCount++
      case game.None: 
      default: hasOpponent = true
      }
    }
    if !hasOpponent && (playerCount >= targetPlayerCount) {
      lines = append(lines, line)
    }
  }
  return lines
}

type line struct {
  locations [game.Size]location
}

type location struct {
  rowNumber byte
  columnNumber byte
}

var sideNumbers [game.Size]byte
var allLines [game.Size * 2 + 2]line

func init() {
  var number byte
  for number = 0; number < game.Size; number++ {
    sideNumbers[number] = number
  }

  var nextLine line
  lineIndex := 0
  for _, rowNumber := range sideNumbers {
    nextLine = line {}
    for index, columnNumber := range sideNumbers {
      nextLine.locations[index] = location { rowNumber: rowNumber, columnNumber: columnNumber }
    }
    allLines[lineIndex] = nextLine
    lineIndex++
  }

  for _, columnNumber := range sideNumbers {
    nextLine = line {}
    for index, rowNumber := range sideNumbers {
      nextLine.locations[index] = location { rowNumber: rowNumber, columnNumber: columnNumber }
    }
    allLines[lineIndex] = nextLine
    lineIndex++
  }

  nextLine = line {}
  for index, rowNumber := range sideNumbers {
    nextLine.locations[index] = location { rowNumber: rowNumber, columnNumber: rowNumber }
  }
  allLines[lineIndex] = nextLine
  lineIndex++

  nextLine = line {}
  for index, rowNumber := range sideNumbers {
    nextLine.locations[index] = location { rowNumber: rowNumber, columnNumber: game.Size - rowNumber - 1 }
  }
  allLines[lineIndex] = nextLine
}
