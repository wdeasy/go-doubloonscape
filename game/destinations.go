package game

import (
    "fmt"
    "math/rand"
    "strings"
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
    game.visitDestination(ATLANTIS_NAME, ATLANTIS_CHANCE, ATLANTIS_DURATION, 2, ATLANTIS_MOD_MAX)	
}

//time modifier
func (game *Game) visitBermuda() {
    game.visitDestination(BERMUDA_NAME, BERMUDA_CHANCE, BERMUDA_DURATION, (BERMUDA_MOD_MAX * -1), BERMUDA_MOD_MAX)
}

//time modifier
func (game *Game) visitDestination(name string, chance int, duration int, lower int, upper int) {

    if _, ok := game.destinations[name]; !ok {
        game.setDestination(name, time.Now(), 0)
        return	
    }

    if time.Now().Before(game.destinations[name].End) {
        return
    }

    if game.destinations[name].Amount != 0 {
        game.destinations[name].Amount = 0
        game.setDestinations()
    }	

    if rand.Intn(100) <= chance {
        amount := RandInt(lower, upper)
        if amount == 0 {
            return
        }

        end := time.Now().Add(time.Minute * time.Duration(duration))        

        game.updateDestination(name, end, amount)
        game.setDestinations()
    }
}

//add the destination info
func (game *Game) setDestination(name string, end time.Time, amount int) {
    var destination storage.Destination

    destination.DB = game.storage.DB
    destination.Name = name
    destination.End = end
    destination.Amount = amount

    game.destinations[destination.Name] = &destination	
}

//update the destination info
func (game *Game) updateDestination(name string, end time.Time, amount int) {
    game.destinations[name].End = end
    game.destinations[name].Amount = amount
}

//generate random int
func RandInt(lower, upper int) int {
    rand.Seed(time.Now().UnixNano())
    rng := upper - lower

    return rand.Intn(rng) + lower
}

//generate random int64
func RandInt64(lower, upper int64) int64 {
    rand.Seed(time.Now().UnixNano())
    rng := upper - lower

    if upper <= 0 {
        return 0
    }

    return rand.Int63n(rng) + lower
}

//update the embed destination info
func (game *Game) setDestinations() {
    game.stats.Destinations = game.destinationsString()
}

//generate destinations string for embed
func (game *Game) destinationsString() (*string) {
    var b strings.Builder

    if time.Now().Before(game.destinations[ATLANTIS_NAME].End) {
        fmt.Fprintf(&b, "` %s  ` ` %s %d `\n", "𝔄𝔱𝔩𝔞𝔫𝔱𝔦𝔰", "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰 𝔪𝔲𝔩𝔱𝔦𝔭𝔩𝔦𝔢𝔡 𝔟𝔶", game.destinations["atlantis"].Amount)
    }

    if time.Now().Before(game.destinations[BERMUDA_NAME].End) {
        fmt.Fprintf(&b, "` %s ` ` %s %d %s `\n", "𝔅𝔢𝔯𝔪𝔲𝔡𝔞", "𝔗𝔦𝔪𝔢 𝔞𝔩𝔱𝔢𝔯𝔢𝔡 𝔟𝔶", game.destinations["bermuda"].Amount, "𝔭𝔢𝔯𝔠𝔢𝔫𝔱")
    }	

    if b.Len() == 0 {
        fmt.Fprintf(&b, "` %s `\n", "𝔗𝔥𝔢 𝔖𝔢𝔳𝔢𝔫 𝔖𝔢𝔞𝔰")               
    }

    String := "**𝔇𝔢𝔰𝔱𝔦𝔫𝔞𝔱𝔦𝔬𝔫𝔰**\n" + b.String()
    
    return &String
}