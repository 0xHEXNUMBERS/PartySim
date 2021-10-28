package mp1

import (
	"fmt"
	"strconv"
)

//Response is a response to any Event.
type Response interface{}

//TODO: Update each event to have Type() method

type EventType int

const (
	//ENUM_EVT_TYPE is the default for responses that are not classified as
	//Ranges, Booleans, Players or ChainSpaces. Responses of this type are typically
	//implement fmt.Stringer to give a verbose definition of what they are,
	//but that is not required.
	ENUM_EVT_TYPE EventType = iota

	//RANGE_EVT_TYPE specifies that responses are integers between a given
	//range. The 0th index of Responses will hold the minimum integer
	//value, while the last index of Responses will hold the maximum
	//integer value.
	RANGE_EVT_TYPE

	//COIN_EVT_TYPE specifies that responses are integers that correspond
	//to a coin amount.
	COIN_EVT_TYPE

	//PLAYER_EVT_TYPE specifies that responses are integers that correspond
	//to player indicies (0 == Player 1, 1 == Player 2, etc.). A response
	//of 4 is only used in Drawable FFA Minigames to indicate that a draw
	//will occur.
	PLAYER_EVT_TYPE

	//MULTIWIN_PLAYER_EVT_TYPE specifies that responses are integers that
	//correspond to a player mask. The nth bit of the integer represents
	//whether player n+1 has won the minigame.
	//Examples:
	//0b0011 --> Players 1 and 2 have won
	//0b1010 --> Players 2 and 4 have won
	//0b0000 --> All players lost
	MULTIWIN_PLAYER_EVT_TYPE

	//CHAINSPACE_EVT_TYPE specifies that responses are chainspaces on the
	//board.
	CHAINSPACE_EVT_TYPE
)

//Event is an action that can be responded to via a Response.
type Event interface {
	//Responses returns a list of all responses that this event can
	//handle.
	Responses() []Response

	//ControllingPlayer returns the player that is responding to the
	//event.
	ControllingPlayer() int

	//Handle handles the current event with the given response onto
	//the given game. Handle must set the Game's NextEvent field.
	Handle(Response, *Game)

	//Type returns what types of responses the caller should expect. Must
	//be one of ENUM_EVT_TYPE, RANGE_EVT_TYPE, or CHAINSPACE_EVT_TYPE.
	Type() EventType

	//Question returns a representation of the struct in question form
	//(e.g. What face did the die land on? Which path will Mario take?
	//Does Yoshi pay 20 coins to perform x action?).
	Question(*Game) string
}

//BranchEvent lets the player decide where to branch off to.
type BranchEvent struct {
	Player int
	Moves  int
	Links  *[]ChainSpace
}

func (b BranchEvent) Type() EventType {
	return CHAINSPACE_EVT_TYPE
}

func (b BranchEvent) Question(g *Game) string {
	return fmt.Sprintf("Which path will the %s take?",
		g.Players[b.Player].Char)
}

//Responses return a slice of landable ChainSpaces that the player can move
//to.
func (b BranchEvent) Responses() []Response {
	ret := []Response{}
	links := *b.Links
	for _, l := range links {
		ret = append(ret, l)
	}
	return ret
}

//Handle moves the player to the selected ChainSpace. The player then
//moves the remaining spaces - 1.
func (b BranchEvent) Handle(r Response, g *Game) {
	newPlayerPos := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = newPlayerPos
	g.MovePlayer(b.Player, b.Moves-1)
}

func (b BranchEvent) ControllingPlayer() int {
	return b.Player
}

//PayRangeEvent allows a player to pay some amount of coins within a
//range. It is mostly contained by other events that need a player to
//pay some amount of coins.
type PayRangeEvent struct {
	Range
	Player int
}

//Handle takes the given number of coins away from the player.
func (p PayRangeEvent) Handle(r Response, g *Game) {
	cost := r.(int)
	g.AwardCoins(p.Player, -cost, false)
}

func (p PayRangeEvent) ControllingPlayer() int {
	return p.Player
}

type MushroomEventResponse int

const (
	RedMushroom MushroomEventResponse = iota
	PoisonMushroom
)

func (m MushroomEventResponse) String() string {
	switch m {
	case RedMushroom:
		return "Red Mushroom"
	case PoisonMushroom:
		return "Poison Mushroom"
	}
	return ""
}

//MushroomEvent occurs when a player lands on a Mushroom Space.
type MushroomEvent struct {
	Player int
}

func (m MushroomEvent) Question(g *Game) string {
	return fmt.Sprintf("What mushroom did %s recieve?",
		g.Players[m.Player].Char)
}

func (m MushroomEvent) Type() EventType {
	return ENUM_EVT_TYPE
}

func (m MushroomEvent) Responses() []Response {
	return []Response{RedMushroom, PoisonMushroom}
}

//Handle sets next players turn. If r == true, then player m.Player goes
//again. If r == false, then m.Player's SkipTurn flag is set before ending
//their turn.
func (m MushroomEvent) Handle(r Response, g *Game) {
	redMushroom := r.(MushroomEventResponse)
	if redMushroom == RedMushroom {
		g.SetDiceBlock()
		return
	}
	g.Players[m.Player].SkipTurn = true
	g.Players[m.Player].LastSpaceType = Mushroom
	g.EndCharacterTurn()
}

func (m MushroomEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//BooCoinsEvent handles the transfer of coins from one player to another
//when a player decides to steal coins via passing by a Boo space.
type BooCoinsEvent struct {
	PayRangeEvent
	RecvPlayer int
	Moves      int
}

func (b BooCoinsEvent) Question(g *Game) string {
	return fmt.Sprintf("How many coins will %s steal from %s",
		g.Players[b.RecvPlayer].Char,
		g.Players[b.PayRangeEvent.Player].Char)
}

func (b BooCoinsEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle transfers r coins from the giving player in PayRangeEvent to the
//receiving player.
func (b BooCoinsEvent) Handle(r Response, g *Game) {
	b.PayRangeEvent.Handle(r, g)
	g.AwardCoins(b.RecvPlayer, r.(int), false)

	if b.Moves != 0 {
		g.MovePlayer(b.RecvPlayer, b.Moves)
	} else {
		g.EndCharacterTurn()
	}
}

//BooEvent lets a passing player decide what action to take when passing
//a Boo space.
type BooEvent struct {
	Player  int
	Players [4]Player
	Moves   int //No call to MovePlayer on 0
	Coins   int
}

//BooStealAction describes an action a player passing Boo may take.
type BooStealAction struct {
	RecvPlayer   int
	GivingPlayer int
	Star         bool
}

func (b BooStealAction) String() string {
	givingPlayer := strconv.Itoa(b.GivingPlayer + 1)
	if b.Star {
		return "Steal 1 star from player " + givingPlayer + " for 50 coins"
	} else {
		return "Steal coins from player " + givingPlayer
	}
}

func (b BooEvent) Question(g *Game) string {
	return fmt.Sprintf("What will %s do with Boo?",
		g.Players[b.Player].Char)
}

func (b BooEvent) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of BooStealActions that b.Player can take.
func (b BooEvent) Responses() []Response {
	res := make([]Response, 0)
	if b.Coins >= 50 {
		for i := 0; i < 4; i++ {
			if i == b.Player {
				continue
			}
			if b.Players[i].Stars > 0 {
				res = append(res, BooStealAction{b.Player, i, true})
			}
		}
	}
	for i := 0; i < 4; i++ {
		if i == b.Player {
			continue
		}
		if b.Players[i].Coins > 0 {
			res = append(res, BooStealAction{b.Player, i, false})
		}
	}
	return res
}

//Handle applies the BooStealAction r for b.Player. After applying r, the
//player moves their remaining spaces.
func (b BooEvent) Handle(r Response, g *Game) {
	steal := r.(BooStealAction)
	if steal.Star {
		g.AwardCoins(steal.RecvPlayer, -50, false)
		g.Players[steal.GivingPlayer].Stars--
	} else {
		maxCoins := 15
		if b.Players[steal.GivingPlayer].Coins <= maxCoins {
			maxCoins = b.Players[steal.GivingPlayer].Coins
		}
		g.NextEvent = BooCoinsEvent{
			PayRangeEvent{Range{1, maxCoins}, steal.GivingPlayer},
			steal.RecvPlayer,
			b.Moves,
		}
		return
	}
	if b.Moves != 0 {
		g.MovePlayer(b.Player, b.Moves)
	} else {
		g.EndCharacterTurn()
	}
}

func (b BooEvent) ControllingPlayer() int {
	return b.Player
}

//DeterminePlayerTeamEvent handles deciding which minigame team a player
//is if said player landed on a *green* space.
type DeterminePlayerTeamEvent struct {
	Player int
}

func (d DeterminePlayerTeamEvent) Question(g *Game) string {
	return fmt.Sprintf("What team was %s chosen to be on?",
		g.Players[d.Player].Char)
}

func (d DeterminePlayerTeamEvent) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns BlueTeam and RedTeam, the two available teams a player
//can be on.
func (d DeterminePlayerTeamEvent) Responses() []Response {
	return []Response{BlueTeam, RedTeam}
}

//Handle sets d.Player's team. If r is true, d.Player's team is blue. If r
//is false, d.Player's team is Red.
func (d DeterminePlayerTeamEvent) Handle(r Response, g *Game) {
	team := r.(MinigameTeam)

	if team == BlueTeam {
		g.Players[d.Player].LastSpaceType = Blue
	} else {
		g.Players[d.Player].LastSpaceType = Red
	}
	g.FindGreenPlayer()
}

func (d DeterminePlayerTeamEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Range is a partial event that generates a range from [Min,Max]
//that the CPU player can respond to. It is mostly used to generate the
//[Min,Max] range for other events.
type Range struct {
	Min int
	Max int
}

func NewRange(min, max int) []Response {
	return Range{min, max}.Responses()
}

func (r Range) Type() EventType {
	return RANGE_EVT_TYPE
}

//Responses returns a list of ints from [c.Min,c.Max].
func (r Range) Responses() []Response {
	var ret []Response
	for i := r.Min; i <= r.Max; i++ {
		ret = append(ret, i)
	}
	return ret
}
