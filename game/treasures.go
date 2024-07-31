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
    // if treasureChance() {
    //     game.giveTreasure(game.currentCaptainID)
    // }

    game.treasure.Increment(game.goldModifier())
}

//generate a treasure string for the logs
func (game *Game) treasureString(UserID string, amount int64) (string) {
    log := "%s looted the Treasure for %d"    

    logLength := getStringLength(log)
    userLength := getStringLength(game.captains[UserID].Name)
    
    amountLength := len(strconv.FormatInt(amount, 10))
    variableLength := (strings.Count(log, "%")*2)

    i := LOG_LINE_LENGTH - amountLength - logLength + variableLength

    if userLength > i {
        userLength = i
    }

    return fmt.Sprintf(log, firstN(game.captains[UserID].Name,userLength), amount)
}

//add the treasure amount to the logs
func (game *Game) logTreasure() {
    game.addToLogs(fmt.Sprintf("The treasure chest has grown to %d", game.treasure.Amount))
}