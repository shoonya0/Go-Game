package core

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyHealthSheet holds the three states of one enemy health bar: full, half
// and empty, like the player's stamina sheet.
const EnemyHealthSheet = "../assets/UI/enemyHealth.png"

// Enemy health is counted in bars, shown above the enemy like the player's
// stamina segments. Player attack damage is tuned in half bars.
const (
	EnemyHealthPerBar       = 20  // hit points in one bar
	MinEnemyHealthBars      = 2   // every enemy has at least this many bars
	EnemyRegenBarsPerSecond = 0.1 // a live enemy regains this many bars per second

	enemyBarW    = 28 // drawn width of one bar
	enemyBarH    = 6  // drawn height of one bar
	enemyBarGap  = 2  // space between bars
	enemyBarLift = 4  // space between the top of the sprite and the bars
)

// enemyBarRects are the full/half/empty source rects inside EnemyHealthSheet,
// measured from the non-transparent pixels.
var enemyBarRects = [3]image.Rectangle{
	image.Rect(42, 9, 668, 101),
	image.Rect(57, 137, 663, 235),
	image.Rect(56, 253, 662, 341),
}

var enemyBars *hudSprite

// healthBars is how many bars the enemy's full health fills.
func (e *EnemyRuntime) healthBars() int {
	return (e.bp.Stats.MaxHealth + EnemyHealthPerBar - 1) / EnemyHealthPerBar
}

// healthBarStates is the full/half/empty state of each health bar. Unlike the
// player's HUD it rounds down to the half bar, so regenerating health doesn't
// refill a bar until a whole half is back; any health left still shows half a
// bar.
func (e *EnemyRuntime) healthBarStates() []segmentState {
	halves := e.Health / (EnemyHealthPerBar / 2)
	if e.Health > 0 && halves == 0 {
		halves = 1
	}
	return halfStepStates(halves, e.healthBars())
}

// regenerateHealth refills a live enemy's health over time, up to its maximum.
// Like water damage, the fractional part carries between ticks so integer
// health rises steadily.
func (e *EnemyRuntime) regenerateHealth() {
	if e.Health >= e.bp.Stats.MaxHealth {
		e.regenAccum = 0
		return
	}
	e.regenAccum += EnemyRegenBarsPerSecond * EnemyHealthPerBar * deltaTime()
	if whole := int(e.regenAccum + 1e-9); whole > 0 { // epsilon: per-tick sums fall just short

		e.Health = min(e.Health+whole, e.bp.Stats.MaxHealth)
		e.regenAccum -= float64(whole)
	}
}

// DrawHealthBar draws the enemy's health bars side by side, centred above its
// sprite. Dead enemies have none.
func (e *EnemyRuntime) DrawHealthBar(screen *ebiten.Image, cam Camera) {
	anim := e.currentAnim()
	if e.State.IsDead() || anim == nil {
		return
	}
	if enemyBars == nil {
		enemyBars = loadHUDSprite(EnemyHealthSheet, enemyBarRects)
	}

	n := e.healthBars()
	body := e.GetBounds()
	width := float64(n)*enemyBarW + float64(n-1)*enemyBarGap
	spriteTop := body.Y + body.Height - float64(anim.FrameHeight)*e.Scale
	x := centerX(body) - width/2 - cam.Pos.X
	y := spriteTop - enemyBarLift - enemyBarH - cam.Pos.Y
	drawSegments(screen, enemyBars, e.healthBarStates(), x, y, enemyBarW, enemyBarH, enemyBarGap)
}
