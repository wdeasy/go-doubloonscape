package game

import (
    "fmt"
    "math"

    "github.com/wdeasy/go-doubloonscape/storage"
)

//main function for handling all the steps of changing out the current captain
func (game *Game) newCaptain(GuildID string, UserID string) {
    err := game.changeRoles(GuildID, UserID)
    if err != nil {
        fmt.Printf("could not change roles: %s\n", err)
    }	

    game.changeCaptains(UserID)
    game.setMessage()   
}

//change out who is the current captain in the game
func (game *Game) changeCaptains(ID string) {
    game.removeCaptains()
    game.addCaptain(ID)
    game.currentCaptainID = ID
}

//set captain status to false for all captains
func (game *Game) removeCaptains(){
    for _, c := range game.captains {
        captain := c
        captain.Captain = false
        game.captains[captain.ID] = captain
    }	
}

//set captain status to true
func (game *Game) addCaptain(ID string) {
    captain := game.captains[ID]
    captain.Captain = true
    game.captains[ID] = captain	
}

//increment captains gold by prestige
func (game *Game) incrementCaptain() {
    if game.currentCaptainID == "" {
        return
    }

    captain := game.captains[game.currentCaptainID]
    captain.Gold = captain.Gold + int(game.goldModifier() * math.Floor(captain.Prestige))
    
    game.captains[game.currentCaptainID] = captain	
}

//create a captain and add it to the map
func (game *Game) createCaptain(ID string, Name string) {
    var captain storage.Captain

    captain.ID = ID 
    captain.Name = Name
    captain.Gold = DEFAULT_GOLD
    captain.Prestige = DEFAULT_PRESTIGE
    captain.Captain = false

    game.captains[ID] = captain
}