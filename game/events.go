package game

import (
	"time"

	"github.com/wdeasy/go-doubloonscape/storage"
)

func (game *Game) getEvents() ([]string) {
    return []string{CAPTAIN_NAME, INCREMENT_NAME, PICKPOCKET_NAME, PRESTIGE_NAME, TREASURE_NAME}
}

func (game *Game) getIdleEvents() ([]string) {
    return []string{PICKPOCKET_NAME, TREASURE_NAME}
}

//iterate over events
func (game *Game) checkEvents() {
    for _, event := range game.events {
        event.Ready()
    }

    for _, event := range game.getIdleEvents() {
        if game.events[event].Roll() {
            game.executeEvent(event, IDLE_EVENT, "")
        }
    }
}

//update the event info
func (game *Game) setEvent(name string) {
    var event storage.Event

    event.DB       = game.storage.DB
    event.Name     = name
    event.Last     = time.Now()
    event.Up       = false
    event.Chance   = game.getEventChance(name)
    event.Max      = game.getEventMax(name)
    event.Cooldown = game.getEventCooldown(name)

    game.events[event.Name] = &event	
}

//process the event reaction
func(game *Game) EventReactionReceived(Name string, UserID string, GuildID string) {
    if !game.events[Name].Up {
        return
    }
    
    game.executeEvent(Name, UserID, GuildID)   
}

//update events from constants
func (game *Game) updateEvents() {
    for _, event := range game.getEvents() {
        if _, ok := game.events[event]; !ok {
            game.setEvent(event)
            continue
        }
        
        game.events[event].Chance   = game.getEventChance(event)
        game.events[event].Max      = game.getEventMax(event)
        game.events[event].Cooldown = game.getEventCooldown(event)        
    }

    game.storage.SaveEvents(game.events)    
}

//execute the appropriate event
func(game *Game) executeEvent(Name string, UserID string, GuildID string) {
    game.events[Name].Reset()

    switch Name {
    case INCREMENT_NAME:
        game.incrementCaptain()
    case CAPTAIN_NAME:
        game.changeCaptainsInGameAndServer(GuildID, UserID) 
    case PRESTIGE_NAME:
        game.setPrestige(UserID)
    case TREASURE_NAME:           
        game.giveTreasure(UserID)                        
    case PICKPOCKET_NAME:
        game.executePickPocket(UserID) 
    default:
        return
    }
    
    game.setMessage()
}

//get the appropriate cooldown from constants
func (game *Game) getEventCooldown(name string) (int) {
    switch name{
    case PICKPOCKET_NAME:
        return PICKPOCKET_COOLDOWN
    case CAPTAIN_NAME:
        return CAPTAIN_COOLDOWN
    case INCREMENT_NAME:
        return INCREMENT_COOLDOWN
    case PRESTIGE_NAME:
        return PRESTIGE_COOLDOWN
    case TREASURE_NAME:
        return TREASURE_COOLDOWN        
    default:
        return DEFAULT_COOLDOWN
    }   
}

//get the appropriate chance from constants
func (game *Game) getEventChance(name string) (int) {
    switch name{
    case PICKPOCKET_NAME:
        return PICKPOCKET_CHANCE
    case CAPTAIN_NAME:
        return CAPTAIN_CHANCE
    case INCREMENT_NAME:
        return INCREMENT_CHANCE
    case PRESTIGE_NAME:
        return PRESTIGE_CHANCE
    case TREASURE_NAME:
        return TREASURE_CHANCE        
    default:
        return DEFAULT_CHANCE
    }   
}

//get the max number to use against the chance roll
func (game *Game) getEventMax(name string) (int) {
    switch name{
    case PICKPOCKET_NAME:
        return PICKPOCKET_MAX
    case CAPTAIN_NAME:
        return CAPTAIN_MAX
    case INCREMENT_NAME:
        return INCREMENT_MAX
    case PRESTIGE_NAME:
        return PRESTIGE_MAX
    case TREASURE_NAME:
        return TREASURE_MAX        
    default:
        return DEFAULT_MAX
    }   
}