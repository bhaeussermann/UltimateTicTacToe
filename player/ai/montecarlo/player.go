package montecarlo

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
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

  totalGames := 0
  totalWins := 0
  totalDraws := 0

  root := node { state: state }
  for 
  deadline := time.Now().Add(time.Second * 2);
  time.Now().Before(deadline) && !root.isFullyExpanded; {
    leaf := selectLeaf(&root)
    children := expand(leaf)
    for _, child := range children {
      winner := play(child.state)

      totalGames++
      if winner == state.GetCurrentPlayer() {
        totalWins++
      }
      if winner == game.Cell_None {
        totalDraws++
      }

      backpropagate(child, winner)

      isChildStateDone, _ := child.state.GetWinState()
      if isChildStateDone { setIsFullyExpanded(child) }
    }
  }

  fmt.Printf("Total games: %v\r\n", totalGames)
  fmt.Printf("Wins: %.1f %%\r\n", float64(totalWins + totalDraws / 2) / float64(totalGames) * 100)

  var maximumWinRatio float32 = -1
  var bestMove game.Move
  for _, node := range root.children {
    winRatio := float32(node.winCount) / float32(node.gameCount)
    if winRatio > maximumWinRatio {
      maximumWinRatio = winRatio
      bestMove = *node.move
    }
  }
  return player.Action_Move, &bestMove
}

func selectLeaf(root *node) *node {
  currentNode := root
  for ; currentNode.children != nil; {
    var maximumChildScore float64 = -1
    var maximumChild *node
    for _, child := range currentNode.children {
      if child.isFullyExpanded { continue }
      nodeScore := float64(child.winCount) / float64(child.gameCount) + explorationFactor * math.Sqrt(math.Log2(float64(currentNode.gameCount)) / float64(child.gameCount))
      if nodeScore > maximumChildScore {
        maximumChildScore = nodeScore
        maximumChild = child
      }
    }
    currentNode = maximumChild
  }
  return currentNode
}

const explorationFactor = math.Sqrt2

func expand(n *node) []*node {
  potentialMoves := getPotentialMoves(n.state)
  n.children = make([]*node, len(potentialMoves))
  for index, move := range potentialMoves {
    childState := n.state.Copy()
    childState.Place(&move)
    n.children[index] = &node { state: childState, move: &move, parent: n }
  }
  return slices.Clone(n.children)
}

func backpropagate(leaf *node, winner game.Player) {
  if winner == game.Cell_None {
    if rand.Intn(2) == 0 { winner = game.Cell_X } else { winner = game.Cell_O }
  }

  for currentNode := leaf; currentNode != nil; currentNode = currentNode.parent {
    if (currentNode.parent != nil) && (currentNode.parent.state.GetCurrentPlayer() == winner) {
      currentNode.winCount++
    }
    currentNode.gameCount++
  }
}

func setIsFullyExpanded(node *node) {
  for ; node != nil; node = node.parent {
    node.isFullyExpanded = true
    if node.parent != nil {
      for _, child := range node.parent.children {
        if !child.isFullyExpanded { return }
      }
    }
  }
}

type node struct {
  state *game.State
  move *game.Move
  winCount int
  gameCount int
  isFullyExpanded bool
  parent *node
  children []*node
}

func play(state *game.State) game.Player {
  done, winner := state.GetWinState()
  if done { return winner }

  stateCopy := state.Copy()
  for ; !done; done, winner = stateCopy.GetWinState() {
    potentialMoves := getPotentialMoves(stateCopy)
    move := potentialMoves[rand.Intn(len(potentialMoves))]
    stateCopy.Place(&move)
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
