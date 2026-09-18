package game

import (
	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
)

type uiPlayer struct {
	selectedMoveChannel chan *game.Move
}

func createUIPlayer() *uiPlayer {
	return &uiPlayer{
		selectedMoveChannel: make(chan *game.Move),
	}
}

func (p *uiPlayer) GetMove(state *game.State, log player.Log) (player.Action, *game.Move) {
	selectedMove := <-p.selectedMoveChannel
	if selectedMove == nil {
		return player.Action_None, nil
	}
	return player.Action_Move, selectedMove
}

func (p *uiPlayer) setMove(selectedMove *game.Move) {
	p.selectedMoveChannel <- selectedMove
}

func (p *uiPlayer) unblock() {
	p.setMove(nil)
}
