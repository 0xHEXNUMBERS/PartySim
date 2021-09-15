package board

import "github.com/0xhexnumbers/partysim/mp1"

//BMMBranchPay is a custom branch event for the player to decide if they
//want to pay 10 coins to take a chance at taking the star path.
type BMMBranchPay struct {
	Player     int
	Moves      int
	BowserPath mp1.ChainSpace
	StarPath   mp1.ChainSpace
}

//Responses returns a slice of bools (true/false).
func (b BMMBranchPay) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (b BMMBranchPay) ControllingPlayer() int {
	return b.Player
}

//Handle executes based on r. If r is true, the player pays 10 coins to
//let chance decide which path they take. Otherwise, they take the bowser
//path.
func (b BMMBranchPay) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(b.Player, -10, false)
		g.NextEvent = BMMBranchDecision{
			b.Player, b.Moves, b.BowserPath, b.StarPath,
		}
	} else {
		g.Players[b.Player].CurrentSpace = b.BowserPath
		g.MovePlayer(b.Player, b.Moves-1)
	}
}

//BMMBranchDecision decides which path the player takes.
type BMMBranchDecision struct {
	Player     int
	Moves      int
	BowserPath mp1.ChainSpace
	StarPath   mp1.ChainSpace
}

//Responses returns a slice of the 2 paths the player can take.
func (b BMMBranchDecision) Responses() []mp1.Response {
	return []mp1.Response{b.BowserPath, b.StarPath}
}

func (b BMMBranchDecision) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle moves the player to the ChainSpace r.
func (b BMMBranchDecision) Handle(r mp1.Response, g *mp1.Game) {
	dest := r.(mp1.ChainSpace)
	g.Players[b.Player].CurrentSpace = dest
	g.MovePlayer(b.Player, b.Moves-1)
}

//BMMBowserRoulette decides if bowser steals a star or 20 coins.
type BMMBowserRoulette struct {
	Player int
	Moves  int
}

//Responses returns a slice of bools (true/false).
func (b BMMBowserRoulette) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (b BMMBowserRoulette) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle executes based on r. If r is true, a star is taken from the
//player. If r is false, 20 coins is taken from the palyer.
func (b BMMBowserRoulette) Handle(r mp1.Response, g *mp1.Game) {
	starSteal := r.(bool)
	if starSteal {
		g.Players[b.Player].Stars--
	} else {
		g.AwardCoins(b.Player, -20, false)
	}
	g.MovePlayer(b.Player, b.Moves)
}
