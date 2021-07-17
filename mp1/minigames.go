package mp1

type MinigameAwards [4]int

var MinigameRewardsFFA = []Response{
	MinigameAwards{10, 0, 0, 0},
	MinigameAwards{0, 10, 0, 0},
	MinigameAwards{0, 0, 10, 0},
	MinigameAwards{0, 0, 0, 10},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards2V2 = []Response{
	MinigameAwards{10, 10, 0, 0},
	MinigameAwards{0, 0, 10, 10},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards1V3 = []Response{
	MinigameAwards{15, -5, -5, -5},
	MinigameAwards{-15, 5, 5, 5},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards1P = []Response{
	MinigameAwards{5, 0, 0, 0},
	MinigameAwards{0, 0, 0, 0},
}

type MinigameType int

const (
	MinigameFFA MinigameType = iota
	Minigame2V2
	Minigame1V3
	Minigame1P
)

type MinigameEvent struct {
	//Player IDs
	//FFA -> [Player0, Player1, Player2, Player3]
	//1V3 -> [Team1, Team2, Team2, Team2]
	//2V2 -> [Team1, Team1, Team2, Team2]
	//1P  -> [Team1, nil, nil, nil]
	PlayerIDs [4]int
	Type      MinigameType
}

func (m MinigameEvent) Responses() []Response {
	switch m.Type {
	case MinigameFFA:
		return MinigameRewardsFFA
	case Minigame2V2:
		return MinigameRewards2V2
	case Minigame1V3:
		return MinigameRewards1V3
	case Minigame1P:
		return MinigameRewards1P
	}
	//Unreachable
	return nil
}

func (m MinigameEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (m MinigameEvent) Handle(r Response, g Game) Game {
	g.ExtraEvent = nil
	awards := r.(MinigameAwards)
	for i, player := range m.PlayerIDs {
		g = AwardCoins(g, player, awards[i], true)
	}
	return g
}

func GetMinigame(g Game) MinigameEvent {
	var blueTeam []int
	var redTeam []int
	for i, p := range g.Players {
		if p.LastSpaceType == Blue {
			blueTeam = append(blueTeam, i)
		} else if p.LastSpaceType == Red {
			redTeam = append(redTeam, i)
		}
	}

	var minigameType MinigameType
	switch len(blueTeam) {
	case 0, 4:
		minigameType = MinigameFFA
	case 1, 3:
		minigameType = Minigame1V3
	case 2:
		minigameType = Minigame2V2
	}

	minigame := MinigameEvent{Type: minigameType}
	var players []int
	if len(redTeam) == 1 { //Put 1 person team in front
		players = append(redTeam, blueTeam...)
	} else {
		players = append(blueTeam, redTeam...)
	}
	for i := range blueTeam {
		minigame.PlayerIDs[i] = players[i]
	}
	return minigame
}

func FindGreenPlayer(g Game) Event {
	for i, p := range g.Players {
		if p.LastSpaceType != Blue && p.LastSpaceType != Red {
			return DeterminePlayerTeamEvent{
				Player: i,
			}
		}
	}
	return nil
}
