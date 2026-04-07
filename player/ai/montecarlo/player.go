package montecarlo

import (
	"math"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
)

type Player struct {
  Difficulty ai.Difficulty
}

type gameCounts struct {
  totalGames uint32
  totalWins uint32
  totalDraws uint32
}

func (p *Player) GetMove(state *game.State, log player.Log) (player.Action, *game.Move) {
  done, _ := state.GetWinState()
  if done {
    return player.Action_None, nil
  }

  root := createNode(state, nil, nil)
  treeMutex := sync.Mutex {}
  gameCounts := gameCounts {}
  isRunning := true
  workerCount := runtime.NumCPU()
  var waitGroup sync.WaitGroup
  waitGroup.Add(workerCount)
  for count := 0; count < workerCount; count++ {
    go search(root, &treeMutex, state, &gameCounts, &isRunning, &waitGroup)
  }

  time.Sleep(p.getTimeoutDuration())
  isRunning = false
  waitGroup.Wait()

  log.Logf("Simulations: %v\r\n", gameCounts.totalGames)
  log.Logf("Wins: %.1f %%\r\n", float64(gameCounts.totalWins + gameCounts.totalDraws / 2) / float64(gameCounts.totalGames) * 100)

  var maximumWinRatio float32 = -1
  var bestMove game.Move
  for _, node := range root.children {
    winRatio := float32(node.winCount) / float32(node.gameCount)
    if winRatio > maximumWinRatio {
      maximumWinRatio = winRatio
      bestMove = *node.lastMove
    }
  }
  return player.Action_Move, &bestMove
}

func search(root *node, treeMutex *sync.Mutex, state *game.State, gameCounts *gameCounts, isRunning *bool, waitGroup *sync.WaitGroup) {
  defer waitGroup.Done()
	for *isRunning {
    treeMutex.Lock()
    leaf := selectLeaf(root)

    var child *node
    var winner game.Player
    isTerminalNode := len(leaf.potentialMovesToExpand) == 0
    if isTerminalNode {
      treeMutex.Unlock()
      child = leaf
      _, winner = child.state.GetWinState()
    } else {
      child = expand(leaf)
      treeMutex.Unlock()
      winner = play(child.state)
    }

    backpropagate(child, winner)

		atomic.AddUint32(&gameCounts.totalGames, 1)
		if winner == state.GetCurrentPlayer() {
			atomic.AddUint32(&gameCounts.totalWins, 1)
		}
		if winner == game.Cell_None {
			atomic.AddUint32(&gameCounts.totalDraws, 1)
		}
	}
}

func (p *Player) getTimeoutDuration() time.Duration {
  switch p.Difficulty {
  case ai.Difficulty_Easy: return time.Second
  case ai.Difficulty_Medium: return time.Second * 2
  default: return time.Second * 5
  }
}

func selectLeaf(root *node) *node {
  currentNode := root
  for ; (len(currentNode.potentialMovesToExpand) == 0) && (len(currentNode.children) != 0); {
    var maximumChildScore float64 = -1
    var maximumChild *node
    for _, child := range currentNode.children {
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

func expand(n *node) *node {
  move := n.potentialMovesToExpand[0]
  n.potentialMovesToExpand = slices.Delete(n.potentialMovesToExpand, 0, 1)

  childState := n.state.Copy()
  childState.Place(move)
  childNode := createNode(childState, move, n)
  n.children = append(n.children, childNode)
  return childNode
}

func backpropagate(leaf *node, winner game.Player) {
  if winner == game.Cell_None {
    if rand.Intn(2) == 0 { winner = game.Cell_X } else { winner = game.Cell_O }
  }

  for currentNode := leaf; currentNode != nil; currentNode = currentNode.parent {
    if (currentNode.parent != nil) && (currentNode.parent.state.GetCurrentPlayer() == winner) {
      atomic.AddUint32(&currentNode.winCount, 1)
    }
    atomic.AddUint32(&currentNode.gameCount, 1)
  }
}

type node struct {
  state *game.State
  lastMove *game.Move
  winCount uint32
  gameCount uint32
  parent *node
  children []*node
  potentialMovesToExpand []*game.Move
}

func createNode(state *game.State, lastMove *game.Move, parent *node) *node {
  return &node {
    state: state,
    lastMove: lastMove,
    parent: parent,
    children: []*node{},
    potentialMovesToExpand: getPotentialMoves(state),
  }
}

func play(state *game.State) game.Player {
  done, winner := state.GetWinState()
  if done { return winner }

  stateCopy := state.Copy()
  for ; !done; done, winner = stateCopy.GetWinState() {
    potentialMoves := getPotentialMoves(stateCopy)
    move := potentialMoves[rand.Intn(len(potentialMoves))]
    stateCopy.Place(move)
  }
  return winner
}

func getPotentialMoves(state *game.State) []*game.Move {
	var potentialMoves []*game.Move

  done, _ := state.GetWinState()
  if done { return potentialMoves }

  activeBoardReference := state.GetActiveBoard()
	if activeBoardReference != nil {
		for _, location := range getPotentialMoveLocations(state.GetBoard(activeBoardReference).Cells) {
			potentialMoves = append(potentialMoves, &game.Move{Board: activeBoardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
		}
	} else {
    for _, potentialMoveBoard := range getPotentialMoveLocations(state.GetSuperBoard()) {
      boardReference := game.BoardReference{RowNumber: potentialMoveBoard.RowNumber, ColumnNumber: potentialMoveBoard.ColumnNumber}
      board := state.GetBoard(&boardReference)
      for _, location := range getPotentialMoveLocations(board.Cells) {
        potentialMoves = append(potentialMoves, &game.Move{Board: &boardReference, RowNumber: location.RowNumber, ColumnNumber: location.ColumnNumber})
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
