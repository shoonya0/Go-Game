package core

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GameVersion is shown on the How to Play board.
const GameVersion = "v1.0"

// How to Play (opened from Settings): pages of controls and rules drawn on the
// menu board, each key as a small plank. The numbers come from the game's
// constants, so retuning an attack updates its line.

// helpRow is one line of a page: the keys to press and what they do. Combo
// keys are pressed together (drawn with a +), others in turn. A row without
// keys is a note across the whole page.
type helpRow struct {
	Keys  []string
	Combo bool
	Text  string
}

// helpPage is one page of How to Play.
type helpPage struct {
	Name string
	Rows []helpRow
}

// bars is hp in enemy health bars, the unit the enemies' health shows.
func bars(hp int) string {
	n := float64(hp) / EnemyHealthPerBar
	if n <= 1 {
		return fmt.Sprintf("%g bar", n)
	}
	return fmt.Sprintf("%g bars", n)
}

// stamina is a power cost in stamina parts (HUD segments).
func stamina(cost float64) string {
	return fmt.Sprintf("%g ST", cost/(MaxPlayerPower/hudSegments))
}

// hearts is health in hearts.
func hearts(hp int) string {
	n := float64(hp) / (MaxPlayerHealth / hudSegments)
	if n == 1 {
		return "1 heart"
	}
	return fmt.Sprintf("%g hearts", n)
}

var helpPages = []helpPage{
	{Name: "MOVES", Rows: []helpRow{
		{Keys: []string{"A", "D"}, Text: "Walk left and right"},
		{Keys: []string{"W", "S"}, Text: "Pick in menus"},
		{Keys: []string{"Space"}, Text: "Jump (" + stamina(JumpPowerCost) + ")"},
		{Keys: []string{"Shift"}, Text: "Hold while walking to run (uses ST)"},
		{Keys: []string{"E"}, Text: "At a portal: save. Hold E: travel"},
		{Keys: []string{"F"}, Text: "Tease the enemies"},
		{Keys: []string{"Esc"}, Text: "Pause the game"},
		{Keys: []string{"F3"}, Text: "Show hit boxes"},
	}},
	{Name: "POWERS", Rows: []helpRow{
		{Text: "Damage in enemy health bars. ST: stamina parts"},
		{Keys: []string{"J"}, Text: "Quick slash: " + bars(WeakAttackDamage) + ", " + stamina(AttackPowerCost)},
		{Keys: []string{"I"}, Text: "High slash: " + bars(StrongAttackDamage) + ", " + stamina(AttackPowerCost)},
		{Keys: []string{"Space", "J"}, Combo: true, Text: "Air slash: " + bars(WeakAttackDamage) + ", " + stamina(AttackPowerCost)},
		{Keys: []string{"Space", "I"}, Combo: true, Text: "Air smash, front and back: " + bars(StrongAirAttackDamage) + ", " + stamina(AttackPowerCost)},
		{Keys: []string{"K", "K"}, Text: fmt.Sprintf("Twin crescents: %s + %s, free", bars(SpecialAttack1Damage1), bars(SpecialAttack1Damage2))},
		{Keys: []string{"L", "L"}, Text: fmt.Sprintf("Blue disc: %s + %s, free", bars(SpecialAttack2Damage1), bars(SpecialAttack2Damage2))},
		{Keys: []string{"U"}, Text: "Ground wave: " + bars(GroundWaveDamage) + " to foes on the ground, " + stamina(SpecialPowerCost)},
		{Keys: []string{"O", "O"}, Text: fmt.Sprintf("Burst, then thrust: %s + %s, %g + %s", bars(SpecialAttack4Damage1), bars(SpecialAttack4Damage2),
			SpecialAttack4PowerCost1/(MaxPlayerPower/hudSegments), stamina(SpecialAttack4PowerCost2))},
		{Text: "K, L, O: press again for the second hit. Air: press while falling"},
	}},
	{Name: "RULES", Rows: []helpRow{
		{Text: fmt.Sprintf("You have %d hearts. A hit of %s or more knocks you down", hudSegments, hearts(HeavyHitDamage))},
		{Text: fmt.Sprintf("Stamina refills by itself. At 0 you rest until %s refill", stamina(RunRecoverPower))},
		{Text: "Water slows you, stops your run and slowly drains health"},
		{Text: "Larva burns " + hearts(LarvaDamage) + " each time you step in"},
		{Text: fmt.Sprintf("Falls of more than %d tiles hurt; into water is safe", FallSafeTiles)},
		{Text: "Save at portals. Use every portal to open the next level"},
		{Text: "Enemies heal over time: finish them fast"},
		{Text: "Beaten enemies return when you enter a level or travel"},
	}},
}

const (
	helpPanelW    = 1220
	helpPanelH    = 720
	helpTitleSize = 36
	helpPageSize  = 18 // the page name and number
	helpTextSize  = 15
	helpKeySize   = 14
	helpKeyH      = 34 // px of a key plank
	helpKeyPadX   = 14 // px around a key's letters
	helpKeyGap    = 8  // px between keys pressed in turn
	helpPlusW     = 30 // px of the + between combo keys
	helpRowH      = 46
	helpRowsTop   = 170 // px from the board's top to the first row's center
	helpTextX     = 280 // px from the board's left to a row's text
	helpPadX      = 80  // px from the board's sides to the keys and notes
	helpBackW     = 260
	helpBackH     = 56
	helpBottom    = 84 // px from the board's bottom to the Back button's center
	helpArrowGap  = 30 // px between Back and an arrow
)

var helpTextColor = color.RGBA{0xf6, 0xe7, 0xc8, 0xff} // cream on the wood

// drawHelpText draws s from x, centered on cy, over a dark shadow so it reads
// on the board's streaks.
func drawHelpText(screen *ebiten.Image, s string, size, x, cy float64) {
	DrawText(screen, s, size, x+2, cy+2, menuTitleShadow)
	DrawText(screen, s, size, x, cy, helpTextColor)
}

// HelpScreen is How to Play: one page at a time, turned with left/right or the
// arrow buttons, closed with Back or Esc.
type HelpScreen struct {
	Page             int
	screenW, screenH int
}

// NewHelpScreen lays How to Play out for a w×h screen.
func NewHelpScreen(w, h int) *HelpScreen { return &HelpScreen{screenW: w, screenH: h} }

// PanelRect is the board, centered on the screen.
func (h *HelpScreen) PanelRect() image.Rectangle {
	x, y := (h.screenW-helpPanelW)/2, (h.screenH-helpPanelH)/2
	return image.Rect(x, y, x+helpPanelW, y+helpPanelH)
}

// BackRect is the Back button, centered at the board's bottom.
func (h *HelpScreen) BackRect() image.Rectangle {
	p := h.PanelRect()
	cx, cy := (p.Min.X+p.Max.X)/2, p.Max.Y-helpBottom
	return image.Rect(cx-helpBackW/2, cy-helpBackH/2, cx+helpBackW/2, cy+helpBackH/2)
}

// ArrowRect is the page arrow on side dir (-1 left, 1 right) of Back.
func (h *HelpScreen) ArrowRect(dir int) image.Rectangle {
	b := h.BackRect()
	x := b.Max.X + helpArrowGap
	if dir < 0 {
		x = b.Min.X - helpArrowGap - helpBackH
	}
	return image.Rect(x, b.Min.Y, x+helpBackH, b.Max.Y)
}

// keyWidth is how wide key k's plank is drawn.
func keyWidth(k string) int {
	return max(helpKeyH, int(textWidth(k, helpKeySize))+2*helpKeyPadX)
}

// keyRects places row r's keys on row j, left to right.
func (h *HelpScreen) keyRects(j int, r helpRow) []image.Rectangle {
	p := h.PanelRect()
	x, cy := p.Min.X+helpPadX, h.rowY(j)
	rects := make([]image.Rectangle, len(r.Keys))
	for i, k := range r.Keys {
		if i > 0 {
			x += helpKeyGap
			if r.Combo {
				x += helpPlusW - helpKeyGap
			}
		}
		w := keyWidth(k)
		rects[i] = image.Rect(x, cy-helpKeyH/2, x+w, cy+helpKeyH/2)
		x += w
	}
	return rects
}

// rowY is the screen y row j is centered on.
func (h *HelpScreen) rowY(j int) int { return h.PanelRect().Min.Y + helpRowsTop + j*helpRowH }

// textX is where row r's text starts: after the key column, or at the left
// edge for a note.
func (h *HelpScreen) textX(r helpRow) int {
	if len(r.Keys) == 0 {
		return h.PanelRect().Min.X + helpPadX
	}
	return h.PanelRect().Min.X + helpTextX
}

// rowRect is the box row j's keys and text cover.
func (h *HelpScreen) rowRect(j int, r helpRow) image.Rectangle {
	cy := h.rowY(j)
	x := h.textX(r)
	box := image.Rect(x, cy-helpKeyH/2, x+int(textWidth(r.Text, helpTextSize)), cy+helpKeyH/2)
	for _, k := range h.keyRects(j, r) {
		box = box.Union(k)
	}
	return box
}

// turn moves to the page dir (±1) away, stopping at the first and the last.
func (h *HelpScreen) turn(dir int) { h.Page = min(len(helpPages)-1, max(0, h.Page+dir)) }

// Update applies one tick of input and reports whether How to Play closes.
func (h *HelpScreen) Update(in *InputState) (back bool) {
	if at, pressed := in.PressAt(); pressed {
		pt := image.Pt(int(at.X), int(at.Y))
		switch {
		case pt.In(h.BackRect()):
			return true
		case pt.In(h.ArrowRect(-1)):
			h.turn(-1)
		case pt.In(h.ArrowRect(1)):
			h.turn(1)
		}
		return false
	}
	switch {
	case in.Back || in.Confirm || in.Interact:
		return true
	case in.LeftJustPressed:
		h.turn(-1)
	case in.RightJustPressed:
		h.turn(1)
	}
	return false
}

// Draw draws the current page on the board, with the page arrows (dark at
// the first and last page) and Back.
func (h *HelpScreen) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(h.screenW), float32(h.screenH), menuShade, false)
	p := h.PanelRect()
	drawSliced(screen, uiBoardRect, uiBoardInsets, p, menuBoardScale, nil, 1)
	cx := float64(p.Min.X+p.Max.X) / 2
	ty := float64(p.Min.Y) + 78
	DrawTextCentered(screen, "HOW TO PLAY", helpTitleSize, cx+3, ty+3, menuTitleShadow)
	DrawTextCentered(screen, "HOW TO PLAY", helpTitleSize, cx, ty, menuTitleColor)
	page := helpPages[h.Page]
	pageName := fmt.Sprintf("%s  %d/%d", page.Name, h.Page+1, len(helpPages))
	DrawTextCentered(screen, pageName, helpPageSize, cx+2, ty+46, menuTitleShadow)
	DrawTextCentered(screen, pageName, helpPageSize, cx, ty+44, helpTextColor)

	for j, r := range page.Rows {
		keys := h.keyRects(j, r)
		for i, k := range keys {
			drawPlank(screen, uiPlankRect, uiPlankInsets, k, nil)
			DrawTextCentered(screen, r.Keys[i], helpKeySize, float64(k.Min.X+k.Max.X)/2,
				float64(k.Min.Y)+menuTextRise*float64(k.Dy()), menuItemColor)
			if r.Combo && i > 0 {
				DrawTextCentered(screen, "+", helpTextSize, float64(k.Min.X)-helpPlusW/2, float64(h.rowY(j)), helpTextColor)
			}
		}
		drawHelpText(screen, r.Text, helpTextSize, float64(h.textX(r)), float64(h.rowY(j)))
	}

	for _, a := range []struct {
		dir  int
		icon image.Rectangle
		end  bool
	}{{-1, uiLeftRect, h.Page == 0}, {1, uiRightRect, h.Page == len(helpPages)-1}} {
		var tint color.Color
		if a.end {
			tint = menuDisabledTint
		}
		drawRotated(screen, uiSprite(a.icon), h.ArrowRect(a.dir), 0, tint, 1)
	}
	b := h.BackRect()
	drawPlank(screen, uiPlankArrowRect, uiPlankArrowInsets, b, nil)
	DrawTextCentered(screen, "Back", menuItemSize, float64(b.Min.X+b.Max.X)/2,
		float64(b.Min.Y)+menuTextRise*float64(b.Dy()), menuItemColor)
	drawHelpText(screen, GameVersion, helpKeySize, float64(p.Max.X-helpPadX)-textWidth(GameVersion, helpKeySize),
		float64(b.Min.Y+b.Max.Y)/2)
}
