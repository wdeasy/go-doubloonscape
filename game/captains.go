package game

import (
    "fmt"
    "math/rand"

    "github.com/wdeasy/go-doubloonscape/storage"
)

//main function for handling all the steps of changing out the current captain
func (game *Game) changeCaptainsInGameAndServer(GuildID string, UserID string) {
    if game.currentCaptainID == UserID {
        return
    }

    err := game.changeRoles(GuildID, UserID)
    if err != nil {
        printLog(fmt.Sprintf("could not change roles: %s\n", err))
    }	

    game.changeCaptainsInGame(UserID)
    printLog(game.captainString())
    game.setMessage()   
}

//change out who is the current captain in the game
func (game *Game) changeCaptainsInGame(ID string) {
    game.removeCaptains()
    game.addCaptain(ID)
    game.currentCaptainID = ID
}

//set captain status to false for all captains
func (game *Game) removeCaptains(){
    for _, c := range game.captains {
        c.DemoteCaptain()
    }	
}

//set captain status to true
func (game *Game) addCaptain(ID string) {
    game.captains[ID].PromoteCaptain()
}

//increment captains gold by prestige
func (game *Game) incrementCaptain() {
    if game.currentCaptainID == "" {
        return
    }

    game.captains[game.currentCaptainID].IncrementDoubloon(game.goldModifier())
}

//create a captain and add it to the map
func (game *Game) createCaptain(ID string, Name string) {
    var captain storage.Captain

    captain.DB = game.storage.DB
    captain.ID = ID 
    captain.Name = Name
    captain.Gold = DEFAULT_GOLD
    captain.Prestige = DEFAULT_PRESTIGE
    captain.Captain = false

    game.captains[ID] = &captain
}

//create a captain in the map if it does not exist
func (game *Game) checkIfCaptainExists(UserID string, GuildID string) (error) {
    if _, ok := game.captains[UserID]; !ok {
        err := game.addCaptainFromDiscordReaction(GuildID, UserID)

        if err != nil {
            return fmt.Errorf("could not add captain %s from discord reaction: %w", UserID, err)
        }		
    }
    
    return nil
}

//get random captain
func (game *Game) randomCaptainID() (string, error) {
    var UserID string
    var captains []string

    for _, c := range game.captains {
        if c.ID != game.currentCaptainID {
            captains = append(captains, c.ID)
        }
    }

    if len(captains) == 0 {
        return UserID, fmt.Errorf("not enough captains")
    }

    i := rand.Intn(len(captains))
    for _, v := range captains {
        if i == 0 {
            UserID = v
        }
        i--
    }

    return UserID, nil
}

//generate a captain string for the logs
func (game *Game) captainString() (string) {
    return fmt.Sprintf("%s 𝔦𝔰 𝔱𝔥𝔢 𝔠𝔞𝔭𝔱𝔞𝔦𝔫 𝔫𝔬𝔴", firstN(game.captains[game.currentCaptainID].Name,10))
}

//create a new captain with info from the discord reaction
func (game *Game) addCaptainFromDiscordReaction(GuildID string, UserID string) (error) {
    m, err := game.dg.GuildMember(GuildID, UserID)
    if err != nil {
        return fmt.Errorf("could not get guild member info %s: %w", UserID, err)
    }
        
    name := getName(m.Nick, m.User.Username)
    
    if m.User.Bot {
        return fmt.Errorf("%s is a bot", name)
    }

    game.createCaptain(UserID, name)	

    return nil
}

//create a new captain with info from the discord message
func (game *Game) addCaptainFromDiscordMessage(UserID string, Nick string, UserName string) {
    name := getName(Nick, UserName)
    game.createCaptain(UserID, name)	
}

//used to update the captains prestige
func (game *Game) setPrestige(UserID string) {
    if UserID != game.currentCaptainID {
        return
    }

    game.captains[UserID].AddPrestige(PRESTIGE_CONVERSION)
    game.setMessage()
    
    if game.captains[UserID].Prestige > game.treasure.Prestige {
        game.treasure.SetPrestige(game.captains[UserID].Prestige)
    }
}