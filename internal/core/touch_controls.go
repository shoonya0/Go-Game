package core

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// UIButtonsSheet holds the wooden UI buttons used by the on-screen touch
// controls.
const UIButtonsSheet = "../assets/UIButtons.png"

// TouchButton identifies one on-screen button of the touch controls.
type TouchButton int

const (
	TouchUp       TouchButton = iota // W
	TouchLeft                        // A
	TouchDown                        // S
	TouchRight                       // D
	TouchSpecial3                    // U
	TouchStrong                      // I
	TouchSpecial4                    // O
	TouchWeak                        // J
	TouchSpecial1                    // K
	TouchSpecial2                    // L
	TouchSprint                      // Left Shift
	TouchMenu
	touchButtonCount
)

// Touch layout in logical screen pixels: movement buttons bottom-left, attack
// buttons bottom-right, the menu button top-left. The clusters are kept small
// and in the outer thirds of the screen, where the camera never puts the player
// (see UpdateCamera), and idle buttons are see-through so they hide little.
const (
	touchButtonSize = 80
	touchButtonGap  = 8
	touchEdge       = 20 // distance of the button clusters from the screen edges
	touchMenuSize   = 64
	touchIconInset  = 15  // percent of the button cut from each side of its symbol, so it draws at 70%
	touchIdleAlpha  = 0.5 // opacity of buttons that are not pressed

	// TouchHUDLeft is where the HUD starts when the menu button takes the
	// top-left corner.
	TouchHUDLeft = hudMargin + touchMenuSize + hudMargin
)

// Source rects inside UIButtonsSheet, measured from the non-transparent pixels.
var (
	uiMenuRect       = image.Rect(125, 1430, 329, 1641)   // square button with the menu icon
	uiRoundUpRect    = image.Rect(2052, 1732, 2268, 1946) // round gold button
	uiRoundDownRect  = image.Rect(1811, 1732, 2027, 1946) // round dark button, used while pressed
	touchIconColor   = color.RGBA{0x5a, 0x24, 0x10, 0xff}
	touchPressedTint = color.RGBA{0xff, 0xd2, 0x6a, 0xff}
)

// TouchControls is the on-screen gamepad shown on touch screens. It turns the
// touch points of a frame into button states and merges them into InputState,
// so the player code never knows whether input came from keys or fingers.
type TouchControls struct {
	Enabled bool // shown and active; set on touch devices or after the first touch

	rects      [touchButtonCount]image.Rectangle
	held, prev [touchButtonCount]bool

	sprites *touchSprites
}

type touchSprites struct {
	menu, up, down *ebiten.Image
	icons          [touchButtonCount]*ebiten.Image // symbols from touch_icons.go; nil for the menu
}

// NewTouchControls lays the buttons out for a logical screen of w×h pixels.
func NewTouchControls(w, h int) *TouchControls {
	tc := &TouchControls{}
	s, g := touchButtonSize, touchButtonGap
	sq := func(x, y, size int) image.Rectangle { return image.Rect(x, y, x+size, y+size) }

	// Movement: a diamond whose middle is empty, W on top and S below it.
	step := s + g/2 // a full button apart, so the diagonal neighbours don't overlap
	cx := touchEdge + step + s/2
	cy := h - touchEdge - step - s/2
	tc.rects[TouchUp] = sq(cx-s/2, cy-step-s/2, s)
	tc.rects[TouchDown] = sq(cx-s/2, cy+step-s/2, s)
	tc.rects[TouchLeft] = sq(cx-step-s/2, cy-s/2, s)
	tc.rects[TouchRight] = sq(cx+step-s/2, cy-s/2, s)

	// Attacks: two rows laid out like the keyboard, U I O over J K L, with
	// sprint centred above them.
	bottom := h - touchEdge - s
	top := bottom - s - g
	for col, b := range [3][2]TouchButton{{TouchSpecial3, TouchWeak}, {TouchStrong, TouchSpecial1}, {TouchSpecial4, TouchSpecial2}} {
		x := w - touchEdge - (3-col)*s - (2-col)*g
		tc.rects[b[0]] = sq(x, top, s)
		tc.rects[b[1]] = sq(x, bottom, s)
	}
	tc.rects[TouchSprint] = sq(tc.rects[TouchStrong].Min.X, top-g-s, s)

	tc.rects[TouchMenu] = sq(hudMargin, hudMargin, touchMenuSize)
	return tc
}

// Rect returns the screen area of button b.
func (tc *TouchControls) Rect(b TouchButton) image.Rectangle { return tc.rects[b] }

// Update records which buttons the given touch points are on. Any touch turns
// the controls on, so touch screens the platform check missed still get them.
func (tc *TouchControls) Update(points []image.Point) {
	if len(points) > 0 {
		tc.Enabled = true
	}
	tc.prev = tc.held
	tc.held = [touchButtonCount]bool{}
	for _, p := range points {
		for b, r := range tc.rects {
			if p.In(r) {
				tc.held[b] = true
			}
		}
	}
}

func (tc *TouchControls) justPressed(b TouchButton) bool { return tc.held[b] && !tc.prev[b] }

// Apply merges the touch buttons into in, on top of the keyboard input already
// there. W also jumps, since the keyboard's jump key (Space) has no button.
func (tc *TouchControls) Apply(in *InputState) {
	if !tc.Enabled {
		return
	}
	h := &tc.held

	if h[TouchUp] && !h[TouchDown] {
		in.Direction.UpDown = 1
	} else if h[TouchDown] && !h[TouchUp] {
		in.Direction.UpDown = -1
	}
	if h[TouchLeft] && !h[TouchRight] {
		in.Direction.LeftRight = -1
	} else if h[TouchRight] && !h[TouchLeft] {
		in.Direction.LeftRight = 1
	}
	in.JumpJustPressed = in.JumpJustPressed || h[TouchUp]
	// Sprint stands in for Left Shift, which both runs and dashes.
	in.RunJustPressed = in.RunJustPressed || h[TouchSprint]
	in.DashJustPressed = in.DashJustPressed || h[TouchSprint]

	in.Skills.WeakAttack = in.Skills.WeakAttack || h[TouchWeak]
	in.Skills.StrongAttack = in.Skills.StrongAttack || h[TouchStrong]
	in.Skills.SpecialAttack1 = in.Skills.SpecialAttack1 || tc.justPressed(TouchSpecial1)
	in.Skills.SpecialAttack2 = in.Skills.SpecialAttack2 || tc.justPressed(TouchSpecial2)
	in.Skills.SpecialAttack3 = in.Skills.SpecialAttack3 || h[TouchSpecial3]
	in.Skills.SpecialAttack4 = in.Skills.SpecialAttack4 || h[TouchSpecial4]

	if tc.justPressed(TouchMenu) {
		in.Menu = !in.Menu
	}
}

// Draw renders the buttons over the game, see-through while idle; pressed
// buttons turn dark and solid.
func (tc *TouchControls) Draw(screen *ebiten.Image) {
	if !tc.Enabled {
		return
	}
	if tc.sprites == nil {
		tc.sprites = loadTouchSprites()
	}
	for b := TouchButton(0); b < touchButtonCount; b++ {
		r := tc.rects[b]
		img := tc.sprites.up
		if b == TouchMenu {
			img = tc.sprites.menu
		} else if tc.held[b] {
			img = tc.sprites.down
		}
		alpha := float32(touchIdleAlpha)
		if tc.held[b] {
			alpha = 1
		}
		drawScaled(screen, img, r, nil, alpha)

		if icon := tc.sprites.icons[b]; icon != nil {
			tint := touchIconColor
			if tc.held[b] {
				tint = touchPressedTint
			}
			inset := r.Dx() * touchIconInset / 100
			drawScaled(screen, icon, r.Inset(inset), tint, alpha)
		}
	}
}

func loadTouchSprites() *touchSprites {
	sheet := LoadImage(UIButtonsSheet)
	s := &touchSprites{
		menu: sheet.SubImage(uiMenuRect).(*ebiten.Image),
		up:   sheet.SubImage(uiRoundUpRect).(*ebiten.Image),
		down: sheet.SubImage(uiRoundDownRect).(*ebiten.Image),
	}
	for b, shapes := range touchIconShapes {
		if len(shapes) > 0 {
			s.icons[b] = newTouchIcon(shapes)
		}
	}
	return s
}

// drawScaled draws img stretched to dst at the given opacity, tinted by tint
// when it is not nil.
func drawScaled(screen, img *ebiten.Image, dst image.Rectangle, tint color.Color, alpha float32) {
	b := img.Bounds()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(float64(dst.Dx())/float64(b.Dx()), float64(dst.Dy())/float64(b.Dy()))
	op.GeoM.Translate(float64(dst.Min.X), float64(dst.Min.Y))
	if tint != nil {
		op.ColorScale.ScaleWithColor(tint)
	}
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(img, op)
}
