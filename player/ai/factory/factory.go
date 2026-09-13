package factory

import (
	"fmt"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/alphabeta"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/montecarlo"
)

func CreateAIPlayer(aiDifficulty ai.Difficulty) player.Player {
	switch aiDifficulty {
  case ai.Difficulty_Easy:
		return &alphabeta.Player{Difficulty: ai.Difficulty_Easy}
  case ai.Difficulty_Medium:
		return &alphabeta.Player{Difficulty: ai.Difficulty_Hard}
	case ai.Difficulty_Hard:
		return &montecarlo.Player{Difficulty: ai.Difficulty_Hard}
  default:
			panic(fmt.Sprintf("Unhandled AI difficulty: %d", aiDifficulty))
	}
}
