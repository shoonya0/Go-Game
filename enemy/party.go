package enemy

import "fmt"

// ========================================================================
// Party System - Social behavior, forming groups, betrayal
// ========================================================================

// Party represents a group of enemies working together
type Party struct {
	ID      string
	Leader  *EnemyRuntime
	Members []*EnemyRuntime
}

// PartyManager handles all parties in the game
type PartyManager struct {
	Parties     map[string]*Party
	nextPartyID int
}

// default party manager is -1
func InitPartyManager() PartyManager {
	return PartyManager{
		Parties:     make(map[string]*Party),
		nextPartyID: 0,
	}
}

func (pm *PartyManager) GeneratePartyID(managerID string) string {
	pm.nextPartyID++
	if managerID != "" {
		return fmt.Sprintf("%s-P-%d", managerID, pm.nextPartyID)
	}
	return fmt.Sprintf("P-%d", pm.nextPartyID)
}
