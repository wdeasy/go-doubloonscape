package game

import (
	"fmt"
	"strconv"
	"strings"
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
}

//create a pickpocketstring for the logs
func (game *Game) pickPocketString(pickpocketeer string, amount int64) (string) {
    log := "%s 𝔭𝔦𝔠𝔨𝔭𝔬𝔠𝔨𝔢𝔱𝔰 %s 𝔣𝔬𝔯 %d"

    logLength := getStringLength(log)
    captainLength := getStringLength(game.captains[game.currentCaptainID].Name)
    pickpocketeerLength := getStringLength(game.captains[pickpocketeer].Name)
    
    maxLength := 42
    amountLength := len(strconv.FormatInt(amount, 10))
    variableLength := (strings.Count(log, "%")*2)

    i := maxLength - amountLength - logLength + variableLength + 6

    if i < 4 {
        i = 4
    }

    captainN := i/2
    pickpocketeerN := i/2

    if captainLength < i/2 {
        captainN = captainLength
        pickpocketeerN = i - captainLength
    }

    if pickpocketeerLength < i/2 {
        pickpocketeerN = pickpocketeerLength
        captainN = i - pickpocketeerLength
    }

    return fmt.Sprintf(log,
            firstN(game.captains[pickpocketeer].Name,pickpocketeerN), 
            firstN(game.captains[game.currentCaptainID].Name,captainN), amount)
}

func getStringLength(String string) (int) {
    var StringLength int
    for range String {
        StringLength++
    }
    return StringLength
}