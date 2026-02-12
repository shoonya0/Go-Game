package enemy

// ========================================================================
// Party System - Social behavior, forming groups, betrayal
// ========================================================================

// Party represents a group of enemies working together
type Party struct {
	ID      int
	Leader  *EnemyRuntime
	Members []*EnemyRuntime
}

// PartyManager handles all parties in the game
type PartyManager struct {
	Parties     map[int]*Party
	NextPartyID int
}
