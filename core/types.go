package core

import "time"

// RiskClass defines the risk profile of a project
type RiskClass int

const (
	RiskLow RiskClass = iota
	RiskHigh
	RiskDegen
)

// Project represents a deployed token or pool
type Project struct {
	ID           string
	Name         string
	Owner        string // PubKey
	BondLocked   float64
	Liquidity    float64
	StartEpoch   int
	LockDuration int
	Status       string // "ACTIVE", "RUGGED", "GRADUATED"
}

// Identity represents a developer's reputation
type Identity struct {
	PubKey string
	Score  int
	Burned bool
}
