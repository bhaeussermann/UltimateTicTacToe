package game

import (
	"bytes"
	"image/color"

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
	drawCross(screen, float32(g.screenWidth) - crossButtonSize - margin, margin)
}

func drawCross(screen *ebiten.Image, x float32, y float32) {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorXf, cursorYf := float32(cursorX), float32(cursorY)
	if (x <= cursorXf && cursorXf <= x + crossButtonSize) && (y <= cursorYf && cursorYf <= y + crossButtonSize) {
		vector.FillRect(screen, x, y, crossButtonSize, crossButtonSize, color.Gray{128}, false)
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
