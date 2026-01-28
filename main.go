package main

import (
	"fmt"
	"github.com/emilyspringerton/lemigog/core"
)

func main() {
	fmt.Println("🍋 LEMIGOG NODE v0.1 (Go)")
	
	chain := core.NewChain()
	
	// Simple Simulation Run
	fmt.Println("\n--- SIMULATION START ---")
	owner := "0xSatoshi"
	pid := chain.DeployProject("Bitcoin 2", owner, 1000.0, core.RiskLow)
	
	chain.AdvanceTime(100)
	chain.ActionClaimBond(pid)
	
	fmt.Println("--- SIMULATION END ---")
}
