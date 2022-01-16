package game

import (
    "fmt"
    "math/rand"
    "strconv"
    "strings"
)

//give a random amount of gold to a captain
func (game *Game) giveTreasure(UserID string) {
    if !game.treasure.Up {
        return 
    }

    game.treasure.Up = false    
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

    game.treasure.Increment(game.goldModifier())
}

//generate a treasure string for the logs
func (game *Game) treasureString(UserID string, amount int) (string) {
    log := "%s 𝔩𝔬𝔬𝔱𝔢𝔡 𝔱𝔥𝔢 𝔗𝔯𝔢𝔞𝔰𝔲𝔯𝔢 𝔣𝔬𝔯 %d"    

    logLength := getStringLength(log)
    userLength := getStringLength(game.captains[UserID].Name)
    
    amountLength := len(strconv.Itoa(amount))
    variableLength := (strings.Count(log, "%")*2)

    i := LOG_LINE_LENGTH - amountLength - logLength + variableLength + 4

    if userLength > i {
        userLength = i
    }

    return fmt.Sprintf(log, firstN(game.captains[UserID].Name,userLength), amount)
}

//add the treasure amount to the logs
func (game *Game) logTreasure() {
    game.addToLogs(fmt.Sprintf("𝔗𝔥𝔢 𝔱𝔯𝔢𝔞𝔰𝔲𝔯𝔢 𝔠𝔥𝔢𝔰𝔱 𝔥𝔞𝔰 𝔤𝔯𝔬𝔴𝔫 𝔱𝔬 %d", game.treasure.Amount))
}