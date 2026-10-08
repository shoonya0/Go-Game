package core

// EnemyState is the behavioural state shared by every enemy, regardless of which
// sprite sheet backs it. The succubus sheet groups its frames exactly this way
// (IDLE/PATROL, ALERT, ATTACK, HURT, DIED), and other enemies reuse the same set
// even when their frame counts and pixel sizes differ. Each enemy's Animations
// map is keyed by these values, so the general update/draw code never needs to
// know which concrete enemy it is driving.
type EnemyState int

const (
	EnemyIdle   EnemyState = iota // standing still at rest
	EnemyPatrol                   // pacing around the spawn point
	EnemyAlert                    // chasing the player after spotting them
	EnemyAttack                   // striking the player (one-shot)
	EnemyHurt                     // reacting to damage (one-shot)
	EnemyDead                     // death animation, then inert (one-shot)
)

// IsDead reports whether the enemy has been killed.
func (s EnemyState) IsDead() bool { return s == EnemyDead }

// IsHurt reports whether the enemy is playing its hurt reaction.
func (s EnemyState) IsHurt() bool { return s == EnemyHurt }

// IsAttacking reports whether the enemy is mid-attack.
func (s EnemyState) IsAttacking() bool { return s == EnemyAttack }

// IsAlert reports whether the enemy is actively chasing the player.
func (s EnemyState) IsAlert() bool { return s == EnemyAlert }

// IsOneShot reports whether the state is a non-looping action that should hand
// back to a follow-up state once its animation finishes (see enemy_anim.go).
func (s EnemyState) IsOneShot() bool {
	return s == EnemyAttack || s == EnemyHurt || s == EnemyDead
}
