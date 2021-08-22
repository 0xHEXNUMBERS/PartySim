package mp1

//SpaceType is an enum type for various Spaces.
type SpaceType int

const (
	Invisible SpaceType = iota
	Blue
	Red
	MinigameSpace
	Happening
	Star
	Chance
	Start
	Mushroom
	Bowser
	BogusItem
	Boo
)

//Space is a physical space on the board that Players can land on and/or
//pass by.
type Space struct {
	Type SpaceType

	//Whether a Hidden Block can appear on this space or not
	HiddenBlock bool

	//For Invisible/Happening Spaces, gets called when a player lands
	//on this space. If SpaceType == Invisible, the player's
	//LastSpaceType needs to be set to a known landable space type.
	StoppingEvent func(game *Game, player int)

	//For Invisible Spaces, gets called when a player is moving through
	//this space. Value returned is the number of moves the simulation
	//needs to process before ending the Player's movement.
	PassingEvent func(game *Game, player, moves int) int
}

//Chain is a sequence of non-branching Spaces
type Chain []Space

//ChainSpace is an index to the board of chains
type ChainSpace struct {
	Chain int
	Space int
}

//ExtraBoardData is any *comparable* piece of data that the Board holds
//onto. The engine does not manipulate this data directly, but board
//specific function calls may manipulate this data.
type ExtraBoardData interface{}

type Movement struct {
	Player int
	Moves  int
	Skip   bool
}

//Board holds all data specifc to an MP1 board.
type Board struct {
	//Chains is a list  of chains on the board.
	Chains *[]Chain

	//Links is a linking between the end of each chain to the
	//ChainSpace they link to.
	Links *map[int]*[]ChainSpace

	//BowserCoins is the amount of coins Bowser takes from a player when
	//passing by a BogusItem space.
	BowserCoins int

	//Data holds the board specific data.
	Data ExtraBoardData

	//EndCharacterTurnEvent is a function that gets called at the end
	//of every player's turn.
	EndCharacterTurnEvent func(game *Game, player int)
}
