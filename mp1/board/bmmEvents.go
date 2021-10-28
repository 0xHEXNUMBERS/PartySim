package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type BMMBranchPayResponse int

const (
	BMMBranchPayPay    BMMBranchPayResponse = iota
	BMMBranchPayIgnore BMMBranchPayResponse = iota
)

func (b BMMBranchPayResponse) String() string {
	switch b {
	case BMMBranchPayPay:
		return "Pay 10 coins to roll die"
	case BMMBranchPayIgnore:
		return "Pass"
	}
	return ""
}

//BMMBranchPay is a custom branch event for the player to decide if they
//want to pay 10 coins to take a chance at taking the star path.
type BMMBranchPay struct {
	Player     int
	Moves      int
	BowserPath mp1.ChainSpace
	StarPath   mp1.ChainSpace
}

func (b BMMBranchPay) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s pay 10 coins to roll the star/bowser die?",
		g.Players[b.Player].Char)
}

func (b BMMBranchPay) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (b BMMBranchPay) Responses() []mp1.Response {
	return []mp1.Response{BMMBranchPayPay, BMMBranchPayIgnore}
}

func (b BMMBranchPay) ControllingPlayer() int {
	return b.Player
}

//Handle executes based on r. If r is true, the player pays 10 coins to
//let chance decide which path they take. Otherwise, they take the bowser
//path.
func (b BMMBranchPay) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(BMMBranchPayResponse)
	if pay == BMMBranchPayPay {
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

func (b BMMBranchDecision) Question(g *mp1.Game) string {
	return fmt.Sprintf("Which space did %s go to?",
		g.Players[b.Player].Char)
}

func (b BMMBranchDecision) Type() mp1.EventType {
	return mp1.CHAINSPACE_EVT_TYPE
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

type BMMBowserRouletteResponse int

const (
	BMMBowserRoulette20Coins BMMBowserRouletteResponse = iota
	BMMBowserRouletteStar
)

func (b BMMBowserRouletteResponse) String() string {
	switch b {
	case BMMBowserRoulette20Coins:
		return "Lose 20 coins"
	case BMMBowserRouletteStar:
		return "Lose 1 star"
	}
	return ""
}

//BMMBowserRoulette decides if bowser steals a star or 20 coins.
type BMMBowserRoulette struct {
	Player int
	Moves  int
}

func (b BMMBowserRoulette) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s lose 20 coins or 1 star?",
		g.Players[b.Player].Char)
}

func (b BMMBowserRoulette) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (b BMMBowserRoulette) Responses() []mp1.Response {
	return []mp1.Response{BMMBowserRoulette20Coins, BMMBowserRouletteStar}
}

func (b BMMBowserRoulette) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle executes based on r. If r is true, a star is taken from the
//player. If r is false, 20 coins is taken from the palyer.
func (b BMMBowserRoulette) Handle(r mp1.Response, g *mp1.Game) {
	starSteal := r.(BMMBowserRouletteResponse)
	if starSteal == BMMBowserRouletteStar {
		g.Players[b.Player].Stars--
	} else {
		g.AwardCoins(b.Player, -20, false)
	}
	g.MovePlayer(b.Player, b.Moves)
}
