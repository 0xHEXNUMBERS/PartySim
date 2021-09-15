# PartySim

PartySim is a simulator for the Mario Party series of games.

## Table of Contents

### [Documentation](#documentation)

- [Mario Party 1](#mario-party-1)

### [Getting Started](#getting-started)

- [git](#git)
- [go get](#go-get)
  
### [Control Flow](#control-flow)

### [Sample Code](#sample-code)

### [Bug Report](#bug-report)

## Documentation

Documentation for any game simulator is available at https://pkg.go.dev/github.com/0xhexnumbers/partysim/[game abbr].

Documentation for any board implementation is available at https://pkg.go.dev/github.com/0xhexnumbers/partysim/[game abbr]/board.

### Mario Party 1

https://pkg.go.dev/github.com/0xhexnumbers/partysim/mp1

https://pkg.go.dev/github.com/0xhexnumbers/partysim/mp1/board

## Getting Started

There are 2 ways to download this package, `git` and `go get`.

### git

To download the latest core for all games, use `git clone https://github.com/0xhexnumbers/partysim`.

### go get

To download the stable core for all games, use `go get github.com/0xhexnumbers/partysim@v0.0.1`.

## Control Flow

The simulator runs events, and each event sets the next event to be run. An event is executed with a response that the event can accept. For each event, there are at least 2 possible responses to said event.

## Sample Code

This sample code simulates a random game of Mario Party 1 on Eternal Star.

```go
package main

import (
	"github.com/0xhexnumbers/mp1"
	"github.com/0xhexnumbers/mp1/board"
)

func PrintPlayer(p mp1.Player) {
        fmt.Printf("%s: %d Stars -- %d Coins\n",
		p.Char, p.Stars, p.Coins)
}

func SimulateGame(r *rand.Rand) {
        g := mp1.InitializeGame(board.ES, mp1.GameConfig{MaxTurns: 20})
        g.Players[0].Char = "Mario"
        g.Players[1].Char = "Luigi"
        g.Players[2].Char = "Peach"
        g.Players[3].Char = "Yoshi"
 
        for g.NextEvent != nil {
                res := g.NextEvent.Responses()
                randInt := r.Intn(len(res))
                //fmt.Printf("%#v\n%#v\n\n", g.ExtraEvent, res[randInt])
                g.HandleEvent(res[randInt])
        }

	//Print results of match 
        fmt.Println()
	PrintPlayer(g.Players[0])
	PrintPlayer(g.Players[1])
	PrintPlayer(g.Players[2])
	PrintPlayer(g.Players[3])
 }
```

## Bug Report

If any bugs or crashes are found in any simulator, open an Github Issue and describe the bug or crash in detail.
