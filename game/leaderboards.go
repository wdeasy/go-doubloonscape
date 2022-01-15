package game

import (
	"fmt"
	"sort"
	"strconv"
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

    n := len(strconv.FormatInt(pl[0].Value,10))
    if n >= 32 {
        n = 32
    }
    if n < 7 {
        n = 7
    }
    o := 34 - n
        
    for j, k := range pl {
        fmt.Fprintf(&b, "` %2d ` ` %-s ` ` %*d `\n", j+1, firstN(game.captains[k.Key].Name,o), n, k.Value)
    }

    String := b.String()
    return &String
}

//update the embed leaderboard info
func (game *Game) setLeaderboard() {
    game.stats.Leaderboard = game.printLeaderboard()
}