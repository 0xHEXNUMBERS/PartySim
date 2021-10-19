package board

import "github.com/0xhexnumbers/partysim/mp1"

//wbcBoardData holds all of the board specific data related to WBC.
type wbcBoardData struct {
	Direction bool
}

//wbcCannonShot sets the player's new chain, and set the next event to set
//the player's new position.
func wbcCannonShot(g *mp1.Game, player, moves int) int {
	newChain := g.Players[player].CurrentSpace.Chain
	data := g.Board.Data.(wbcBoardData)
	if data.Direction {
		newChain = (newChain + 3) % 4
	} else {
		newChain = (newChain + 1) % 4
	}
	g.NextEvent = WBCCannon{
		player, moves, newChain,
	}
	return moves
}

//wbcReverseCannons reverses the cannons' direction.
func wbcReverseCannons(g *mp1.Game, player int) {
	data := g.Board.Data.(wbcBoardData)
	data.Direction = !data.Direction
	g.Board.Data = data
}

//wbcLoadPlayerInBowserCannon sets the next event to choosing a random position
//for the player to land on.
func wbcLoadPlayerInBowserCannon(g *mp1.Game, player, moves int) int {
	g.NextEvent = WBCBowserCannon{player, moves}
	return moves
}

//wbcShyGuy occurs when a player passes shyguy. If the player has >=10
//coins, then the next event is set for the player to respond to the
//shyguy.
func wbcShyGuy(g *mp1.Game, player, moves int) int {
	if g.Players[player].Coins >= 10 {
		g.NextEvent = WBCShyGuyEvent{player, moves}
	}
	return moves
}

//WBC holds the data for Wario's Battle Canyon.
var WBC = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Bottom Left
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Start},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Bottom Right
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: wbcReverseCannons},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Top Left
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Boo},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Top Right
			{Type: mp1.Red},
			{Type: mp1.Bowser},
			{Type: mp1.Red},
			{Type: mp1.Red},
			{Type: mp1.Red},
			{Type: mp1.Invisible, PassingEvent: wbcShyGuy},
			{Type: mp1.Red},
			{Type: mp1.Bowser},
			{Type: mp1.Red},
			{Type: mp1.Happening, StoppingEvent: wbcReverseCannons},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Center
			{Type: mp1.MinigameSpace},
			{Type: mp1.MinigameSpace},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Star},
			{Type: mp1.MinigameSpace},
			{Type: mp1.MinigameSpace},
			{Type: mp1.MinigameSpace},
			{Type: mp1.BogusItem},
			{Type: mp1.Invisible, PassingEvent: wbcLoadPlayerInBowserCannon},
		},
	},
	Links:       nil,
	BowserCoins: 20,
	Data:        wbcBoardData{},
}
