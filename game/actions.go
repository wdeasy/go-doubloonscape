package game

import (
    "fmt"
    "math/rand"
)

//convert gold to prestige
func (game *Game) addPrestige(UserID string) {
    captain := game.captains[UserID]
    captain.Prestige = captain.Prestige + (float32(captain.Gold) * PRESTIGE_CONVERSION)
    captain.Gold = 0
    game.captains[UserID] = captain
}

//give a random amount of gold to a captain
func (game *Game) giveTreasure(UserID string) {
    treasure := rand.Intn(TREASURE_MAX)

    captain := game.captains[UserID]
    captain.Gold = captain.Gold + treasure
    game.captains[UserID] = captain

    stats.Event = fmt.Sprintf("%s looted Treasure worth %d gold!", captain.Name, treasure)
}