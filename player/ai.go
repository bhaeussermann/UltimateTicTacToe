package player

import (
	"slices"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
)

type AI struct {}

func (*AI) GetMove(state *game.State) (*game.Move, bool) {
  moveLocation := getMoveLocation(state)
  return &game.Move { RowNumber: moveLocation.rowNumber, ColumnNumber: moveLocation.columnNumber }, true
}

func getMoveLocation(state *game.State) location {
  me := state.GetCurrentPlayer()
  var opponent game.Player
  if me == game.X { opponent = game.O } else { opponent = game.X }
  
  board := state.GetBoard()

  winLocation := getWinLocation(board, me)
  if winLocation != nil {
    return *winLocation
  }

  opponentWinLocation := getWinLocation(board, opponent)
  if opponentWinLocation != nil {
    return *opponentWinLocation
  }

  forkLocation := getForkLocation(board, me)
  if forkLocation != nil {
    return *forkLocation
  }

  location := getFirstLocation(board, []location {
    location { rowNumber: 1, columnNumber: 1 },
    location { rowNumber: 0, columnNumber: 0 },
    location { rowNumber: 0, columnNumber: game.Size - 1 },
    location { rowNumber: game.Size - 1, columnNumber: 0 },
    location { rowNumber: game.Size - 1, columnNumber: game.Size - 1 },
  })
  if (location != nil) {
    return *location
  }

  return *getLocation(board)
}

func getWinLocation(board *game.Board, player game.Player) *location {
  winLines := getOpenLines(board, player, 2)
  if len(winLines) == 0 {
    return nil
  }
  return &winLines[0].getNoneLocations(board)[0]
}

func getForkLocation(board *game.Board, player game.Player) *location {
  openLines := getOpenLines(board, player, 1)
  var noneLocations []location
  for _, line := range openLines {
    for _, noneLocationInLine := range line.getNoneLocations(board) {
      if slices.Contains(noneLocations, noneLocationInLine) {
        return &noneLocationInLine
      }
      noneLocations = append(noneLocations, noneLocationInLine)
    }
  }
  return nil
}

func getFirstLocation(board *game.Board, locations []location) *location {
  for _, location := range locations {
    if board[location.rowNumber][location.columnNumber] == game.None {
      return &location
    }
  }
  return nil
}

func getLocation(board *game.Board) *location {
  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      if board[rowNumber][columnNumber] == game.None {
        return &location { rowNumber: rowNumber, columnNumber: columnNumber }
      }
    }
  }
  return nil
}

func (line *line) getNoneLocations(board *game.Board) []location {
  var noneLocations []location
  for _, location := range line.locations {
    if board[location.rowNumber][location.columnNumber] == game.None {
      noneLocations = append(noneLocations, location)
    }
  }
  return noneLocations
}

func getOpenLines(board *game.Board, player game.Player, targetPlayerCount byte) []line {
  var lines []line
  for _, line := range allLines {
    var playerCount byte
    var hasOpponent bool
    for _, location := range line.locations {
      switch game.Player(board[location.rowNumber][location.columnNumber]) {
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
