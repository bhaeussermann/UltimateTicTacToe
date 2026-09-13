package game

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/gui/scenes"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	textFaceSource *text.GoTextFaceSource
	screenWidth, screenHeight int
	exitToScene scenes.GetNextScene

	playerSelection game.Player
	playerX, playerO player.Player
	state *game.State
	activeUiPlayer *uiPlayer
	isClicking bool
}

func NewGame(playerSelection game.Player, aiDifficulty ai.Difficulty, exitToScene scenes.GetNextScene) (scenes.Scene, error) {
	textFaceSource, error := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if error != nil {
		return nil, error
	}

	playerX, playerO := getPlayers(playerSelection, aiDifficulty)
	game := &Game{
		textFaceSource: textFaceSource,
		exitToScene: exitToScene,
		playerSelection: playerSelection,
		playerX: playerX,
		playerO: playerO,
		state: game.CreateState(),
	}
	go game.play()
	return game, nil
}

func (g *Game) SetScreenSize(width int, height int) {
	g.screenWidth = width
	g.screenHeight = height
}

func (g *Game) Update() scenes.SceneChange {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.isClicking = false
	} else if !g.isClicking {
		g.isClicking = true
		cursorX, cursorY := ebiten.CursorPosition()
		cursorXf, cursorYf := float32(cursorX), float32(cursorY)
		exitClicked := (float32(g.screenWidth) - exitButtonSize - margin <= cursorXf && cursorXf <= float32(g.screenWidth) - margin) && (margin <= cursorYf && cursorYf <= exitButtonSize + margin)
		if exitClicked {
			if g.activeUiPlayer != nil {
				g.activeUiPlayer.unblock()
			}
			return scenes.SceneChange{GetNextScene: g.exitToScene }
		}

		if g.activeUiPlayer != nil {
			hoveredCell := g.getHoveredCell()
			if (hoveredCell != nil) && g.state.CanPlaceIn(hoveredCell.Board) && g.state.CanPlace(hoveredCell) {
				g.activeUiPlayer.setMove(hoveredCell)
			}
		}
	}

	return scenes.SceneChange{}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawSuperBoard(screen)
	drawExitButton(screen, float32(g.screenWidth) - exitButtonSize - margin, margin)
}

func (g *Game) drawSuperBoard(screen *ebiten.Image) {
	g.drawHoveredCell(screen)
	
	boardDrawSize := g.getBoardDrawSize()
	cellDrawSize := g.getCellDrawSize()
	for boardRow := range game.Size {
		boardYOffset := float32(boardRow) * boardDrawSize + cellDrawSize
		for boardColumn := range game.Size {
			boardXOffset := float32(boardColumn) * boardDrawSize + cellDrawSize

			boardReference := &game.BoardReference{RowNumber: byte(boardRow), ColumnNumber: byte(boardColumn)}
			boardForegroundColor, boardCrossColor, boardNaughtColor := g.getBoardColors(boardReference)

			g.drawGrid(screen, boardXOffset, boardYOffset, boardForegroundColor)
			board := g.state.GetBoard(boardReference)
			g.drawBoardSymbols(screen, board, boardXOffset, boardYOffset, boardCrossColor, boardNaughtColor)
		}
	}
}

func (g *Game) drawHoveredCell(screen *ebiten.Image) {
	if g.activeUiPlayer != nil {
		hoveredCell := g.getHoveredCell()
		if (hoveredCell != nil) && g.state.CanPlaceIn(hoveredCell.Board) && g.state.CanPlace(hoveredCell) {
			boardDrawSize := g.getBoardDrawSize()
			cellDrawSize := g.getCellDrawSize()
			x := float32(hoveredCell.Board.ColumnNumber) * boardDrawSize + (1 + float32(hoveredCell.ColumnNumber)) * cellDrawSize
			y := float32(hoveredCell.Board.RowNumber) * boardDrawSize + (1 + float32(hoveredCell.RowNumber)) * cellDrawSize
			vector.FillRect(screen, x, y, cellDrawSize, cellDrawSize, highlightColor, false)
		}
	}
}

func (g *Game) getBoardColors(boardReference *game.BoardReference) (ebiten.ColorScale, ebiten.ColorScale, ebiten.ColorScale) {
	if g.state.GetBoard(boardReference).Done {
		return completedForegroundColor, completedCrossColor, completedNaughtColor
	}
	if g.state.CanPlaceIn(boardReference) {
		return foregroundColor, crossColor, naughtColor
	}
	return disabledForegroundColor, disabledCrossColor, disabledNaughtColor
}

func (g *Game) drawGrid(screen *ebiten.Image, x, y float32, boardForegroundColor ebiten.ColorScale) {
	gameSize := float32(game.Size)
	cellDrawSize := g.getCellDrawSize()

	for lineNumber := range game.Size - 1 {
		lineOffset := (1 + float32(lineNumber)) * cellDrawSize
		path := vector.Path{}
		path.MoveTo(x, y + lineOffset)
		path.LineTo(x + gameSize * cellDrawSize, y + lineOffset)
		path.MoveTo(x + lineOffset, y)
		path.LineTo(x + lineOffset, y + gameSize * cellDrawSize)
		vector.StrokePath(
			screen,
			&path,
			&vector.StrokeOptions{Width: smallStrokeWidth},
			&vector.DrawPathOptions{AntiAlias: true, ColorScale: boardForegroundColor})
	}
}

func (g *Game) drawBoardSymbols(screen *ebiten.Image, board *game.Board, boardXOffset, boardYOffset float32, boardCrossColor ebiten.ColorScale, boardNaughtColor ebiten.ColorScale) {
	cellDrawSize := g.getCellDrawSize()
	cellSymbolDrawSize := cellDrawSize * (1 - cellPaddingRatio * 2)
	
	var cellRow, cellColumn byte
	for cellRow = range game.Size {
		cellY := boardYOffset + (float32(cellRow) + cellPaddingRatio) * cellDrawSize
		for cellColumn = range game.Size {
			cellX := boardXOffset + (float32(cellColumn) + cellPaddingRatio) * cellDrawSize
			cell := board.Cells.GetCell(cellRow, cellColumn)
			switch cell {
			case game.Cell_X:
				drawCross(screen, cellX, cellY, cellSymbolDrawSize, smallStrokeWidth, boardCrossColor)
			case game.Cell_O:
				drawNaught(screen, cellX, cellY, cellSymbolDrawSize, smallStrokeWidth, boardNaughtColor)
			}
		}
	}

	if board.Done {
		boardGridDrawSize := cellDrawSize * game.Size
		symbolX := boardXOffset + cellPaddingRatio * boardGridDrawSize
		symbolY := boardYOffset + cellPaddingRatio * boardGridDrawSize
		symbolDrawSize := boardGridDrawSize * (1 - cellPaddingRatio * 2)
		switch board.Owner {
		case game.Cell_X:
			drawCross(screen, symbolX, symbolY, symbolDrawSize, largeStrokeWidth, crossColor)
		case game.Cell_O:
			drawNaught(screen, symbolX, symbolY, symbolDrawSize, largeStrokeWidth, naughtColor)
		case game.Cell_None:
			drawDash(screen, symbolX, symbolY, symbolDrawSize, largeStrokeWidth, foregroundColor)
		}
	}
}

func drawCross(screen *ebiten.Image, x, y, size float32, strokeWidth float32, color ebiten.ColorScale) {
	path := vector.Path{}
	path.MoveTo(x, y)
	path.LineTo(x + size, y + size)
	path.MoveTo(x + size, y)
	path.LineTo(x, y + size)
	vector.StrokePath(
		screen,
		&path,
		&vector.StrokeOptions{Width: strokeWidth},
		&vector.DrawPathOptions{AntiAlias: true, ColorScale: color})
}

func drawNaught(screen *ebiten.Image, x, y, size float32, strokeWidth float32, color ebiten.ColorScale) {
	path := vector.Path{}
	path.Arc(x + size / 2, y + size / 2, size / 2, 0, math.Pi * 2, vector.Clockwise)
	vector.StrokePath(
		screen,
		&path,
		&vector.StrokeOptions{Width: strokeWidth},
		&vector.DrawPathOptions{AntiAlias: true, ColorScale: color})
}

func drawDash(screen *ebiten.Image, x, y, size float32, strokeWidth float32, color ebiten.ColorScale) {
	path := vector.Path{}
	path.MoveTo(x, y + size / 2)
	path.LineTo(x + size, y + size / 2)
	vector.StrokePath(
		screen,
		&path,
		&vector.StrokeOptions{Width: strokeWidth},
		&vector.DrawPathOptions{AntiAlias: true, ColorScale: color})
}

func drawExitButton(screen *ebiten.Image, x, y float32) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float32(cursorX), float32(cursorY)
	if (x <= cursorXf && cursorXf <= x + exitButtonSize) && (y <= cursorYf && cursorYf <= y + exitButtonSize) {
		vector.FillRect(screen, x, y, exitButtonSize, exitButtonSize, highlightColor, false)
	}

	drawCross(screen, x + padding, y + padding, exitCrossSize, smallStrokeWidth, foregroundColor)
}

func (g *Game) getHoveredCell() *game.Move {
	cellDrawSize := g.getCellDrawSize()
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float32(cursorX), float32(cursorY)
	hoveringCellAbsoluteRow, hoveringCellAbsoluteColumn := byte(cursorYf / cellDrawSize), byte(cursorXf / cellDrawSize)
	hoveringBoardRowNumber, hoveringBoardColumnNumber := hoveringCellAbsoluteRow / (1 + game.Size), hoveringCellAbsoluteColumn / (1 + game.Size)
	if (hoveringBoardRowNumber >= game.Size) || (hoveringBoardColumnNumber >= game.Size) { return nil }

	hoveringCellRow, hoveringCellColumn := hoveringCellAbsoluteRow % (1 + game.Size), hoveringCellAbsoluteColumn % (1 + game.Size)
	if (hoveringCellRow == 0) || (hoveringCellColumn == 0) { return nil }

	return &game.Move {
		Board: &game.BoardReference{RowNumber: hoveringBoardRowNumber, ColumnNumber: hoveringBoardColumnNumber},
		RowNumber: hoveringCellRow - 1,
		ColumnNumber: hoveringCellColumn - 1,
	}
}

func (g *Game) getCellDrawSize() float32 {
	superBoardDrawSize := float32(math.Min(float64(g.screenWidth) - float64(exitCrossSize + 2 * margin), float64(g.screenHeight)))
	cellCount := game.Size * (1 + game.Size) + 1
	cellDrawSize := superBoardDrawSize / float32(cellCount)
	return cellDrawSize
}

func (g *Game) getBoardDrawSize() float32 {
	return (float32(game.Size) + 1) * g.getCellDrawSize()
}

func getPlayers(playerSelection game.Player, aiDifficulty ai.Difficulty) (player.Player, player.Player) {
	uiPlayer := createUIPlayer()
	aiPlayer := factory.CreateAIPlayer(aiDifficulty)
	if playerSelection == game.Cell_X {
    return uiPlayer, aiPlayer
	} else {
    return aiPlayer, uiPlayer
	}
}

func (g *Game) play() {
	for gameDone := false; !gameDone; gameDone, _ = g.state.GetWinState() {
		var currentPlayer player.Player
		if g.state.GetCurrentPlayer() == game.Cell_X {
			currentPlayer = g.playerX
		} else {
			currentPlayer = g.playerO
		}

		isUIPlayer := g.state.GetCurrentPlayer() == g.playerSelection
		if isUIPlayer {
			g.activeUiPlayer = currentPlayer.(*uiPlayer)
		} else {
			g.activeUiPlayer = nil
		}

		action, move := currentPlayer.GetMove(g.state, player.NilLog)
		switch action {
		case player.Action_None:
			return
		case player.Action_Move:
			g.state.Place(move)
		default:
			panic(fmt.Sprintf("Unhandled action: %d", action))
		}
	}
}

const margin = float32(20)
const padding = float32(10)
const exitCrossSize = float32(30)
const exitButtonSize = exitCrossSize + padding * 2
const smallStrokeWidth = float32(3)
const largeStrokeWidth = float32(6)
const cellPaddingRatio = float32(0.2)

var foregroundColor = ebiten.ColorScale{}
var crossColor = ebiten.ColorScale{}
var naughtColor = ebiten.ColorScale{}

var disabledForegroundColor = ebiten.ColorScale{}
var disabledCrossColor = ebiten.ColorScale{}
var disabledNaughtColor = ebiten.ColorScale{}

var completedForegroundColor = ebiten.ColorScale{}
var completedCrossColor = ebiten.ColorScale{}
var completedNaughtColor = ebiten.ColorScale{}

var highlightColor = color.Gray{64}

func init() {
	setCrossColor(&crossColor)
	setNaughtColor(&naughtColor)
	
	disabledColorScale := ebiten.ColorScale{}
	disabledColorScale.SetR(0.5)
	disabledColorScale.SetG(0.5)
	disabledColorScale.SetB(0.5)

	disabledForegroundColor.ScaleWithColorScale(disabledColorScale)
	setCrossColor(&disabledCrossColor)
	disabledCrossColor.ScaleWithColorScale(disabledColorScale)
	setNaughtColor(&disabledNaughtColor)
	disabledNaughtColor.ScaleWithColorScale(disabledColorScale)
	
	completedColorScale := ebiten.ColorScale{}
	completedColorScale.SetR(0.2)
	completedColorScale.SetG(0.2)
	completedColorScale.SetB(0.2)

	completedForegroundColor.ScaleWithColorScale(completedColorScale)
	setCrossColor(&completedCrossColor)
	completedCrossColor.ScaleWithColorScale(completedColorScale)
	setNaughtColor(&completedNaughtColor)
	completedNaughtColor.ScaleWithColorScale(completedColorScale)
}

func setCrossColor(color *ebiten.ColorScale) {
	color.SetR(1)
	color.SetG(0)
	color.SetB(0)
}

func setNaughtColor(color *ebiten.ColorScale) {
	color.SetR(0.4)
	color.SetG(0.4)
	color.SetB(1)
}
