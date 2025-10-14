package game

const Size = 3

type Board [Size][Size]Cell

type Cell byte

const (
  None = iota
  X
  O
)

func (board *Board) ToString() string {
  boardAsString := ""
  for rowIndex, row := range board {
    for cellIndex, cell := range row {
      switch cell {
        case None: boardAsString += "  "
        case X: boardAsString += " X"
        case O: boardAsString += " O"
      }
      if (cellIndex < Size - 1) {
        boardAsString += " |"
      }
    }
    if (rowIndex < Size - 1) {
      boardAsString += "\r\n"
      for i := 0; i < Size * 4 - 1; i++ {
        boardAsString += "—"
      }
      boardAsString += "\r\n"
    }
  }
  return boardAsString
}
