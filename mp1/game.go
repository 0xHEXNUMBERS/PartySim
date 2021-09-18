package mp1

//GameConfig holds the configuration settings of the current game.
type GameConfig struct {
	MaxTurns     uint8
	NoBonusStars bool
	NoKoopa      bool
	NoBoo        bool
	RedDice      bool
	BlueDice     bool
	WarpDice     bool
	EventsDice   bool
}

//Game is the structure that holds all game information.
type Game struct {
	Board
	Config        GameConfig
	StarSpaces    StarData
	Players       [4]Player
	Turn          uint8
	CurrentPlayer int
	NextEvent     Event

	//Every 10 passes, Koopa rewards 20 coins to the passing player.
	KoopaPasses int
}

//Responses returns the valid responses for the next event.
func (g *Game) Responses() []Response {
	return g.NextEvent.Responses()
}

//HandleEvent executes the next event using the given Response r.
func (g *Game) HandleEvent(r Response) {
	g.NextEvent.Handle(r, g)
}

//SetDiceBlock looks at the GameConfig to see if there are any special
//dice in play. If there are, the next Event is set to pick a dice block.
//Otherwise, the next Event is set to the normal dice block.
func (g *Game) SetDiceBlock() {
	if g.Turn != 0 && (g.Config.RedDice || g.Config.BlueDice || g.Config.WarpDice || g.Config.EventsDice) {
		g.NextEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
	} else {
		g.NextEvent = NormalDiceBlock{Range{1, 10}, g.CurrentPlayer}
	}
}

//InitializeGame returns a new game given a Board and a GameConfig.
//This function goes through 2 main steps.
//1. It finds the start space and sets all player positions to that space.
//If there is no start space, player positions are set to ChainSpace{0, 0}.
//2. It looks for all of the star spaces and intializes StarData.
func InitializeGame(b Board, config GameConfig) *Game {
	g := &Game{
		Board: b,
	}
	var chainSpaces []ChainSpace
	var count uint8
	var startSpace ChainSpace
	for ci, chain := range *b.Chains {
		for si, space := range chain {
			if space.Type == Star {
				chainSpaces = append(
					chainSpaces,
					ChainSpace{ci, si},
				)
				count++
			} else if space.Type == Start {
				startSpace = ChainSpace{ci, si}
			}
		}
	}
	if count > 0 {
		g.StarSpaces.IndexToPosition = &chainSpaces
		g.StarSpaces.StarSpaceCount = count
	}
	g.Config = config

	for i := 0; i < len(g.Players); i++ {
		g.Players[i].CurrentSpace = startSpace
		g.Players[i].LastSpaceType = Start
		g.Players[i].Coins = 10
	}

	if g.StarSpaces.StarSpaceCount == 1 { //Go ahead and set star pos
		g.StarSpaces.CurrentStarSpace = (*g.StarSpaces.IndexToPosition)[0]
	}
	if g.StarSpaces.StarSpaceCount <= 1 {
		g.SetDiceBlock()
	} else {
		g.NextEvent = StarLocationEvent{g.StarSpaces, 0, 0}
	}
	return g
}

//LastFiveTurns returns true if the game is in its' final 5 turns.
func (g *Game) LastFiveTurns() bool {
	return g.Config.MaxTurns-g.Turn <= 5
}

//AwardCoins awards a player with coins. It handles min/maxing
//of coins values and Coin/Minigame bonus stars. It returns the
//number of coins the player received.
func (g *Game) AwardCoins(player, coins int, minigame bool) int {
	coins0 := g.Players[player].Coins
	g.Players[player].Coins += coins
	g.Players[player].Coins = max(g.Players[player].Coins, 0)
	coinsGiven := g.Players[player].Coins - coins0
	if minigame {
		g.Players[player].MinigameCoins += coins
	}
	g.Players[player].MaxCoins = max(
		g.Players[player].MaxCoins,
		g.Players[player].Coins,
	)
	return coinsGiven
}

//GiveCoins transfers coins from one player to another.
func (g *Game) GiveCoins(givingPlayer, takingPlayer, coins int, minigame bool) {
	coinsTaken := -g.AwardCoins(givingPlayer, -coins, minigame)
	g.AwardCoins(takingPlayer, coinsTaken, minigame)
}

//CheckLinks looks through the boards linkage system to determine if a
//player needs to make a decision where to branch off to, or setting the
//player's new position if there's <=1 link at the player's current chain.
func (g *Game) CheckLinks(player, chain, moves int) (branch bool) {
	//Check for links
	if g.Board.Links != nil {
		chainLinks := *(g.Board.Links)
		linksPtr, ok := chainLinks[chain]
		if !ok {
			g.Players[player].CurrentSpace.Space = 0
			return
		}
		links := *linksPtr
		switch len(links) {
		case 1:
			g.Players[player].CurrentSpace = links[0]
		default:
			g.NextEvent = BranchEvent{
				player,
				moves - 1,
				linksPtr,
			}
			return true
		}
	}
	return false
}

//ActivateSpace performs the action of a player landing on a space.
func (g *Game) ActivateSpace(player int) {
	//Activate Space
	chains := *g.Board.Chains
	playerPos := g.Players[player].CurrentSpace
	curSpace := chains[playerPos.Chain][playerPos.Space]
	//Star space requires preprocessing
	if curSpace.Type == Star {
		starIndex := g.StarSpaces.GetIndex(playerPos)
		if g.StarSpaces.AbsoluteVisited&(1<<starIndex) > 0 {
			g.Players[player].LastSpaceType = Chance
		} else {
			g.Players[player].LastSpaceType = Blue
		}
	}
	//Perform space action
	switch g.Players[player].LastSpaceType {
	case Invisible:
		//Stopping Event should set LastSpaceType
		curSpace.StoppingEvent(g, player)
		g.ActivateSpace(player)
	case Blue:
		if g.LastFiveTurns() {
			g.AwardCoins(player, 6, false)
		} else {
			g.AwardCoins(player, 3, false)
		}
		g.EndCharacterTurn()
	case Red:
		if g.LastFiveTurns() {
			g.AwardCoins(player, -6, false)
		} else {
			g.AwardCoins(player, -3, false)
		}
		g.EndCharacterTurn()
	case Mushroom:
		g.NextEvent = MushroomEvent{Boolean{}, player}
	case Happening:
		g.Players[player].HappeningCount++
		g.NextEvent = nil
		curSpace.StoppingEvent(g, player)
		if g.NextEvent == nil {
			g.EndCharacterTurn()
		}
	case Bowser:
		g.PreBowserCheck(player)
	case MinigameSpace:
		g.NextEvent = Minigame1PSelector{player}
	case Chance:
		g.NextEvent = ChanceTime{Player: player}
	}
}

//MovePlayer moves the player x many spaces through the board. It handles
//branching and passing events.
func (g *Game) MovePlayer(playerIdx, moves int) {
	chains := *g.Board.Chains
	playerPos := &g.Players[playerIdx].CurrentSpace
	for moves > 0 {
		playerPos.Space++
		if playerPos.Space >= len(chains[playerPos.Chain]) {
			if g.CheckLinks(playerIdx, playerPos.Chain, moves) {
				return
			}
		}
		curSpace := chains[playerPos.Chain][playerPos.Space]
		switch curSpace.Type {
		case Invisible:
			if curSpace.PassingEvent != nil {
				g.NextEvent = nil
				moves = curSpace.PassingEvent(g, playerIdx, moves)
				if g.NextEvent != nil {
					return
				}
			} else {
				moves--
			}
		case Start:
			if !g.Config.NoKoopa {
				g.KoopaPasses++
				if g.KoopaPasses%10 == 0 {
					g.AwardCoins(playerIdx, 20, false)
				} else {
					g.AwardCoins(playerIdx, 10, false)
				}
			}
		case Star:
			if *playerPos == g.StarSpaces.CurrentStarSpace &&
				g.Players[playerIdx].Coins >= 20 {
				g.Players[playerIdx].Stars++
				g.AwardCoins(playerIdx, -20, false)
				if g.StarSpaces.StarSpaceCount > 1 {
					g.NextEvent = StarLocationEvent{
						g.StarSpaces,
						playerIdx,
						moves,
					}
					return
				}
			}
			moves--
		case Boo:
			if !g.Config.NoBoo {
				booEvt := BooEvent{
					playerIdx,
					g.Players,
					moves,
					g.Players[playerIdx].Coins,
				}
				if len(booEvt.Responses()) != 0 {
					g.NextEvent = booEvt
					return
				}
			}
		case BogusItem:
			g.AwardCoins(playerIdx, -g.Board.BowserCoins, false)
		default:
			moves--
		}
	}
	curSpace := chains[playerPos.Chain][playerPos.Space]
	g.Players[playerIdx].LastSpaceType = curSpace.Type
	if g.Config.EventsDice &&
		(curSpace.Type == Blue || curSpace.HiddenBlock) {
		g.NextEvent = HiddenBlockEvent{playerIdx}
		return
	}
	g.ActivateSpace(playerIdx)
}

//AwardBonusStars looks through the players' aggregated statistics to
//award bonus stars.
func (g *Game) AwardBonusStars() {
	if g.Config.NoBonusStars {
		return
	}

	maxCoins := g.Players[0].MaxCoins
	for i := 1; i < 4; i++ {
		maxCoins = max(maxCoins, g.Players[i].MaxCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MaxCoins == maxCoins {
			g.Players[i].Stars++
		}
	}

	maxMinigameCoins := g.Players[0].MinigameCoins
	for i := 1; i < 4; i++ {
		maxMinigameCoins = max(maxMinigameCoins, g.Players[i].MinigameCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MinigameCoins == maxMinigameCoins {
			g.Players[i].Stars++
		}
	}

	maxHappening := g.Players[0].HappeningCount
	for i := 1; i < 4; i++ {
		maxHappening = max(maxHappening, g.Players[i].HappeningCount)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].HappeningCount == maxHappening {
			g.Players[i].Stars++
		}
	}
}

//Winners returns a list of the winning player indexes at the current game
//state.
func (g *Game) Winners() []int {
	maxStarHolders := []int{}
	maxStars := g.Players[0].Stars
	for i := 1; i < 4; i++ {
		maxStars = max(maxStars, g.Players[i].Stars)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].Stars == maxStars {
			maxStarHolders = append(maxStarHolders, i)
		}
	}

	winners := []int{}
	maxCoins := g.Players[maxStarHolders[0]].Coins
	for i := 1; i < len(maxStarHolders); i++ {
		maxCoins = max(maxCoins, g.Players[maxStarHolders[i]].Coins)
	}
	for i := 0; i < len(maxStarHolders); i++ {
		if g.Players[maxStarHolders[i]].Coins == maxCoins {
			winners = append(winners, maxStarHolders[i])
		}
	}
	return winners
}

//EndGameTurn ends the game turn. It handles awarding bonus stars at the
//end of the game, skipping player 0's turn in case of poison mushroom,
//and setting the next diceblock.
func (g *Game) EndGameTurn() {
	g.Turn++
	if g.Turn == g.Config.MaxTurns {
		g.AwardBonusStars()
		//Game is over, no more events
		g.NextEvent = nil
	} else {
		if g.Players[g.CurrentPlayer].SkipTurn {
			g.Players[g.CurrentPlayer].SkipTurn = false
			g.EndCharacterTurn()
			return
		}
		g.SetDiceBlock()
	}
}

//StartMinigamePrep starts preparation for the next end of turn minigame.
func (g *Game) StartMinigamePrep() {
	g.FindGreenPlayer()
}

//EndCharacterTurn handles events that occur at the end of a character's
//turn. It handles skipping player turns if the next player received a
//poison mushroom, Starting minigame preparation if player 3 just
//finished, and calling the board's specifc end of turn event.
func (g *Game) EndCharacterTurn() {
	if g.Board.EndCharacterTurnEvent != nil {
		g.Board.EndCharacterTurnEvent(g, g.CurrentPlayer)
	}
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	if g.CurrentPlayer == 0 {
		g.StartMinigamePrep()
		return
	}
	if g.Players[g.CurrentPlayer].SkipTurn {
		g.Players[g.CurrentPlayer].SkipTurn = false
		g.EndCharacterTurn()
		return
	}
	g.SetDiceBlock()
}
