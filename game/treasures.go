package game

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

//give a random amount of gold to a captain
func (game *Game) giveTreasure(UserID string) {
    if UserID == IDLE_EVENT {
        UserID = game.currentCaptainID
    }

    captain := game.captains[UserID]
    captain.Gold = captain.Gold + int64(game.events[TREASURE_NAME].Amount)
    game.captains[UserID] = captain

    game.addToLogs(game.treasureString(UserID, game.events[TREASURE_NAME].Amount))
    game.setMessage()	

    game.events[TREASURE_NAME].Amount = TREASURE_START
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
    game.addToLogs(fmt.Sprintf("The treasure chest has grown to %d", game.events[TREASURE_NAME].Amount))
}

//increment treasure by 1 each minute
func (game *Game) incrementTreasure() {
    game.events[TREASURE_NAME].Amount += int64(game.goldModifier() * math.Floor(game.getPrestige()))
}

//find highest prestige
func (game *Game) getPrestige() (float64){
    prestige := 1.0

    for _, c := range game.captains {
        if c.Prestige > prestige {
            prestige = c.Prestige
        }
    }

    return prestige
}
