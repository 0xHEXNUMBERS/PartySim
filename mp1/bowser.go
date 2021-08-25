package mp1

//PreBowserCheck is a check that happens before the player visits bowser.
//If certain conditions are met, bowser will perform different actions
//instead of picking one of his normal actions.
func (g *Game) PreBowserCheck(player int) {
	//Special events when player has 0 coins
	if g.Players[player].Coins == 0 {
		if g.Players[player].Stars > 0 {
			g.AwardCoins(player, 10, false)
			g.Players[player].Stars--
		} else {
			g.AwardCoins(player, 20, false)
		}
		g.EndCharacterTurn()
	} else {
		g.NextEvent = BowserEvent{player}
	}
}

//BowserEvent is the normal event that occurs when a player lands on a
//bowser space.
type BowserEvent struct {
	Player int
}

//BowserResponse is an enumeration for all of Bowser's normal actions.
type BowserResponse int

const (
	CoinsForBowser BowserResponse = iota
	BowserBalloonBurst
	BowsersFaceLift
	BowsersTugoWar
	BashnCash
	BowserRevolution
	BowsersChanceTime
	StarPresent
)

//Responses returns a slice of all of Bowser's actions.
func (b BowserEvent) Responses() []Response {
	return []Response{
		CoinsForBowser,
		BowserBalloonBurst,
		BowsersFaceLift,
		BowsersTugoWar,
		BashnCash,
		BowserRevolution,
		BowsersChanceTime,
		StarPresent,
	}
}

func (b BowserEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle takes the action r and executes it, setting a new event if needed.
func (b BowserEvent) Handle(r Response, g *Game) {
	choice := r.(BowserResponse)
	switch choice {
	case CoinsForBowser:
		coinsLost := GetBowserMinigameCoinLoss(g.Turn)
		g.AwardCoins(b.Player, -coinsLost, false)
		g.EndCharacterTurn()
	case BowserBalloonBurst:
		g.NextEvent = BowserBalloonBurstEvent{}
	case BowsersFaceLift:
		g.NextEvent = BowsersFaceLiftEvent{b.Player}
	case BowsersTugoWar:
		g.NextEvent = BowsersTugoWarEvent{b.Player}
	case BashnCash:
		g.NextEvent = BowsersBashnCash{b.Player, g.Players[b.Player].Coins}
	case BowserRevolution:
		coins := 0
		for i := range g.Players {
			coins += g.Players[i].Coins
		}
		coins /= 4
		for i := range g.Players {
			g.Players[i].Coins = coins
		}
		g.EndCharacterTurn()
	case BowsersChanceTime:
		g.NextEvent = BowsersChanceTimeEvent{}
	}
}

func GetBowserMinigameCoinLoss(turn uint8) int {
	if turn <= 9 {
		return 10
	} else if turn <= 19 {
		return 20
	} else if turn <= 29 {
		return 30
	}
	return 40
}

//BowserBalloonBurstEvent holds the implementation for Bowser's Balloon Burst.
type BowserBalloonBurstEvent struct{}

//Responses returns a slice of ints from [0, 4].
func (b BowserBalloonBurstEvent) Responses() []Response {
	return CPURangeEvent{0, 4}.Responses()
}

func (b BowserBalloonBurstEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle takes coins from every player except for player r. If r == 4,
//then every player loses 20 coins.
func (b BowserBalloonBurstEvent) Handle(r Response, g *Game) {
	winner := r.(int)
	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	if winner == 4 {
		for p := range g.Players {
			g.AwardCoins(p, -20, true)
		}
	} else {
		for p := range g.Players {
			if p == winner {
				continue
			}
			g.AwardCoins(p, coinLoss, true)
		}
	}
	g.EndCharacterTurn()
}

//BowsersFaceLiftEvent holds the implementation for Bowser's Face Lift.
type BowsersFaceLiftEvent struct {
	Player int
}

//Responses returns a slice of ints from [1, 15].
func (b BowsersFaceLiftEvent) Responses() []Response {
	return CPURangeEvent{1, 15}.Responses()
}

func (b BowsersFaceLiftEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates coin loses for each player depending on the avlue of r.
//Each bit p in r, if set to 0, loses player p some number of coins. If r
//is 15, then b.Player loses 50 coins.
func (b BowsersFaceLiftEvent) Handle(r Response, g *Game) {
	results := r.(int)
	if results == 15 { //All players won
		g.AwardCoins(b.Player, -50, true)
		return
	}

	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	for p := range g.Players {
		if results&(1<<p) == 0 {
			g.AwardCoins(p, coinLoss, true)
		}
	}
	g.EndCharacterTurn()
}

//BowsersTugoWarEvent holds the implementation for Bowser's Tug o War.
type BowsersTugoWarEvent struct {
	Player int
}

//BowsersTugoWarResult is an enumeration of ending results of Bowser's
//Tug o War.
type BowsersTugoWarResult int

const (
	BTWDraw BowsersTugoWarResult = iota
	BTW1TWin
	BTW3TWin
)

var BTWResults = []Response{
	BTWDraw,
	BTW1TWin,
	BTW3TWin,
}

//Responses returns a slice of the valid end results of Bowser's Tug o War.
func (b BowsersTugoWarEvent) Responses() []Response {
	return BTWResults
}

func (b BowsersTugoWarEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates coin loses given a valid end result r.
func (b BowsersTugoWarEvent) Handle(r Response, g *Game) {
	results := r.(BowsersTugoWarResult)
	switch results {
	case BTWDraw:
		for p := range g.Players {
			g.AwardCoins(p, -30, true)
		}
	case BTW1TWin:
		coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
		for p := range g.Players {
			if p != b.Player {
				g.AwardCoins(p, coinLoss, true)
			}
		}
	case BTW3TWin:
		g.AwardCoins(b.Player, -10, true)
	}
	g.EndCharacterTurn()
}

//BowsersBashnCash holds the implementation of Bowser's BashnCash.
type BowsersBashnCash struct {
	Player int
	Coins  int
}

//Responses returns a slice of ints from [1, max], where max is the
//maximum number of hits the solo player can take before losing all of
//their coins.
func (b BowsersBashnCash) Responses() []Response {
	max := b.Coins / 5
	max += b.Coins % 5
	return CPURangeEvent{1, max}.Responses()
}

func (b BowsersBashnCash) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates the coin loss for the solo player based on the number
//of hits they took.
func (b BowsersBashnCash) Handle(r Response, g *Game) {
	timesHit := r.(int)
	coinsLost := 0
	if b.Coins/5 < timesHit {
		coinsLost += b.Coins - (b.Coins % 5)
		timesHit -= b.Coins / 5
		coinsLost += timesHit
	} else {
		coinsLost += timesHit * 5
	}
	g.AwardCoins(b.Player, -coinsLost, true)
	g.EndCharacterTurn()
}

//BowsersChanceTimeEvent holds the implementation for Bowser's Chance Time.
type BowsersChanceTimeEvent struct{}

//BCTResponse is a valid response to Bowser's Chance Time Event.
type BCTResponse struct {
	Player int
	Coins  int
}

var BCTResponses = []Response{
	BCTResponse{0, 10},
	BCTResponse{0, 20},
	BCTResponse{0, 30},
	BCTResponse{1, 10},
	BCTResponse{1, 20},
	BCTResponse{1, 30},
	BCTResponse{2, 10},
	BCTResponse{2, 20},
	BCTResponse{2, 30},
}

//Responses return the valid responses to Bowser's Chance Time Event.
func (b BowsersChanceTimeEvent) Responses() []Response {
	return BCTResponses
}

func (b BowsersChanceTimeEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates coin loss given a valid Response r.
func (b BowsersChanceTimeEvent) Handle(r Response, g *Game) {
	res := r.(BCTResponse)
	g.AwardCoins(res.Player, -res.Coins, false)
	g.EndCharacterTurn()
}
