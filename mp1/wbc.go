package mp1

//wbcBoardData holds all of the board specific data related to WBC.
type wbcBoardData struct {
	Direction bool
}

//wbcCannonShot sets the player's new chain, and set the next event to set
//the player's new position.
func wbcCannonShot(g *Game, player, moves int) int {
	newChain := g.Players[player].CurrentSpace.Chain
	data := g.Board.Data.(wbcBoardData)
	if data.Direction {
		newChain = (newChain + 3) % 4
	} else {
		newChain = (newChain + 1) % 4
	}
	g.ExtraEvent = wbcCannon{
		player, moves, newChain,
	}
	return moves
}

//wbcReverseCannons reverses the cannons' direction.
func wbcReverseCannons(g *Game, player int) {
	data := g.Board.Data.(wbcBoardData)
	data.Direction = !data.Direction
	g.Board.Data = data
}

//wbcLoadPlayerInBowserCannon sets the next event to choosing a chain for
//the player to land on.
func wbcLoadPlayerInBowserCannon(g *Game, player, moves int) int {
	g.ExtraEvent = wbcBowserCannon{player, moves}
	return moves
}

//wbcShyGuy occurs when a player passes shyguy. If the player has >=10
//coins, then the next event is set for the player to respond to the
//shyguy.
func wbcShyGuy(g *Game, player, moves int) int {
	if g.Players[player].Coins >= 10 {
		g.ExtraEvent = wbcShyGuyEvent{player, moves}
	}
	return moves
}

//WBC holds the data for Wario's Battle Canyon.
var WBC = Board{
	Chains: &[]Chain{
		{ //Bottom Left
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Bottom Right
			{Type: Blue},
			{Type: Mushroom},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: wbcReverseCannons},
			{Type: Star},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Top Left
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Boo},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Top Right
			{Type: Red},
			{Type: Bowser},
			{Type: Red},
			{Type: Red},
			{Type: Red},
			{Type: Invisible, PassingEvent: wbcShyGuy},
			{Type: Red},
			{Type: Bowser},
			{Type: Red},
			{Type: Happening, StoppingEvent: wbcReverseCannons},
			{Type: Star},
			{Type: Blue},
			{Type: Red},
			{Type: Invisible, PassingEvent: wbcCannonShot},
		},
		{ //Center
			{Type: MinigameSpace},
			{Type: MinigameSpace},
			{Type: MinigameSpace},
			{Type: Star},
			{Type: MinigameSpace},
			{Type: MinigameSpace},
			{Type: MinigameSpace},
			{Type: BogusItem},
			{Type: Invisible, PassingEvent: wbcLoadPlayerInBowserCannon},
		},
	},
	Links:       nil,
	BowserCoins: 20,
	Data:        wbcBoardData{},
}
