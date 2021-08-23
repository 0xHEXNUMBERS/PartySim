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

Documentation for any simulator is available at https://pkg.go.dev/github.com/0xhexnumbers/partysim/[game abbr].

### Mario Party 1

https://pkg.go.dev/github.com/0xhexnumbers/partysim/mp1

## Getting Started

There are 2 ways to download this package, `git` and `go get`.

### git

To download the core for all games, use `git clone https://github.com/0xhexnumbers/partysim`.

### go get

To download the core for a particular game, use `go get github.com/0xhexnumbers/partysim/[game abbr]`. For example: to download the core for Mario Party 1, use `go get github.com/0xhexnumbers/partysim/mp1`.

## Control Flow

The simulator runs events, and each event sets the next event to be run. An event is executed with a response that the event can accept. For each event, there are at least 2 possible responses to said event.

## Sample Code

This sample code simulates a random game of Mario Party 1 on Eternal Star.

```go
func SimulateGame(r *rand.Rand) {
        g := mp1.InitializeGame(mp1.ES, mp1.GameConfig{MaxTurns: 20})
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
 
        fmt.Println()
        fmt.Printf("P1: %#v\n", g.Players[0])
        fmt.Printf("P2: %#v\n", g.Players[1])
        fmt.Printf("P3: %#v\n", g.Players[2])
        fmt.Printf("P4: %#v\n", g.Players[3])
 }
```

## Bug Report

If any bugs or crashes are found in any simulator, open an Github Issue and describe the bug or crash in detail.