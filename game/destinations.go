package game

import (
    "fmt"
    "math/rand"
    "strings"
    "time"

    "github.com/wdeasy/go-doubloonscape/storage"
)

func (game *Game) getDestinations() ([]string) {
    return []string{ATLANTIS_NAME, BERMUDA_NAME}
}

//iterate over destinations
func (game *Game) visitDestinations() {
    for _, destination := range game.destinations {
        printLog(fmt.Sprintf("visiting: %s\n", destination.Name))
        game.visitDestination(destination)
    }    
}

//time modifier
func (game *Game) visitDestination(destination *storage.Destination) {
    if time.Now().Before(destination.End) {
        printLog(fmt.Sprintf("time is before end: %s\n", destination.Name))
        return
    }

    if destination.Amount != 0 {
        printLog(fmt.Sprintf("amount is not zero: %s\n", destination.Name))
        destination.Amount = 0
        game.setDestinations()
    }	

    if rand.Intn(destination.Max) <= destination.Chance {
        printLog(fmt.Sprintf("passed the roll: %s\n", destination.Name))
        amount := RandInt(destination.Low, destination.ModMax)
        if amount == 0 {
            printLog(fmt.Sprintf("amount == zero: %s\n", destination.Name))
            return
        }

        end := time.Now().Add(time.Minute * time.Duration(destination.Duration))        

        game.updateDestination(*destination, end, amount)
        game.setDestinations()
    }
}

//add the destination info
func (game *Game) setDestination(name string) {
    var destination storage.Destination

    destination.DB       = game.storage.DB
    destination.Name     = name
    destination.End      = time.Now()
    destination.Amount   = 0
    destination.Chance   = game.getDestinationChance(name)
    destination.Low      = game.getDestinationLow(name)
    destination.Max      = game.getDestinationMax(name)
    destination.ModMax   = game.getDestinationModMax(name)
    destination.Duration = game.getDestinationDuration(name)
    

    game.destinations[destination.Name] = &destination	
}

//update the destination info
func (game *Game) updateDestination(destination storage.Destination, end time.Time, amount int) {
    destination.End = end
    destination.Amount = amount
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
        fmt.Fprintf(&b, "` %s  ` ` %s %d `\n", "𝔄𝔱𝔩𝔞𝔫𝔱𝔦𝔰", "𝔇𝔬𝔲𝔟𝔩𝔬𝔬𝔫𝔰 𝔪𝔲𝔩𝔱𝔦𝔭𝔩𝔦𝔢𝔡 𝔟𝔶", game.destinations[ATLANTIS_NAME].Amount)
    }

    if time.Now().Before(game.destinations[BERMUDA_NAME].End) {
        fmt.Fprintf(&b, "` %s ` ` %s %d %s `\n", "𝔅𝔢𝔯𝔪𝔲𝔡𝔞", "𝔗𝔦𝔪𝔢 𝔞𝔩𝔱𝔢𝔯𝔢𝔡 𝔟𝔶", game.destinations[BERMUDA_NAME].Amount, "𝔭𝔢𝔯𝔠𝔢𝔫𝔱")
    }	

    if b.Len() == 0 {
        fmt.Fprintf(&b, "` %s `\n", "𝔗𝔥𝔢 𝔖𝔢𝔳𝔢𝔫 𝔖𝔢𝔞𝔰")               
    }

    String := "**𝔇𝔢𝔰𝔱𝔦𝔫𝔞𝔱𝔦𝔬𝔫𝔰**\n" + b.String()
    
    printLog(fmt.Sprintf("destination string: %s\n", String))
    return &String
}

//update events from constants
func (game *Game) updateDestinations() {
    for _, destination := range game.getDestinations() {
        if _, ok := game.destinations[destination]; !ok {
            game.setDestination(destination)
            continue
        }
        
        game.destinations[destination].Chance   = game.getDestinationChance(destination)
        game.destinations[destination].Max      = game.getDestinationMax(destination)
        game.destinations[destination].ModMax   = game.getDestinationModMax(destination)
        game.destinations[destination].Duration = game.getDestinationDuration(destination)        
    }

    game.storage.SaveEvents(game.events)    
}

//get the appropriate destination chance from constants
func (game *Game) getDestinationChance(name string) (int) {
    switch name{
    case ATLANTIS_NAME:
        return ATLANTIS_CHANCE
    case BERMUDA_NAME:
        return BERMUDA_CHANCE  
    default:
        return DEFAULT_CHANCE
    }   
}

//get the appropriate destination max from constants
func (game *Game) getDestinationLow(name string) (int) {
    switch name{
    case ATLANTIS_NAME:
        return ATLANTIS_LOW
    case BERMUDA_NAME:
        return BERMUDA_LOW 
    default:
        return DEFAULT_LOW
    }   
}

//get the appropriate destination max from constants
func (game *Game) getDestinationMax(name string) (int) {
    switch name{
    case ATLANTIS_NAME:
        return ATLANTIS_MAX
    case BERMUDA_NAME:
        return BERMUDA_MAX  
    default:
        return DEFAULT_MAX
    }   
}

//get the appropriate destination mod max from constants
func (game *Game) getDestinationModMax(name string) (int) {
    switch name{
    case ATLANTIS_NAME:
        return ATLANTIS_MOD_MAX
    case BERMUDA_NAME:
        return BERMUDA_MOD_MAX  
    default:
        return DEFAULT_MOD_MAX
    }   
}

//get the appropriate destination duration from constants
func (game *Game) getDestinationDuration(name string) (int) {
    switch name{
    case ATLANTIS_NAME:
        return ATLANTIS_MOD_MAX
    case BERMUDA_NAME:
        return BERMUDA_MOD_MAX  
    default:
        return DEFAULT_MOD_MAX
    }   
}