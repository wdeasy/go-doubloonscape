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

    if !game.events[name].Ready(game.getCooldown(name)) {
        return false
    }

    if rand.Intn(100) > chance {
        return false
    }

    game.resetEvent(name)
    return true
}

//update the event info
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

//reset the event
func (game *Game) resetEvent(name string) {
    game.events[name].Last = time.Now()
    game.events[name].Up = false
}

//add event to current events
func (game *Game) addCurrentEvent(name string) {
    game.currentEvents[name] = game.events[name]
}

//remove event from current events
func (game *Game) removeCurrentEvent(name string) {
    delete(game.currentEvents, name)
}


func(game *Game) EventReactionReceived(Name string, UserID string) {
    if !game.events[Name].Up {
        return
    }
    
    game.resetEvent(Name)
    
    if UserID == game.currentCaptainID {
        return
    }
    
     game.executeEvent(Name, UserID)   
}

func(game *Game) executeEvent(Name string, UserID string) {
    switch Name {
    case PICKPOCKET_NAME:
        game.executePickPocket(UserID) 
    default:
        return
    }
    
    game.setMessage()	
    game.removeCurrentEvent(Name)      
}

func(game *Game) loadCurrentEvents(){
    for _, e := range game.events {
        if e.Up {
            game.addCurrentEvent(e.Name)
        }
    }     
}