package game

type State struct {
  Board Board
  CurrentPlayer Player
  Done bool
  Winner Player
}

type Player byte

func CreateState() *State {
  var state State
  state.CurrentPlayer = X
  return &state
}

func (state *State) Place(rowNumber byte, columnNumber byte) bool {
  if state.Board[rowNumber][columnNumber] != None {
    return false
  }

  state.Board[rowNumber][columnNumber] = Cell(state.CurrentPlayer)
  if state.CurrentPlayer == X {
    state.CurrentPlayer = O
  } else {
    state.CurrentPlayer = X
  }
  return true
}
