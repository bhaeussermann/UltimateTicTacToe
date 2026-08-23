package game

import (
	"bytes"
	"image/color"
	"math"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/gui/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	textFaceSource *text.GoTextFaceSource
	screenWidth, screenHeight int
	exitToScene scenes.GetNextScene
}

func NewGame(exitToScene scenes.GetNextScene) (scenes.Scene, error) {
	textFaceSource, error := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if error != nil {
		return nil, error
	}

	return &Game{
		textFaceSource: textFaceSource,
		exitToScene: exitToScene,
	}, nil
}

func (g *Game) SetScreenSize(width int, height int) {
	g.screenWidth = width
	g.screenHeight = height
}

func (g *Game) Update() scenes.SceneChange {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		cursorXf, cursorYf := float32(cursorX), float32(cursorY)
		if (float32(g.screenWidth) - crossButtonSize - margin <= cursorXf && cursorXf <= float32(g.screenWidth) - margin) && (margin <= cursorYf && cursorYf <= crossButtonSize + margin) {
			return scenes.SceneChange{GetNextScene: g.exitToScene }
		}
	}

	return scenes.SceneChange{}
}


func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
	drawCross(screen, float32(g.screenWidth) - crossButtonSize - margin, margin)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	gameSize := float32(game.Size)
	cellDrawSize := g.getCellDrawSize()
	boardDrawSize := (gameSize + 1) * cellDrawSize

	hoveredCell := g.getHoveredCell()
	if hoveredCell != nil {
		x := float32(hoveredCell.Board.ColumnNumber) * boardDrawSize + (1 + float32(hoveredCell.ColumnNumber)) * cellDrawSize
		y := float32(hoveredCell.Board.RowNumber) * boardDrawSize + (1 + float32(hoveredCell.RowNumber)) * cellDrawSize
		vector.FillRect(screen, x, y, cellDrawSize, cellDrawSize, highlightColor, false)
	}

	for boardRow := range game.Size {
		boardYOffset := float32(boardRow) * boardDrawSize + cellDrawSize
		for boardColumn := range game.Size {
			boardXOffset := float32(boardColumn) * boardDrawSize + cellDrawSize
			for lineNumber := range game.Size - 1 {
				lineOffset := (1 + float32(lineNumber)) * cellDrawSize
				path := vector.Path{}
				path.MoveTo(boardXOffset, boardYOffset + lineOffset)
				path.LineTo(boardXOffset + gameSize * cellDrawSize, boardYOffset + lineOffset)
				path.MoveTo(boardXOffset + lineOffset, boardYOffset)
				path.LineTo(boardXOffset + lineOffset, boardYOffset + gameSize * cellDrawSize)
				vector.StrokePath(
					screen,
					&path,
					&vector.StrokeOptions{Width: strokeWidth},
					&vector.DrawPathOptions{AntiAlias: true})
			}
		}
	}
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
	boardDrawSize := float32(math.Min(float64(g.screenWidth) - float64(crossSize + 2 * margin), float64(g.screenHeight)))
	cellCount := game.Size * (1 + game.Size) + 1
	cellDrawSize := boardDrawSize / float32(cellCount)
	return cellDrawSize
}

func drawCross(screen *ebiten.Image, x float32, y float32) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float32(cursorX), float32(cursorY)
	if (x <= cursorXf && cursorXf <= x + crossButtonSize) && (y <= cursorYf && cursorYf <= y + crossButtonSize) {
		vector.FillRect(screen, x, y, crossButtonSize, crossButtonSize, highlightColor, false)
	}

	foregroundColorScale := ebiten.ColorScale{}
	foregroundColorScale.Scale(1, 1, 1, 1)
	path := vector.Path{}
	path.MoveTo(x + padding, y + padding)
	path.LineTo(x + padding + crossSize, y + padding + crossSize)
	path.MoveTo(x + padding + crossSize, y + padding)
	path.LineTo(x + padding, y + padding + crossSize)
	vector.StrokePath(
		screen,
		&path,
		&vector.StrokeOptions{Width: strokeWidth},
		&vector.DrawPathOptions{AntiAlias: true, ColorScale: foregroundColorScale})
}

const margin = float32(20)
const padding = float32(10)
const crossSize = float32(30)
const crossButtonSize = crossSize + padding * 2
const strokeWidth = float32(3)

var highlightColor = color.Gray{64}
