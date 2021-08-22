package mp1

//bmmBranchPay is a custom branch event for the player to decide if they
//want to pay 10 coins to take a chance at taking the star path.
type bmmBranchPay struct {
	Player     int
	Moves      int
	BowserPath ChainSpace
	StarPath   ChainSpace
}

//Responses returns a slice of bools (true/false).
func (b bmmBranchPay) Responses() []Response {
	return []Response{true, false}
}

func (b bmmBranchPay) ControllingPlayer() int {
	return b.Player
}

//Handle executes based on r. If r is true, the player pays 10 coins to
//let chance decide which path they take. Otherwise, they take the bowser
//path.
func (b bmmBranchPay) Handle(r Response, g *Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(b.Player, -10, false)
		g.ExtraEvent = bmmBranchDecision{
			b.Player, b.Moves, b.BowserPath, b.StarPath,
		}
	} else {
		g.Players[b.Player].CurrentSpace = b.BowserPath
		g.MovePlayer(b.Player, b.Moves-1)
	}
}

//bmmBranchDecision decides which path the player takes.
type bmmBranchDecision struct {
	Player     int
	Moves      int
	BowserPath ChainSpace
	StarPath   ChainSpace
}

//Responses returns a slice of the 2 paths the player can take.
func (b bmmBranchDecision) Responses() []Response {
	return []Response{b.BowserPath, b.StarPath}
}

func (b bmmBranchDecision) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle moves the player to the ChainSpace r.
func (b bmmBranchDecision) Handle(r Response, g *Game) {
	dest := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = dest
	g.MovePlayer(b.Player, b.Moves-1)
}

//bmmBowserRoulette decides if bowser steals a star or 20 coins.
type bmmBowserRoulette struct {
	Player int
	Moves  int
}

//Responses returns a slice of bools (true/false).
func (b bmmBowserRoulette) Responses() []Response {
	return []Response{true, false}
}

func (b bmmBowserRoulette) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle executes based on r. If r is true, a star is taken from the
//player. If r is false, 20 coins is taken from the palyer.
func (b bmmBowserRoulette) Handle(r Response, g *Game) {
	starSteal := r.(bool)
	if starSteal {
		g.Players[b.Player].Stars--
	} else {
		g.AwardCoins(b.Player, -20, false)
	}
	g.MovePlayer(b.Player, b.Moves)
}
