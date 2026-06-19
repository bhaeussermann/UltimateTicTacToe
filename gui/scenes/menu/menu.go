package menu

import (
	"bytes"
	"image/color"
	"math"
	"time"

	"github.com/bhaeussermann/ultimate-tic-tac-toe/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/gui/scenes"
	gamescene "github.com/bhaeussermann/ultimate-tic-tac-toe/gui/scenes/game"
	"github.com/bhaeussermann/ultimate-tic-tac-toe/player/ai"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

type TitleScreen struct {
	regularTextFaceSource, boldTextFaceSource *text.GoTextFaceSource

	backgroundPixels []byte
	screenWidth, screenHeight int
	isMouseButtonPressed bool

	selectedMenuItemIndex int

	playerSelection game.Player
	aiDifficulty ai.Difficulty
}

func NewTitleScreen() (scenes.Scene, error) {
	regularTextFaceSource, error := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if error != nil {
		return nil, error
	}

	boldTextFaceSource, error := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if error != nil {
		return nil, error
	}

	return &TitleScreen{
		regularTextFaceSource: regularTextFaceSource,
		boldTextFaceSource: boldTextFaceSource,
		backgroundPixels: []byte{},
		selectedMenuItemIndex: -1,
		playerSelection: game.Cell_X,
		aiDifficulty: ai.Difficulty_Easy,
	}, nil
}

func (t *TitleScreen) SetScreenSize(width int, height int) {
	t.screenWidth = width
	t.screenHeight = height

	if len(t.backgroundPixels) != width * height * 4 {
		t.backgroundPixels = make([]byte, width * height * 4)
	}
}

func (t *TitleScreen) Update() scenes.SceneChange {
	t.selectedMenuItemIndex = t.getSelectedMenuItemIndex()

	if !t.isMouseButtonPressed && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		switch t.selectedMenuItemIndex {
		case 0: {
			if t.playerSelection == game.Cell_X { t.playerSelection = game.Cell_O } else { t.playerSelection = game.Cell_X }
		}
		case 1: {
      switch t.aiDifficulty {
      case ai.Difficulty_Easy:
        t.aiDifficulty = ai.Difficulty_Medium
      case ai.Difficulty_Medium:
        t.aiDifficulty = ai.Difficulty_Hard
      default:
        t.aiDifficulty = ai.Difficulty_Easy
      }
		}
		case 3: {
			return scenes.SceneChange{GetNextScene: func() (scenes.Scene, error) { return gamescene.NewGame(NewTitleScreen) }}
		}
		case 4: {
			return scenes.SceneChange{Terminate: true}	
		}
		}
	}

	t.isMouseButtonPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	return scenes.SceneChange{}
}

func (t *TitleScreen) getSelectedMenuItemIndex() int {
	menuItems := t.getMenuItems()
	menuItemsHeight := 0.0
	for _, menuItem := range menuItems {
		if menuItem.itemType == menuItemType_Gap {
			menuItemsHeight += textMarginHeight
		} else {
			menuItemsHeight += menuItemTextSize + textMarginHeight * 2
		}
	}
		
	_, cursorY := ebiten.CursorPosition()
	titleHeight := titleTopMargin + titleTextSize + textMarginHeight * 2
	currentMenuItemTop := (float64(t.screenHeight) - menuItemsHeight - titleHeight) / 2 + titleHeight
	if cursorY < int(currentMenuItemTop) {
		return -1
	}

	for menuItemIndex, menuItem := range menuItems {
		if menuItem.itemType == menuItemType_Gap {
			currentMenuItemTop += textMarginHeight
		} else {
			currentMenuItemBottom := currentMenuItemTop + menuItemTextSize + textMarginHeight * 2
			if int(currentMenuItemTop) <= cursorY && cursorY <= int(currentMenuItemBottom) {
				return menuItemIndex
			}
			currentMenuItemTop = currentMenuItemBottom
		}
	}
	return -1
}

func (t *TitleScreen) Draw(screen *ebiten.Image) {
	const colorRange int = 50

	if len(t.backgroundPixels) == screen.Bounds().Dx() * screen.Bounds().Dy() * 4 {
		timeOffset := int((math.MaxInt64 - time.Now().UnixMilli()) / 50)
		for x := 0; x < t.screenWidth; x++ {
			for y := 0; y < t.screenHeight; y++ {
				colorOffset := byte(revolvingMod((max(revolvingMod(x / 2, colorRange), revolvingMod(y, colorRange)) + timeOffset), colorRange))
				pixelOffset := (y * t.screenWidth + x) * 4
				t.backgroundPixels[pixelOffset] = colorOffset
				t.backgroundPixels[pixelOffset + 1] = colorOffset
				t.backgroundPixels[pixelOffset + 2] = 50 + colorOffset
				t.backgroundPixels[pixelOffset + 3] = 255
			}
		}
		screen.WritePixels(t.backgroundPixels)
	}

	t.drawTitle(screen)
	t.drawMenuItems(screen, t.getMenuItems())
}

func revolvingMod(n int, m int) int {
	_2m := 2 * m
	return abs(n % _2m - m)
}

func abs(n int) int { if n < 0 { return -n } else { return n } }

func max(a int, b int) int { if a > b { return a } else { return b }}

func (t *TitleScreen) drawTitle(screen *ebiten.Image) {
	t.drawTextLineBackground(screen, titleTopMargin, titleTextSize + textMarginHeight * 2)

	textFace := text.GoTextFace{
		Source: t.boldTextFaceSource,
		Size: titleTextSize,
	}
	drawText(screen, float64(screen.Bounds().Dx()) / 2, titleTopMargin + textMarginHeight, "Ultimate Tic Tac Toe", &textFace, text.AlignCenter)
}

func (t *TitleScreen) getMenuItems() []menuItem {
	var playerSelection string
	if t.playerSelection == game.Cell_X { playerSelection = "X" } else { playerSelection = "O" }
	var difficultySelection string
	switch t.aiDifficulty {
	case ai.Difficulty_Easy: difficultySelection = "Easy"
	case ai.Difficulty_Medium: difficultySelection = "Medium"
	default: difficultySelection = "Hard"
	}

	return []menuItem {
		*createSelectionItem("Player selection: ", playerSelection),
		*createSelectionItem("AI difficulty: ", difficultySelection),
		*createGapItem(),
		*createActionItem("Start game"),
		*createActionItem("Exit"),
	}
}

func (t *TitleScreen) drawMenuItems(screen *ebiten.Image, menuItems []menuItem) {
	menuItemLabelTextFace := text.GoTextFace{
		Source: t.regularTextFaceSource,
		Size: menuItemTextSize,
	}
	menuItemSelectionTextFace := text.GoTextFace{
		Source: t.boldTextFaceSource,
		Size: menuItemTextSize,
	}

	menuItemsTextRuns := make([]menuItemTextRuns, len(menuItems))
	for menuItemIndex, menuItem := range menuItems {
		currentMenuItemTextRuns := menuItemTextRuns {}
		if menuItem.itemType == menuItemType_Selection {
			currentMenuItemTextRuns.label = &textRun{menuItem.label, &menuItemLabelTextFace}
		}
		if (menuItem.itemType == menuItemType_Selection) || (menuItem.itemType == menuItemType_Action) {
			currentMenuItemTextRuns.selection = &textRun{menuItem.selection, &menuItemSelectionTextFace}
		}
		menuItemsTextRuns[menuItemIndex] = currentMenuItemTextRuns
	}

	maximumOptionLabelWidth := 0.0
	maximumOptionSelectionWidth := 0.0
	optionItemsHeight := 0.0
	for _, menuItemTextRun := range menuItemsTextRuns {
		if menuItemTextRun.selection == nil {
			optionItemsHeight += textMarginHeight
			continue
		}
		optionItemsHeight += menuItemTextSize + textMarginHeight * 2
		if menuItemTextRun.label != nil {
			maximumOptionLabelWidth = math.Max(maximumOptionLabelWidth, menuItemTextRun.label.getWidth())
			maximumOptionSelectionWidth = math.Max(maximumOptionSelectionWidth, menuItemTextRun.selection.getWidth())
		}
	}

	titleHeight := titleTopMargin + titleTextSize + textMarginHeight * 2
	currentItemTop := (float64(screen.Bounds().Dy()) - titleHeight - optionItemsHeight) / 2 + titleHeight
	boundaryLineX := (float64(screen.Bounds().Dx()) - (maximumOptionLabelWidth + maximumOptionSelectionWidth)) / 2 + maximumOptionLabelWidth
	for menuItemIndex, menuItemTextRuns := range menuItemsTextRuns {
		if t.selectedMenuItemIndex == menuItemIndex {
			t.drawTextLineBackground(screen, currentItemTop, menuItemTextSize + textMarginHeight * 2)
		}
		if menuItemTextRuns.selection == nil {
			currentItemTop += textMarginHeight
			continue
		}
		if menuItemTextRuns.label == nil {
			drawText(screen, float64(screen.Bounds().Dx() / 2), currentItemTop + textMarginHeight, menuItemTextRuns.selection.text, menuItemTextRuns.selection.textFace, text.AlignCenter)
		} else {
			drawText(screen, boundaryLineX, currentItemTop + textMarginHeight, menuItemTextRuns.label.text, menuItemTextRuns.label.textFace, text.AlignEnd)
			drawText(screen, boundaryLineX, currentItemTop + textMarginHeight, menuItemTextRuns.selection.text, menuItemTextRuns.selection.textFace, text.AlignStart)
		}
		currentItemTop += menuItemTextSize + textMarginHeight * 2
	}
}

func drawText(screen *ebiten.Image, x float64, y float64, textString string, textFace *text.GoTextFace, alignment text.Align) {
	drawGeom := ebiten.GeoM{}
	drawGeom.Translate(x, y)
	drawOptions := &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: drawGeom,
		},
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: alignment,
		},
	}
	text.Draw(screen, textString, textFace, drawOptions)
}

func (textRun *textRun) getWidth() float64 {
	textWidth, _ := text.Measure(textRun.text, textRun.textFace, 0)
	return textWidth
}

func (t *TitleScreen) drawTextLineBackground(screen *ebiten.Image, top float64, height float64) {
	vector.FillRect(screen, 0, float32(top), float32(t.screenWidth), float32(height), color.Alpha{80}, false)
}

type menuItem struct {
	label string
	selection string
	itemType byte
}

const (
	menuItemType_Selection = iota
	menuItemType_Action
	menuItemType_Gap
)

func createSelectionItem(label string, selection string) *menuItem {
	return &menuItem{label, selection, menuItemType_Selection}
}

func createActionItem(text string) *menuItem {
	return &menuItem{selection: text, itemType: menuItemType_Action}
}

func createGapItem() *menuItem {
	return &menuItem{itemType: menuItemType_Gap}
}

type textRun struct {
	text string
	textFace *text.GoTextFace
}

type menuItemTextRuns struct {
	label *textRun
	selection *textRun
}

const titleTopMargin = float64(20)
const titleTextSize = float64(42)
const menuItemTextSize = float64(36)
const textMarginHeight = float64(20)
