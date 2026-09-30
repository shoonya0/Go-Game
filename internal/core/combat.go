package core

// Combat holds the player's health and power pools, each on a 0-100 scale.
type Combat struct {
	Health    float64 // current health
	MaxHealth float64 // maximum health
	Power     float64 // current power
	MaxPower  float64 // maximum power
}
