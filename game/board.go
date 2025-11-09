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
  None = iota
  X
  O
)

type CellGrid interface {
  GetCell(rowNumber byte, columnNumber byte) Cell
}

func (superBoard SuperBoard) GetCell(rowNumber byte, columnNumber byte) Cell {
  return superBoard[rowNumber][columnNumber].Owner
}

func (boardCells BoardCells) GetCell(rowNumber byte, columnNumber byte) Cell {
  return boardCells[rowNumber][columnNumber]
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
          if board.Owner == X {
            if rowIndex == 0 { boardAsString += "   ╲   ╱  " }
            if rowIndex == 1 { boardAsString += "     ╳    " }
            if rowIndex == 2 { boardAsString += "   ╱   ╲  " }
          }
          if board.Owner == O {
            if rowIndex == 0 { boardAsString += "  ╭─────╮ " }
            if rowIndex == 1 { boardAsString += "  │     │ " }
            if rowIndex == 2 { boardAsString += "  ╰─────╯ " }
          }
          if board.Owner == None {
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
            if board.Owner == X {
              if rowIndex == 0 { boardAsString += "    ╲ ╱    " }
              if rowIndex == 1 { boardAsString += "    ╱ ╲    " }
            }
            if board.Owner == O {
              if rowIndex == 0 { boardAsString += "  │     │  " }
              if rowIndex == 1 { boardAsString += "  │     │  " }
            }
            if board.Owner == None {
              boardAsString += "          "
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
  case None: return " "
  case X: return "X"
  case O: return "O"
  }
  return ""
}
