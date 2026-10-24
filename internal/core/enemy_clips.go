package core

// ClipName is a standard enemy animation. Every species is described with the
// same names, so the shared AI and drawing code never ask which species they
// run. A species without art for an optional clip falls back to a related one.
type ClipName int

const (
	ClipIdle ClipName = iota
	ClipWalk
	ClipAlert
	ClipRun
	ClipAttack1
	ClipAttack2
	ClipAttackLoop // a phased attack's middle, e.g. a dive's fall
	ClipAttackEnd  // a phased attack's follow-through, e.g. a dive's impact
	ClipGuard
	ClipHurt
	ClipDeath
)

var clipNames = [...]string{"idle", "walk", "alert", "run", "attack_1", "attack_2", "attack_loop", "attack_end", "guard", "hurt", "death"}

func (c ClipName) String() string {
	if c < 0 || int(c) >= len(clipNames) {
		return "clip?"
	}
	return clipNames[c]
}

// requiredClips are the clips every blueprint must have art for.
var requiredClips = []ClipName{ClipIdle, ClipWalk, ClipAttack1, ClipHurt, ClipDeath}

// clipFallback is the clip to try when a species has no art for a clip. Run
// falls back through alert, so a species with only an alert loop chases with
// it. Guard has no fallback: without its art the enemy
// doesn't guard.
var clipFallback = map[ClipName]ClipName{
	ClipRun:     ClipAlert,
	ClipAlert:   ClipWalk,
	ClipAttack2: ClipAttack1,
}

// stateClips is the clip each behavioural state plays. Chasing (EnemyAlert)
// plays run: running is a speed, not a state of its own.
var stateClips = map[EnemyState]ClipName{
	EnemyIdle:   ClipIdle,
	EnemyPatrol: ClipWalk,
	EnemyAlert:  ClipRun,
	EnemyAttack: ClipAttack1,
	EnemyHurt:   ClipHurt,
	EnemyGuard:  ClipGuard, // optional: only species with guard art enter it
	EnemyDead:   ClipDeath,
}

// Has reports whether the species has its own art for clip c. The art is the
// capability: a species with a guard clip can guard, one with attack_2
// alternates two attacks. A clip it only plays through a fallback doesn't count.
func (bp *EnemyBlueprint) Has(c ClipName) bool {
	_, ok := bp.Clips[c]
	return ok
}

// capabilities describes what the species can do, one cell per standard clip
// in clip order: the clip's name when it has the art, "(other)" when it plays
// another clip instead, "-" when it can't do it at all.
func (bp *EnemyBlueprint) capabilities() []string {
	cells := make([]string, len(clipNames))
	for c := range ClipName(len(clipNames)) {
		switch got, ok := bp.resolveClip(c); {
		case !ok:
			cells[c] = "-"
		case got == c:
			cells[c] = c.String()
		default:
			cells[c] = "(" + got.String() + ")"
		}
	}
	return cells
}

// resolveClip returns the clip the species has art for in place of name,
// following fallbacks. ok is false when neither name nor any fallback exists.
func (bp *EnemyBlueprint) resolveClip(name ClipName) (_ ClipName, ok bool) {
	for {
		if bp.Has(name) {
			return name, true
		}
		next, ok := clipFallback[name]
		if !ok {
			return name, false
		}
		name = next
	}
}

// clipFor is the clip the enemy plays in state s. An attack plays the clip its
// swing picked.
func (e *EnemyRuntime) clipFor(s EnemyState) ClipName {
	name := stateClips[s]
	if s == EnemyAttack {
		name = e.swing.clip()
	}
	name, _ = e.bp.resolveClip(name)
	return name
}

// animFor is the animation the enemy plays in state s, or nil if it has none.
func (e *EnemyRuntime) animFor(s EnemyState) *Animation {
	return e.Animations[e.clipFor(s)]
}

// currentAnim is the animation for the enemy's current state.
func (e *EnemyRuntime) currentAnim() *Animation { return e.animFor(e.State) }

// currentClip is the clip description for the enemy's current state.
func (e *EnemyRuntime) currentClip() EnemyClip { return e.bp.Clips[e.clipFor(e.State)] }

// hold is how many base frame durations frame i of the clip stays up.
func (c EnemyClip) hold(i int) float64 {
	if i >= 0 && i < len(c.Hold) && c.Hold[i] > 0 {
		return c.Hold[i]
	}
	return 1
}

// duration is how long one play of the clip takes, in seconds.
func (c EnemyClip) duration() float64 { return c.timeTo(c.Frames) }

// timeTo is how long the clip plays before frame f comes up, in seconds.
func (c EnemyClip) timeTo(f int) float64 {
	frames := 0.0
	for i := range f {
		frames += c.hold(i)
	}
	return frames / c.FPS
}

// held returns the clip with per-frame Hold multipliers, in frame order.
func (c EnemyClip) held(hold ...float64) EnemyClip {
	c.Hold = hold
	return c
}

// moving returns the clip with draw-time motion m.
func (c EnemyClip) moving(m motionFX) EnemyClip {
	c.Motion = m
	return c
}

// hitsOn returns the attack clip with its strike landing on frame f.
func (c EnemyClip) hitsOn(f int) EnemyClip {
	c.HitFrame = f
	return c
}
