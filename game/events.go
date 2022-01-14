package game

import (
    "math/rand"
    "time"

    "github.com/wdeasy/go-doubloonscape/storage"
)

//iterate over events
func (game *Game) checkEvents() {
    game.checkPickPocket()
}

//check if event is up
func (game *Game) checkEvent(name string, cooldown int, chance int) (bool){
    if _, ok := game.events[name]; !ok {
        game.setEvent(name)
        return false	
    }

    if game.events[name].Ready(game.getCooldown(name)) {
        return false
    }

    if rand.Intn(100) > chance {
        return false
    }

    game.setEvent(name)
    return true
}

//update the destination info
func (game *Game) setEvent(name string) {
    var event storage.Event

    event.DB = game.storage.DB
    event.Name = name
    event.Last = time.Now()
    event.Up = false

    game.events[event.Name] = &event	
}

//get the appropriate cooldown from constants
func (game *Game) getCooldown(name string) (int) {
    switch name{
    case PICKPOCKET_NAME:
        return PICKPOCKET_COOLDOWN
    default:
        return DEFAULT_COOLDOWN
    }   
}


