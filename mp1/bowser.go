package mp1

import (
	"fmt"
	"strconv"
)

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

func (b BowserResponse) String() string {
	switch b {
	case CoinsForBowser:
		return "Coins For Bowser"
	case BowserBalloonBurst:
		return "Bowser's Balloon Burst"
	case BowsersFaceLift:
		return "Bowser's Face Lift"
	case BowsersTugoWar:
		return "Bowser's Tug o' War"
	case BashnCash:
		return "Bowser's Bash n' Cash"
	case BowserRevolution:
		return "Bowser Revolution"
	case BowsersChanceTime:
		return "Bowser's Chance Time"
	case StarPresent:
		return "Star Present"
	}
	return ""
}

func (b BowserEvent) Question(g *Game) string {
	return "What did Bowser pick as the punishment?"
}

func (b BowserEvent) Type() EventType {
	return ENUM_EVT_TYPE
}

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
		g.NextEvent = BowserBalloonBurstEvent{Range{0, 4}}
	case BowsersFaceLift:
		g.NextEvent = BowsersFaceLiftEvent{Range{0, 15}, b.Player}
	case BowsersTugoWar:
		g.NextEvent = BowsersTugoWarEvent{b.Player}
	case BashnCash:
		g.NextEvent = NewBowsersBashnCash(b.Player, g.Players[b.Player].Coins)
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
type BowserBalloonBurstEvent struct{ Range }

func (b BowserBalloonBurstEvent) Question(g *Game) string {
	return "Which player popped the ballon?"
}

func (b BowserBalloonBurstEvent) Type() EventType {
	return PLAYER_EVT_TYPE
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
	Range
	Player int
}

func (b BowsersFaceLiftEvent) Question(g *Game) string {
	return "Which players won Bowser's Face Lift?"
}

func (b BowsersFaceLiftEvent) Type() EventType {
	return MULTIWIN_PLAYER_EVT_TYPE
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

func (b BowsersTugoWarResult) String() string {
	switch b {
	case BTWDraw:
		return "Draw"
	case BTW1TWin:
		return "Single Player Wins"
	case BTW3TWin:
		return "Team of 3 Wins"
	}
	return ""
}

var BTWResults = []Response{
	BTWDraw,
	BTW1TWin,
	BTW3TWin,
}

func (b BowsersTugoWarEvent) Question(g *Game) string {
	return "Which Team won Bowser's Tug o War."
}

func (b BowsersTugoWarEvent) Type() EventType {
	return ENUM_EVT_TYPE
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
	Range
	Player int
	Coins  int
}

func NewBowsersBashnCash(player, coins int) BowsersBashnCash {
	max := coins / 5
	max += coins % 5
	return BowsersBashnCash{
		Range{1, max},
		player,
		coins,
	}
}

func (b BowsersBashnCash) Question(g *Game) string {
	return fmt.Sprintf("How many times did %s get hit?",
		g.Players[b.Player].Char)
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

func (b BCTResponse) String() string {
	givingPlayer := strconv.Itoa(b.Player + 1)
	coins := strconv.Itoa(b.Coins)
	return "Player " + givingPlayer + " loses " + coins + " coins"
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

func (b BowsersChanceTimeEvent) Question(g *Game) string {
	return "What is the result of Bowser's Chance Time?"
}

func (b BowsersChanceTimeEvent) Type() EventType {
	return ENUM_EVT_TYPE
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
