package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Chain struct {
	Treasury     float64
	CurrentEpoch int
	Projects     map[string]*Project
	Identities   map[string]*Identity
}

func NewChain() *Chain {
	return &Chain{
		Treasury:     0,
		CurrentEpoch: 0,
		Projects:     make(map[string]*Project),
		Identities:   make(map[string]*Identity),
	}
}

func (c *Chain) GetReputation(pubKey string) *Identity {
	if _, exists := c.Identities[pubKey]; !exists {
		c.Identities[pubKey] = &Identity{PubKey: pubKey, Score: 0, Burned: false}
	}
	return c.Identities[pubKey]
}

func (c *Chain) CalculateBond(liquidity float64, risk RiskClass) float64 {
	multiplier := 0.1 // 10%
	switch risk {
	case RiskHigh:
		multiplier = 0.5
	case RiskDegen:
		multiplier = 1.0 // 100% Bond
	}
	return liquidity * multiplier
}

func (c *Chain) DeployProject(name, owner string, liquidity float64, risk RiskClass) string {
	id := c.GetReputation(owner)
	if id.Burned {
		fmt.Printf("🛑 REJECTED: Identity %s is BURNED.\n", owner[:8])
		return ""
	}

	bond := c.CalculateBond(liquidity, risk)
	
	// Reputation Discount (Max 50%)
	discount := float64(id.Score) * 0.01
	if discount > 0.5 {
		discount = 0.5
	}
	bond = bond * (1.0 - discount)

	// Hash ID
	hash := sha256.Sum256([]byte(name + fmt.Sprint(time.Now().UnixNano())))
	projID := hex.EncodeToString(hash[:])[:8]

	c.Projects[projID] = &Project{
		ID:           projID,
		Name:         name,
		Owner:        owner,
		BondLocked:   bond,
		Liquidity:    liquidity,
		StartEpoch:   c.CurrentEpoch,
		LockDuration: 90,
		Status:       "ACTIVE",
	}

	fmt.Printf("🚀 LAUNCH: '%s' (ID: %s) | Bond: %.2f ETH\n", name, projID, bond)
	return projID
}

func (c *Chain) AdvanceTime(epochs int) {
	c.CurrentEpoch += epochs
}

// CheckRugConditions is the deterministic oracle
func (c *Chain) CheckRugConditions(projID string, action string, amount float64) bool {
	p, exists := c.Projects[projID]
	if !exists {
		return false
	}

	if action == "REMOVE_LIQUIDITY" {
		// Rug Condition: Removing > 10% before unlock
		allowed := p.Liquidity * 0.10
		if amount > allowed {
			return true
		}
	}
	return false
}

func (c *Chain) Slash(projID string, reason string) {
	p := c.Projects[projID]
	if p.Status == "RUGGED" {
		return
	}

	fmt.Printf("⚔️ SLASHING '%s' | Reason: %s\n", p.Name, reason)
	
	// Confiscate
	c.Treasury += p.BondLocked
	p.BondLocked = 0
	p.Status = "RUGGED"

	// Burn Identity
	id := c.GetReputation(p.Owner)
	id.Burned = true
	id.Score = 0
	fmt.Printf("   💀 Identity Burned: %s...\n", p.Owner[:8])
}

// ActionWithdrawLiquidity attempts to remove funds
func (c *Chain) ActionWithdrawLiquidity(projID string, amount float64) bool {
	p := c.Projects[projID]
	fmt.Printf("🔄 '%s' attempting to withdraw %.2f Liq...\n", p.Name, amount)

	if c.CheckRugConditions(projID, "REMOVE_LIQUIDITY", amount) {
		c.Slash(projID, "Liquidity Rug Pull Detected")
		return false // Failed withdrawal (Rug blocked/punished)
	}

	p.Liquidity -= amount
	fmt.Println("   ✅ Withdrawal Approved.")
	return true
}

// ActionClaimBond attempts to unlock bond
func (c *Chain) ActionClaimBond(projID string) {
	p := c.Projects[projID]
	if p.Status != "ACTIVE" {
		return
	}

	alive := c.CurrentEpoch - p.StartEpoch
	if alive >= p.LockDuration {
		amt := p.BondLocked
		p.BondLocked = 0
		p.Status = "GRADUATED"
		
		// Rep Gain
		id := c.GetReputation(p.Owner)
		id.Score += int(amt)
		
		fmt.Printf("💰 Bond Returned: %.2f ETH | New Rep: %d\n", amt, id.Score)
	}
}
