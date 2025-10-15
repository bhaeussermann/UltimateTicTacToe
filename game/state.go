package game

type State struct {
  board Board
  currentPlayer Player
  done bool
  winner Player
}

type Player byte

func CreateState() *State {
  var state State
  state.currentPlayer = X
  state.winner = None
  return &state
}

func (state *State) GetBoard() *Board {
  return &state.board
}

func (state *State) GetCurrentPlayer() Player {
  return state.currentPlayer
}

func (state *State) GetWinState() (bool, Player) {
  return state.done, state.winner
}

func (state *State) Place(rowNumber byte, columnNumber byte) bool {
  if state.done || state.board[rowNumber][columnNumber] != None {
    return false
  }

  state.board[rowNumber][columnNumber] = Cell(state.currentPlayer)
  state.updateWinState()
  if state.currentPlayer == X {
    state.currentPlayer = O
  } else {
    state.currentPlayer = X
  }
  return true
}

func (state *State) updateWinState() {
  if state.board.hasAnyLineFilled(X) {
    state.winner = X
    state.done = true
  }
  if state.board.hasAnyLineFilled(O) {
    state.winner = O
    state.done = true
  }
  if !state.board.hasEmptyCell() {
    state.done = true
  }
}

func (board *Board) hasAnyLineFilled(player Player) bool {
  return board.hasAnyRowFilled(player) || board.hasAnyColumnFilled(player) || board.hasDiagonal1Filled(player) || board.hasDiagonal2Filled(player)
}

func (board *Board) hasAnyRowFilled(player Player) bool {
  for _, rowNumber := range sideNumbers {
    if board.hasRowFilled(rowNumber, player) {
      return true;
    }
  }
  return false;
}

func (board *Board) hasRowFilled(rowNumber byte, player Player) bool {
  for _, columnNumber := range sideNumbers {
    if Player(board[rowNumber][columnNumber]) != player {
      return false;
    }
  }
  return true;
}

func (board *Board) hasAnyColumnFilled(player Player) bool {
  for _, columnNumber := range sideNumbers {
    if board.hasColumnFilled(columnNumber, player) {
      return true
    }
  }
  return false
}

func (board *Board) hasColumnFilled(columnNumber byte, player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(board[rowNumber][columnNumber]) != player {
      return false;
    }
  }
  return true;
}

func (board *Board) hasDiagonal1Filled(player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(board[rowNumber][rowNumber]) != player {
      return false;
    }
  }
  return true;
}

func (board *Board) hasDiagonal2Filled(player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(board[rowNumber][Size - rowNumber - 1]) != player {
      return false;
    }
  }
  return true;
}

func (board *Board) hasEmptyCell() bool {
  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      if board[rowNumber][columnNumber] == None {
        return true;
      }
    }
  }
  return false;
}

var sideNumbers [Size]byte

func init() {
  var number byte
  for number = 0; number < Size; number++ {
    sideNumbers[number] = number
  }
}
