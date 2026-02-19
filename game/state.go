package game

type State struct {
  superBoard SuperBoard
  currentPlayer Player
  activeBoard *BoardReference
  done bool
  winner Player
}

type Player byte

type Move struct {
  Board *BoardReference
  RowNumber byte
  ColumnNumber byte
}

func CreateState() *State {
  return &State {
    currentPlayer: Cell_X,
    winner: Cell_None,
  }
}

func (state *State) Copy() *State {
  return &State {
    superBoard: *state.superBoard.Copy(),
    currentPlayer: state.currentPlayer,
    activeBoard: state.activeBoard,
    done: state.done,
    winner: state.winner,
  }
}

func (state *State) GetSuperBoard() *SuperBoard {
  return &state.superBoard
}

func (state *State) GetActiveBoard() *BoardReference {
  return state.activeBoard
}

func (state *State) GetBoard(boardReference *BoardReference) *Board {
  return &state.superBoard[boardReference.RowNumber][boardReference.ColumnNumber]
}

func (state *State) GetCurrentPlayer() Player {
  return state.currentPlayer
}

func (state *State) GetWinState() (bool, Player) {
  return state.done, state.winner
}

func (state *State) CanPlaceIn(boardReference *BoardReference) bool {
  if (state.activeBoard == nil) {
    return !state.GetBoard(boardReference).Done
  }
  return &state.activeBoard == &boardReference
}

func (state *State) CanPlace(move *Move) bool {
  if state.done || (state.activeBoard == nil && move.Board == nil) {
    return false
  }
  
  board := state.getBoard(move)
  return !board.Done && board.Cells[move.RowNumber][move.ColumnNumber] == Cell_None
}

func (state *State) Place(move *Move) bool {
  if !state.CanPlace(move) {
    return false
  }

  board := state.getBoard(move)
  board.Cells[move.RowNumber][move.ColumnNumber] = Cell(state.currentPlayer)
  board.updateBoardOwner()
  state.updateWinState()

  if state.superBoard[move.RowNumber][move.ColumnNumber].Done {
    state.activeBoard = nil
  } else {
    state.activeBoard = &BoardReference{ RowNumber: move.RowNumber, ColumnNumber: move.ColumnNumber }
  }

  if state.currentPlayer == Cell_X {
    state.currentPlayer = Cell_O
  } else {
    state.currentPlayer = Cell_X
  }
  return true
}

func (state *State) getBoard(move *Move) *Board {
  var boardReference *BoardReference
  if state.activeBoard != nil { boardReference = state.activeBoard } else { boardReference = move.Board }
  return &state.superBoard[boardReference.RowNumber][boardReference.ColumnNumber]
}

func (board *Board) updateBoardOwner() {
  if hasAnyLineFilled(board.Cells, Cell_X) {
    board.Done = true
    board.Owner = Cell_X
  }
  if hasAnyLineFilled(board.Cells, Cell_O) {
    board.Done = true
    board.Owner = Cell_O
  }
  if !hasEmptyCell(board.Cells) {
    board.Done = true
    board.Owner = Cell_None
  }
}

func (state *State) updateWinState() {
  if hasAnyLineFilled(state.superBoard, Cell_X) {
    state.winner = Cell_X
    state.done = true
  }
  if hasAnyLineFilled(state.superBoard, Cell_O) {
    state.winner = Cell_O
    state.done = true
  }
  if !hasEmptyCell(state.superBoard) {
    state.done = true
  }
}

func hasAnyLineFilled(cellGrid CellGrid, player Player) bool {
  return hasAnyRowFilled(cellGrid, player) || hasAnyColumnFilled(cellGrid, player) || hasDiagonal1Filled(cellGrid, player) || hasDiagonal2Filled(cellGrid, player)
}

func hasAnyRowFilled(cellGrid CellGrid, player Player) bool {
  for _, rowNumber := range sideNumbers {
    if hasRowFilled(cellGrid, rowNumber, player) {
      return true;
    }
  }
  return false;
}

func hasRowFilled(cellGrid CellGrid, rowNumber byte, player Player) bool {
  for _, columnNumber := range sideNumbers {
    if Player(cellGrid.GetCell(rowNumber, columnNumber)) != player {
      return false;
    }
  }
  return true;
}

func hasAnyColumnFilled(cellGrid CellGrid, player Player) bool {
  for _, columnNumber := range sideNumbers {
    if hasColumnFilled(cellGrid, columnNumber, player) {
      return true
    }
  }
  return false
}

func hasColumnFilled(cellGrid CellGrid, columnNumber byte, player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(cellGrid.GetCell(rowNumber, columnNumber)) != player {
      return false;
    }
  }
  return true;
}

func hasDiagonal1Filled(cellGrid CellGrid, player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(cellGrid.GetCell(rowNumber, rowNumber)) != player {
      return false;
    }
  }
  return true;
}

func hasDiagonal2Filled(cellGrid CellGrid, player Player) bool {
  for _, rowNumber := range sideNumbers {
    if Player(cellGrid.GetCell(rowNumber, Size - rowNumber - 1)) != player {
      return false;
    }
  }
  return true;
}

func hasEmptyCell(cellGrid CellGrid) bool {
  for _, rowNumber := range sideNumbers {
    for _, columnNumber := range sideNumbers {
      if cellGrid.IsEmpty(rowNumber, columnNumber) {
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
