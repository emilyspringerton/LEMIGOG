package core

import (
	"strings"
	"testing"
)

func TestHonestLaunch(t *testing.T) {
	chain := NewChain()
	owner := "0xHonestBuilder" + strings.Repeat("b", 32)

	// 1. Launch
	pid := chain.DeployProject("SafeDeFi", owner, 100.0, RiskHigh)
	if pid == "" {
		t.Fatal("Failed to deploy project")
	}

	// 2. Wait
	chain.AdvanceTime(90)

	// 3. Claim
	chain.ActionClaimBond(pid)

	// Assert
	id := chain.GetReputation(owner)
	if id.Score <= 0 {
		t.Errorf("Expected positive reputation, got %d", id.Score)
	}
	if chain.Projects[pid].Status != "GRADUATED" {
		t.Errorf("Expected status GRADUATED, got %s", chain.Projects[pid].Status)
	}
}

func TestRugPullScenario(t *testing.T) {
	chain := NewChain()
	owner := "0xScammer" + strings.Repeat("x", 32)

	// 1. Launch Degen Pool (100% Bond required)
	// Liquidity: 100, Bond: 100
	pid := chain.DeployProject("RugCoin", owner, 100.0, RiskDegen)

	// 2. Try to Rug Early
	chain.AdvanceTime(10)
	
	// Attempt to pull 90% of liquidity
	success := chain.ActionWithdrawLiquidity(pid, 90.0)

	// Assert
	if success {
		t.Error("Rug pull should have been blocked/slashed")
	}

	p := chain.Projects[pid]
	if p.Status != "RUGGED" {
		t.Errorf("Project status should be RUGGED, got %s", p.Status)
	}
	
	if p.BondLocked > 0 {
		t.Error("Bond should have been slashed to 0")
	}

	if chain.Treasury != 100.0 {
		t.Errorf("Treasury should have 100.0 ETH, got %.2f", chain.Treasury)
	}

	id := chain.GetReputation(owner)
	if !id.Burned {
		t.Error("Scammer identity should be burned")
	}
}

func TestBurnedIdentityCannotDeploy(t *testing.T) {
	chain := NewChain()
	owner := "0xBadActor"

	// 1. Burn him
	id := chain.GetReputation(owner)
	id.Burned = true

	// 2. Try deploy
	pid := chain.DeployProject("NewScam", owner, 10.0, RiskLow)

	if pid != "" {
		t.Error("Burned identity should not be able to deploy")
	}
}
