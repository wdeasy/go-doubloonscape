package game

import (
    "fmt"
)

//check if pickpocket is available
func (game *Game) checkPickPocket() {
    if !game.checkEvent(PICKPOCKET_NAME, PICKPOCKET_COOLDOWN, PICKPOCKET_CHANCE) {
        return
    }

    if len(game.captains) < 2 {
        return
    }

    PickPocketeer, err := game.randomCaptainID()
    if err != nil {
        printLog(fmt.Sprintf("could not execute pickpocket: %s\n", err))
        return
    }

    game.executePickPocket(PickPocketeer)
    
}

//execute the pickpocket
func (game *Game) executePickPocket(pickpocketeer string) {
    max := float64(game.captains[game.currentCaptainID].Gold) * (float64(PICKPOCKET_MAX) * .01)
    amount := RandInt64(1, int64(max))

    if amount == 0 {
        return
    }

    game.captains[game.currentCaptainID].TakeDoubloons(amount)
    game.captains[pickpocketeer].GiveDoubloons(amount)

    game.addToLogs(game.pickPocketString(pickpocketeer, amount))
    game.setMessage()	

    game.removeCurrentEvent(PICKPOCKET_NAME)    
}

//create a pickpocketstring for the logs
func (game *Game) pickPocketString(pickpocketeer string, amount int64) (string) {
    return fmt.Sprintf("%s 𝔭𝔦𝔠𝔨𝔭𝔬𝔠𝔨𝔢𝔱𝔰 %s 𝔣𝔬𝔯 %d", 
            firstN(game.captains[pickpocketeer].Name,10), 
            firstN(game.captains[game.currentCaptainID].Name,10), amount)
}