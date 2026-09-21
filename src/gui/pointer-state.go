package gui

import "github.com/hajimehoshi/ebiten/v2"

type PointerState struct {
	isMobileFlagInitialized, isMobileFlag bool
	isSelecting bool
}

func (s *PointerState) GetSelectState() (int, int, bool) {
	if s.isMobile() {
		ebiten.AppendTouchIDs(nil)
		touchIDs := ebiten.AppendTouchIDs(nil)
		if len(touchIDs) == 0 {
			s.isSelecting = false
			return 0, 0, false
		}
		cursorX, cursorY := ebiten.TouchPosition(touchIDs[0])
		didSelectionStart := !s.isSelecting
		s.isSelecting = true
		return cursorX, cursorY, didSelectionStart
	}

	isSelecting := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	didSelectionStart := !s.isSelecting && isSelecting
	s.isSelecting = isSelecting
	cursorX, cursorY := ebiten.CursorPosition()
	return cursorX, cursorY, didSelectionStart
}

func (s *PointerState) GetHoverState() (int, int) {
	return ebiten.CursorPosition()
}

func (s *PointerState) isMobile() bool {
	if !s.isMobileFlagInitialized {
		cursorX, cursorY := ebiten.CursorPosition()
		s.isMobileFlag = (cursorX == 0) && (cursorY == 0)
	}
	return s.isMobileFlag
}
