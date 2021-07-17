package mp1

type BowserEvent struct {
	Player int
}

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

func (b BowserEvent) Handle(r Response, g Game) Game {
	choice := r.(BowserResponse)
	//TODO: BowserChanceTime
	switch choice {
	case CoinsForBowser:
		maxCoins := min(g.Players[b.Player].Coins, 30)
		minCoins := min(g.Players[b.Player].Coins, 10)
		g.ExtraEvent = PayRangeEvent{b.Player, minCoins, maxCoins, 0}
	case BowserBalloonBurst:
		g.ExtraEvent = BowserBalloonBurstEvent{}
	case BowsersFaceLift:
		g.ExtraEvent = BowsersFaceLiftEvent{b.Player}
	case BowsersTugoWar:
		g.ExtraEvent = BowsersTugoWarEvent{b.Player}
	case BashnCash:
		g.ExtraEvent = BashnCashEvent{b.Player, g.Players[b.Player].Coins}
	case BowserRevolution:
		coins := 0
		for i := range g.Players {
			coins += g.Players[i].Coins
		}
		coins /= 4
		for i := range g.Players {
			g.Players[i].Coins = coins
		}
	}
	return g
}

type BowserBalloonBurstEvent struct{}

type BowserBalloonBurstResult int

const (
	BBBDraw BowserBalloonBurstResult = iota
	BBBP1Win
	BBBP2Win
	BBBP3Win
	BBBP4Win
)

func GetBowserMinigameCoinLoss(turn uint) int {
	if turn <= 9 {
		return 10
	} else if turn <= 19 {
		return 20
	} else if turn <= 29 {
		return 30
	}
	return 40
}

var BBBResults = []Response{
	BBBDraw,
	BBBP1Win,
	BBBP2Win,
	BBBP3Win,
	BBBP4Win,
}

func (b BowserBalloonBurstEvent) Responses() []Response {
	return BBBResults
}

func (b BowserBalloonBurstEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowserBalloonBurstEvent) Handle(r Response, g Game) Game {
	results := r.(BowserBalloonBurstResult)
	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	switch results {
	case BBBDraw:
		for p := range g.Players {
			g = AwardCoins(g, p, -20, true)
		}
	case BBBP1Win:
		g = AwardCoins(g, 1, coinLoss, true)
		g = AwardCoins(g, 2, coinLoss, true)
		g = AwardCoins(g, 3, coinLoss, true)
	case BBBP2Win:
		g = AwardCoins(g, 0, coinLoss, true)
		g = AwardCoins(g, 2, coinLoss, true)
		g = AwardCoins(g, 3, coinLoss, true)
	case BBBP3Win:
		g = AwardCoins(g, 0, coinLoss, true)
		g = AwardCoins(g, 1, coinLoss, true)
		g = AwardCoins(g, 3, coinLoss, true)
	case BBBP4Win:
		g = AwardCoins(g, 0, coinLoss, true)
		g = AwardCoins(g, 1, coinLoss, true)
		g = AwardCoins(g, 2, coinLoss, true)
	}
	return g
}

type BowsersFaceLiftEvent struct {
	Player int
}

func (b BowsersFaceLiftEvent) Responses() []Response {
	return CPURangeEvent{1, 15}.Responses()
}

func (b BowsersFaceLiftEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersFaceLiftEvent) Handle(r Response, g Game) Game {
	results := r.(int)
	if results == 15 { //All players won
		g = AwardCoins(g, b.Player, -50, true)
		return g
	}

	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	for p := range g.Players {
		if results&(1<<p) == 0 {
			g = AwardCoins(g, p, coinLoss, true)
		}
	}
	return g
}

type BowsersTugoWarEvent struct {
	Player int
}

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

func (b BowsersTugoWarEvent) Responses() []Response {
	return BTWResults
}

func (b BowsersTugoWarEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersTugoWarEvent) Handle(r Response, g Game) Game {
	results := r.(BowsersTugoWarResult)
	switch results {
	case BTWDraw:
		for p := range g.Players {
			g = AwardCoins(g, p, -30, true)
		}
	case BTW1TWin:
		coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
		for p := range g.Players {
			if p != b.Player {
				g = AwardCoins(g, p, coinLoss, true)
			}
		}
	case BTW3TWin:
		g = AwardCoins(g, b.Player, -10, true)
	}
	return g
}

type BashnCashEvent struct {
	Player int
	Coins  int
}

func (b BashnCashEvent) Responses() []Response {
	max := b.Coins / 5
	max += b.Coins % 5
	return CPURangeEvent{1, max}.Responses()
}

func (b BashnCashEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BashnCashEvent) Handle(r Response, g Game) Game {
	timesHit := r.(int)
	coinsLost := 0
	if b.Coins/5 < timesHit {
		coinsLost += b.Coins - (b.Coins % 5)
		timesHit -= b.Coins / 5
		coinsLost += timesHit
	} else {
		coinsLost += timesHit * 5
	}
	g = AwardCoins(g, b.Player, -coinsLost, true)
	return g
}
