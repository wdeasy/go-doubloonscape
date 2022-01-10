package game

import (
    "math/rand"
    "time"

    "github.com/wdeasy/go-doubloonscape/storage"
)

//iterate over destinations
func (game *Game) visitDestinations() {
    game.visitAtlantis()
    game.visitBermuda()
}

//gold modifier
func (game *Game) visitAtlantis() {
    game.visitDestination("atlantis", ATLANTIS_CHANCE, ATLANTIS_DURATION, 1, ATLANTIS_MOD_MAX)	
}

//time modifier
func (game *Game) visitBermuda() {
    game.visitDestination("bermuda", BERMUDA_CHANCE, BERMUDA_DURATION, (BERMUDA_MOD_MAX * -1), BERMUDA_MOD_MAX)
}

//time modifier
func (game *Game) visitDestination(name string, chance int, duration int, lower int, upper int) {
    if _, ok := game.destinations[name]; !ok {
        game.setDestination(name, time.Now(), 0)	
    }

    if time.Now().Before(game.destinations[name].End) {
        return
    }

    if time.Now().After(game.destinations[name].End) && game.destinations[name].Amount != 0 {
        game.setDestination(name, time.Now(), 0)
        return
    }	

    if rand.Intn(100) <= chance {
        end := time.Now().Add(time.Minute * time.Duration(duration))
        amount := RandInt(lower, upper)

        game.setDestination(name, end, amount)
    }
}

//update the destination info
func (game *Game) setDestination(name string, end time.Time, amount int) {
    var destination storage.Destination

    destination.Name = name
    destination.End = time.Now()
    destination.Amount = 0

    game.destinations[destination.Name] = destination	
}

//generate random int
func RandInt(lower, upper int) int {
    rand.Seed(time.Now().UnixNano())
    rng := upper - lower
    return rand.Intn(rng) + lower
}

