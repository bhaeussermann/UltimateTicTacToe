package game

const Size = 3

type SuperBoard [Size][Size]Board

type BoardReference struct {
  RowNumber byte
  ColumnNumber byte
}

type Board struct {
  Done bool
  Owner Cell
  Cells BoardCells
}

type BoardCells [Size][Size]Cell

type Cell byte

const (
  Cell_None = iota
  Cell_X
  Cell_O
)

type CellGrid interface {
  GetCell(rowNumber byte, columnNumber byte) Cell
  IsEmpty(rowNumber byte, columnNumber byte) bool
}

func (superBoard SuperBoard) Copy() *SuperBoard {
  superBoardCopy := SuperBoard{}
  for boardRowIndex := 0; boardRowIndex < Size; boardRowIndex++ {
    for boardColumnIndex := 0; boardColumnIndex < Size; boardColumnIndex++ {
      board := &superBoard[boardRowIndex][boardColumnIndex]
      boardCopy := &superBoardCopy[boardRowIndex][boardColumnIndex]
      boardCopy.Done = board.Done
      boardCopy.Owner = board.Owner
      for rowIndex := 0; rowIndex < Size; rowIndex++ {
        for columnIndex := 0; columnIndex < Size; columnIndex++ {
          boardCopy.Cells[rowIndex][columnIndex] = board.Cells[rowIndex][columnIndex]
        }
      }
    }
  }
  return &superBoardCopy
}

func (superBoard SuperBoard) GetCell(rowNumber byte, columnNumber byte) Cell {
  return superBoard[rowNumber][columnNumber].Owner
}

func (superBoard SuperBoard) IsEmpty(rowNumber byte, columnNumber byte) bool {
  return !superBoard[rowNumber][columnNumber].Done
}

func (boardCells BoardCells) GetCell(rowNumber byte, columnNumber byte) Cell {
  return boardCells[rowNumber][columnNumber]
}

func (boardCells BoardCells) IsEmpty(rowNumber byte, columnNumber byte) bool {
  return boardCells[rowNumber][columnNumber] == Cell_None
}

func (superBoard *SuperBoard) GetHorizontalLine() string {
  return repeat("═", Size * (Size * 4 + 2))
}

func (superBoard *SuperBoard) ToString(activeBoard *BoardReference) string {
  boardAsString := ""
  for boardsRowIndex, boardsRow := range superBoard {
    for boardIndex := range boardsRow {
      if activeBoard.is(boardsRowIndex, boardIndex) {
        boardAsString += "┌" + repeat("─", Size * 4 - 1) + "┐  "
      } else {
        boardAsString += repeat(" ", Size * 4 + 3)
      }
    }
    boardAsString += "\r\n"

    for rowIndex := 0; rowIndex < Size; rowIndex++ {
      for boardIndex, board := range boardsRow {
        if activeBoard.is(boardsRowIndex, boardIndex) {
          boardAsString += "│"
        } else {
          boardAsString += " "
        }
        if board.Done {
          if board.Owner == Cell_X {
            if rowIndex == 0 { boardAsString += "   ╲   ╱  " }
            if rowIndex == 1 { boardAsString += "     ╳    " }
            if rowIndex == 2 { boardAsString += "   ╱   ╲  " }
          }
          if board.Owner == Cell_O {
            if rowIndex == 0 { boardAsString += "  ╭─────╮ " }
            if rowIndex == 1 { boardAsString += "  │     │ " }
            if rowIndex == 2 { boardAsString += "  ╰─────╯ " }
          }
          if board.Owner == Cell_None {
            if rowIndex == 0 { boardAsString += "          " }
            if rowIndex == 1 { boardAsString += "  ──────  " }
            if rowIndex == 2 { boardAsString += "          " }
          }
        } else {
          for cellIndex, cell := range board.Cells[rowIndex] {
            boardAsString += " " + cell.toString()
            if cellIndex < Size - 1 {
              boardAsString += " │"
            }
          }
        }
        if activeBoard.is(boardsRowIndex, boardIndex) {
          boardAsString += " │"
        } else {
          boardAsString += "  "
        }
        if boardIndex < Size - 1 {
          boardAsString += "  "
        }
      }
      if rowIndex < Size - 1 {
        boardAsString += "\r\n"
        for boardIndex, board := range boardsRow {
          if activeBoard.is(boardsRowIndex, boardIndex) {
            boardAsString += "│"
          } else {
            boardAsString += " "
          }
          if board.Done {
            if board.Owner == Cell_X {
              if rowIndex == 0 { boardAsString += "    ╲ ╱    " }
              if rowIndex == 1 { boardAsString += "    ╱ ╲    " }
            }
            if board.Owner == Cell_O {
              if rowIndex == 0 { boardAsString += "  │     │  " }
              if rowIndex == 1 { boardAsString += "  │     │  " }
            }
            if board.Owner == Cell_None {
              boardAsString += "           "
            }
          } else {
            for cellIndex := 0; cellIndex < Size; cellIndex++ {
              boardAsString += "───"
              if cellIndex < Size - 1 {
                boardAsString += "┼"
              } 
            }
          }
          if activeBoard.is(boardsRowIndex, boardIndex) {
            boardAsString += "│"
          } else {
            boardAsString += " "
          }
          if boardIndex < Size - 1 {
            boardAsString += "  "
          }
        }
        boardAsString += "\r\n"
      }
    }
    boardAsString += "\r\n"
    
    for boardIndex := range boardsRow {
      if activeBoard.is(boardsRowIndex, boardIndex) {
        boardAsString += "└" + repeat("─", Size * 4 - 1) + "┘  "
      } else {
        boardAsString += repeat(" ", Size * 4 + 3)
      }
    }
    if boardsRowIndex < Size - 1 {
      boardAsString += "\r\n\r\n"
    }
  }
  boardAsString += "\r\n"
  return boardAsString
}

func (activeBoard *BoardReference) is(rowNumber int, columnNumber int) bool {
  return (activeBoard != nil) && (byte(rowNumber) == activeBoard.RowNumber) && (byte(columnNumber) == activeBoard.ColumnNumber)
}

func (board *Board) ToString() string {
  boardAsString := ""
  for rowIndex, row := range board.Cells {
    for cellIndex, cell := range row {
      boardAsString += " " + cell.toString()
      if cellIndex < Size - 1 {
        boardAsString += " │"
      }
    }
    if (rowIndex < Size - 1) {
      boardAsString += "\r\n"
      for cellIndex := range row {
        boardAsString += "───"
        if cellIndex < Size - 1 {
          boardAsString += "┼"
        } 
      }
      boardAsString += "\r\n"
    }
  }
  return boardAsString
}

func repeat(character string, length int) string {
  result := ""
  for count := 0; count < length; count++ {
    result += character
  }
  return result
}

func (cell Cell) toString() string {
  switch cell {
  case Cell_None: return " "
  case Cell_X: return "X"
  case Cell_O: return "O"
  }
  return ""
}
