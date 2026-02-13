package enemy

const (
	StateIdle int = iota
	StatePatrolling
	StateHunting      // Actively tracking player
	StateFleeing      // Running away (health critical)
	StateResting      // Regenerating health
	StateBerserk      // High damage mode, reduced IQ
	StateFormingParty // Looking for allies
	StateInParty      // Part of a group
	StateBetraying    // About to kill ally for stats
	StateJumping
	StateFalling
	StateLanding
	StateAttacking
	StateDefending
	StateDead
)

// Checking Enemy State
func (s *EnemyState) IsEnemyIdle() bool {
	return s.Current == StateIdle
}

func (s *EnemyState) IsEnemyPatrolling() bool {
	return s.Current == StatePatrolling
}

func (s *EnemyState) IsEnemyHunting() bool {
	return s.Current == StateHunting
}

func (s *EnemyState) IsEnemyFleeing() bool {
	return s.Current == StateFleeing
}

func (s *EnemyState) IsEnemyResting() bool {
	return s.Current == StateResting
}

func (s *EnemyState) IsEnemyBerserk() bool {
	return s.Current == StateBerserk
}

func (s *EnemyState) IsEnemyFormingParty() bool {
	return s.Current == StateFormingParty
}

func (s *EnemyState) IsEnemyInParty() bool {
	return s.Current == StateInParty
}

func (s *EnemyState) IsEnemyBetraying() bool {
	return s.Current == StateBetraying
}

func (s *EnemyState) IsEnemyJumping() bool {
	return s.Current == StateJumping
}

func (s *EnemyState) IsEnemyFalling() bool {
	return s.Current == StateFalling
}

func (s *EnemyState) IsEnemyAttacking() bool {
	return s.Current == StateAttacking
}

func (s *EnemyState) IsEnemyDefending() bool {
	return s.Current == StateDefending
}

func (s *EnemyState) IsEnemyDead() bool {
	return s.Current == StateDead
}

func (s *EnemyState) IsEnemyGrounded() bool {
	return s.Current != StateJumping && s.Current != StateFalling
}

func (s *EnemyState) IsEnemyState(state int) bool {
	return s.Current == state
}

// Setting Enemy State
func (s *EnemyState) SetEnemyState(newState int) {
	s.Previous = s.Current
	s.Current = newState
}

// PartyStatus represents the enemy's social status
type PartyStatus int

const (
	PartySolo   PartyStatus = iota // Alone
	PartyMember                    // Part of a group but not leader
	PartyLeader                    // Leading a group
)

// ActionType represents possible actions the AI can take
type ActionType int

const (
	ActionIdle ActionType = iota
	ActionMoveLeft
	ActionMoveRight
	ActionJump
	ActionJumpLeft
	ActionJumpRight
	ActionAttack
	ActionDefend
	ActionFlee
	ActionEnterBerserk
	ActionExitBerserk
	ActionFormParty
	ActionBetrayAlly
	ActionRest
	ActionHuntPlayer
)

// check if AI can take an action
func (s *EnemyState) CanEnemyMove() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateAttacking
}

func (s *EnemyState) CanEnemyJump() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFalling &&
		s.Current != StateJumping && s.Current != StateAttacking && s.Current != StateDefending
}

func (s *EnemyState) CanEnemyAttack() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFleeing
}

func (s *EnemyState) CanEnemyDefend() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFleeing && s.Current != StateFalling && s.Current != StateJumping
}

func (s *EnemyState) CanEnemyFlee() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFleeing && s.Current != StateFalling && s.Current != StateJumping && s.Current != StateAttacking && s.Current != StateDefending
}

func (s *EnemyState) CanEnemyRest() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFleeing && s.Current != StateFalling && s.Current != StateJumping && s.Current != StateAttacking && s.Current != StateDefending
}

func (s *EnemyState) CanEnemyEnterBerserk() bool {
	return s.Current != StateDead && s.Current != StateResting && s.Current != StateFleeing && s.Current != StateFalling && s.Current != StateJumping && s.Current != StateAttacking && s.Current != StateDefending
}

func (s *EnemyState) GetEnemyState() int {
	return s.Current
}
