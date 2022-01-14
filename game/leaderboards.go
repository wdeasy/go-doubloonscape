package game

import (
    "fmt"
    "sort"
    "strings"
)

//leaderboard sorting
type Pair struct {
    Key string
    Value int64
  }

//leaderboard sorting  
type PairList []Pair

//leaderboard sorting
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

//format leaderboard to string
func (game *Game) printLeaderboard() (*string) {

    pl := make(PairList, len(game.captains))
    i := 0
    for k, v := range game.captains {
        pl[i] = Pair{k, v.Gold}
        i++
    }

    sort.Sort(sort.Reverse(pl))

    var b strings.Builder
    fmt.Fprintf(&b, "**𝔏𝔢𝔞𝔡𝔢𝔯𝔅𝔬𝔞𝔯𝔡**\n") 
    for j, k := range pl {
        fmt.Fprintf(&b, "` %2d ` ` %-27s ` ` %7d `\n", j+1, firstN(game.captains[k.Key].Name,27), k.Value)
    }

    String := b.String()
    return &String
}

//update the embed leaderboard info
func (game *Game) setLeaderboard() {
    game.stats.Leaderboard = game.printLeaderboard()
}