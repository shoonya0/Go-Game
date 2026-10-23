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
	ClipGuard
	ClipHurt
	ClipDeath
)

var clipNames = [...]string{"idle", "walk", "alert", "run", "attack_1", "attack_2", "guard", "hurt", "death"}

func (c ClipName) String() string {
	if c < 0 || int(c) >= len(clipNames) {
		return "clip?"
	}
	return clipNames[c]
}

// requiredClips are the clips every blueprint must have art for.
var requiredClips = []ClipName{ClipIdle, ClipWalk, ClipAttack1, ClipHurt, ClipDeath}

// clipFallback is the clip to try when a species has no art for a clip. Run
// falls back through alert, so a species with only an alert loop (the
// succubus) chases with it. Guard has no fallback: without its art the enemy
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
	EnemyDead:   ClipDeath,
}

// resolveClip returns the clip the species has art for in place of name,
// following fallbacks. ok is false when neither name nor any fallback exists.
func (bp *EnemyBlueprint) resolveClip(name ClipName) (_ ClipName, ok bool) {
	for {
		if _, ok := bp.Clips[name]; ok {
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
	if s == EnemyAttack && e.swing.second {
		name = ClipAttack2
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
func (c EnemyClip) duration() float64 {
	frames := 0.0
	for i := range c.Frames {
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
