package mp1

type Event interface {
}

type BranchEvent struct {
	Player int
	Chain  int
	Moves  int
}

type MushroomEvent struct {
	Player int
}
