package game

import (
    "fmt"
    "math/rand"
)

//give a random amount of gold to a captain
func (game *Game) giveTreasure(UserID string) {
    captain := game.captains[UserID]
    captain.Gold = captain.Gold + int64(game.treasure.Amount)
    game.captains[UserID] = captain

    game.addToLogs(game.treasureString(UserID, game.treasure.Amount))
    game.setMessage()	

    game.treasure.Reset()
}

//roll for the treasure
func treasureChance() (bool) {
    return rand.Intn(1000) <= TREASURE_CHANCE
}

//do the treasure turn
func (game *Game) checkTreasure() {
    if treasureChance() {
        game.giveTreasure(game.currentCaptainID)
    }

    game.treasure.Increment()
}

//generate a treasure string for the logs
func (game *Game) treasureString(UserID string, amount int) (string) {
    return fmt.Sprintf("%s 𝔩𝔬𝔬𝔱𝔢𝔡 𝔱𝔥𝔢 𝔗𝔯𝔢𝔞𝔰𝔲𝔯𝔢 𝔣𝔬𝔯 %d", 
            firstN(game.captains[UserID].Name,10), amount)
}

//add the treasure amount to the logs
func (game *Game) logTreasure() {
    game.addToLogs(fmt.Sprintf("𝔗𝔥𝔢 𝔱𝔯𝔢𝔞𝔰𝔲𝔯𝔢 𝔠𝔥𝔢𝔰𝔱 𝔥𝔞𝔰 𝔤𝔯𝔬𝔴𝔫 𝔱𝔬 %d", game.treasure.Amount))
}