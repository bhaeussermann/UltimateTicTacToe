package game

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	gui "github.com/bhaeussermann/ultimate-tic-tac-toe/gui"
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
	pointerState *gui.PointerState
	exitToScene scenes.GetNextScene

	playerSelection game.Player
	playerX, playerO player.Player
	state *game.State
	logText string
	isAiThinking bool
	activeUiPlayer *uiPlayer
}

func NewGame(playerSelection game.Player, aiDifficulty ai.Difficulty, exitToScene scenes.GetNextScene) (scenes.Scene, error) {
	textFaceSource, error := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if error != nil {
		return nil, error
	}

	playerX, playerO := getPlayers(playerSelection, aiDifficulty)
	game := &Game{
		textFaceSource: textFaceSource,
		pointerState: &gui.PointerState{},
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
	cursorX, cursorY, didSelect := g.pointerState.GetSelectState()

	if didSelect {
		cursorXf, cursorYf := float32(cursorX), float32(cursorY)
		exitClicked := (float32(g.screenWidth) - exitButtonSize - margin <= cursorXf && cursorXf <= float32(g.screenWidth) - margin) && (margin <= cursorYf && cursorYf <= exitButtonSize + margin)
		if exitClicked {
			if g.activeUiPlayer != nil {
				g.activeUiPlayer.unblock()
			}
			return scenes.SceneChange{GetNextScene: g.exitToScene}
		}

		if g.activeUiPlayer != nil {
			selectedCell := g.getCellAtLocation(cursorX, cursorY)
			if (selectedCell != nil) && g.state.CanPlaceIn(selectedCell.Board) && g.state.CanPlace(selectedCell) {
				g.activeUiPlayer.setMove(selectedCell)
			}
		}
	}

	return scenes.SceneChange{}
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawExitButton(screen, float32(g.screenWidth) - exitButtonSize - margin, margin)
	g.drawSuperBoard(screen)
	g.writeStatusMessage(screen)
	g.writeMessageLog(screen)
	g.drawBusyIndicator(screen)
}

func drawExitButton(screen *ebiten.Image, x, y float32) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float32(cursorX), float32(cursorY)
	if (x <= cursorXf && cursorXf <= x + exitButtonSize) && (y <= cursorYf && cursorYf <= y + exitButtonSize) {
		vector.FillRect(screen, x, y, exitButtonSize, exitButtonSize, highlightColor, false)
	}

	drawCross(screen, x + padding, y + padding, exitCrossSize, smallStrokeWidth, foregroundColor)
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
		hoveredCell := g.getCellAtLocation(g.pointerState.GetHoverState())
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

func (g *Game) writeStatusMessage(screen *ebiten.Image) {
	done, winner := g.state.GetWinState()
	if !done { return }

	switch winner {
	case game.Cell_X:
		g.writeStatusMessageText(screen, "Cross is the winner!")
	case game.Cell_O:
		g.writeStatusMessageText(screen, "Naughts is the winner!")
	default:
		g.writeStatusMessageText(screen, "It's a tie.")
	}
}

func (g *Game) writeStatusMessageText(screen *ebiten.Image, messageText string) {
	var textYRatio float32
	hasLogText := len(g.logText) != 0
	if hasLogText {
		textYRatio = 1.0 / 3.0
	} else {
		textYRatio = 0.5
	}
	gameFieldSize, textLocation := g.getGameFieldSize()
	gameFieldSize -= g.getCellDrawSize()
	var textX, textY float32
	if textLocation == TextLocation_Side {
		textX = gameFieldSize + (float32(g.screenWidth) - gameFieldSize) / 2
		textY = float32(g.screenHeight) * textYRatio
	} else {
		textX = float32(g.screenWidth) / 2
		textY = gameFieldSize + (float32(g.screenHeight) - gameFieldSize) * textYRatio
	}

	drawGeom := ebiten.GeoM{}
	drawGeom.Translate(float64(textX), float64(textY))
	drawOptions := &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: drawGeom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
		},
	}
	textFace := text.GoTextFace{
		Source: g.textFaceSource,
		Size: statusMessageTextSize,
	}
	text.Draw(screen, messageText, &textFace, drawOptions)
}

func (g *Game) writeMessageLog(screen *ebiten.Image) {
	if len(g.logText) == 0 { return }

	textX, textY := g.getStatusDrawLocation()
	drawGeom := ebiten.GeoM{}
	drawGeom.Translate(float64(textX), float64(textY))
	drawOptions := &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: drawGeom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
			LineSpacing: logTextSize,
		},
	}
	textFace := text.GoTextFace{
		Source: g.textFaceSource,
		Size: logTextSize,
	}
	text.Draw(screen, g.logText, &textFace, drawOptions)
}

func (g *Game) drawBusyIndicator(screen *ebiten.Image) {
	if !g.isAiThinking { return }

	x, y := g.getStatusDrawLocation()
	animationFrame := time.Now().UnixMilli() / 100 % 10
	drawBusyIndicatorCircle(screen, x, y, animationFrame)
	drawBusyIndicatorCircle(screen, x, y, animationFrame + 10)
}

func drawBusyIndicatorCircle(screen *ebiten.Image, x float32, y float32, animationFrame int64) {
	path := vector.Path{}
	path.Arc(x, y, float32(animationFrame), 0, math.Pi * 2, vector.Clockwise)
	color := ebiten.ColorScale{}
	color.ScaleAlpha(1 - float32(animationFrame) / 20)
	vector.StrokePath(screen, &path, &vector.StrokeOptions{Width: 1}, &vector.DrawPathOptions{AntiAlias: true, ColorScale: color})
}

func (g *Game) getStatusDrawLocation() (float32, float32) {
	var textYRatio float32
	isStatusMessageVisible, _ := g.state.GetWinState()
	if isStatusMessageVisible {
		textYRatio = 2.0 / 3.0
	} else {
		textYRatio = 0.5
	}
	gameFieldSize, textLocation := g.getGameFieldSize()
	gameFieldSize -= g.getCellDrawSize()
	var textX, textY float32
	if textLocation == TextLocation_Side {
		textX = gameFieldSize + (float32(g.screenWidth) - gameFieldSize) / 2
		textY = float32(g.screenHeight) * textYRatio
	} else {
		textX = float32(g.screenWidth) / 2
		textY = gameFieldSize + (float32(g.screenHeight) - gameFieldSize) * textYRatio
	}
	return textX, textY
}

func (g *Game) getCellAtLocation(locationX, locationY int) *game.Move {
	cellDrawSize := g.getCellDrawSize()
	locationXf, locationYf := float32(locationX), float32(locationY)
	hoveringCellAbsoluteRow, hoveringCellAbsoluteColumn := byte(locationYf / cellDrawSize), byte(locationXf / cellDrawSize)
	hoveringBoardRowNumber, hoveringBoardColumnNumber := hoveringCellAbsoluteRow / (1 + game.Size), hoveringCellAbsoluteColumn / (1 + game.Size)
	if (hoveringBoardRowNumber >= game.Size) || (hoveringBoardColumnNumber >= game.Size) { return nil }

	hoveringCellRow, hoveringCellColumn := hoveringCellAbsoluteRow % (1 + game.Size), hoveringCellAbsoluteColumn % (1 + game.Size)
	if (hoveringCellRow == 0) || (hoveringCellColumn == 0) { return nil }

	return &game.Move{
		Board: &game.BoardReference{RowNumber: hoveringBoardRowNumber, ColumnNumber: hoveringBoardColumnNumber},
		RowNumber: hoveringCellRow - 1,
		ColumnNumber: hoveringCellColumn - 1,
	}
}

func (g *Game) getCellDrawSize() float32 {
	superBoardDrawSize, _ := g.getGameFieldSize()
	cellCount := game.Size * (1 + game.Size) + 1
	cellDrawSize := float32(superBoardDrawSize) / float32(cellCount)
	return cellDrawSize
}

func (g *Game) getGameFieldSize() (float32, byte) {
	if g.screenWidth > g.screenHeight {
		if float32(g.screenWidth) > float32(g.screenHeight) + minimumTextAreaWidth {
			return float32(g.screenHeight), TextLocation_Side
		}
		return float32(g.screenWidth) - minimumTextAreaWidth, TextLocation_Side
	}

	availableWidth := float32(g.screenWidth) - exitCrossSize - 2 * margin
	availableHeight := float32(g.screenHeight) - minimumTextAreaHeight
	return minFloat32(availableWidth, availableHeight), TextLocation_Under
}

func minFloat32(x, y float32) float32 {
	if x < y {
		return x
	} else {
		return y
	}
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
			g.isAiThinking = true
		}

		messageLog := player.CreateLog()
		action, move := currentPlayer.GetMove(g.state, messageLog)
		g.updateLogText(messageLog)
		g.isAiThinking = false

		switch action {
		case player.Action_None:
			return
		case player.Action_Move:
			g.state.Place(move)
		default:
			panic(fmt.Sprintf("Unhandled action: %d", action))
		}
	}
	g.activeUiPlayer = nil
}

func (g *Game) updateLogText(messageLog *player.MessageLog) {
	g.logText = ""
	messages := messageLog.GetMessages()
	if len(messages) == 0 {
		return
	}

	for index, message := range messages {
		if index != 0 {
			g.logText += "\n"
		}
		g.logText += message
	}
}

const margin = float32(20)
const padding = float32(10)
const exitCrossSize = float32(30)
const exitButtonSize = exitCrossSize + padding * 2
const smallStrokeWidth = float32(3)
const largeStrokeWidth = float32(6)
const cellPaddingRatio = float32(0.2)

const minimumTextAreaWidth = float32(350)
const minimumTextAreaHeight = float32(180)
const statusMessageTextSize = float64(36)
const logTextSize = float64(18)

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

const (
	TextLocation_Side = iota
	TextLocation_Under
)

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
