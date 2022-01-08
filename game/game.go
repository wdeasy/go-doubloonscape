package game

import (
    "fmt"
    "sort"
    "strings"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/wdeasy/go-doubloonscape/storage"
)

var (
    stats Stats
    treasure bool
)

type Game struct {
    storage *storage.Storage
    dg *discordgo.Session

    captains map[string]storage.Captain
    currentCaptainID string
}

type Stats struct {
    Leaderboard string
    Event string
}

//start the game
func InitGame(storage *storage.Storage) (*Game, error) {
    var game Game
    
    dg, err := game.InitBot()
    if err != nil {
        return nil, fmt.Errorf("could not initialize bot: %w", err)
    }

    game.dg = dg
    game.storage = storage

    game.captains, err = storage.LoadCaptains()
    if err != nil {
        return nil, fmt.Errorf("could not load captains: %w", err)
    }

    currentCaptain, err := storage.LoadCurrentCaptain()
    if err != nil {
        fmt.Printf("could not load current captain: %s\n", err)
    }
    game.currentCaptainID = currentCaptain.ID

    // Timer
    game.GameTimer()

    return &game, nil	
}

//save all info to storage
func (game *Game) SaveGame() {
    // start := time.Now()

    game.storage.SaveCaptains(game.captains)

    // end := time.Now()
    // diff := end.Sub(start)
    // fmt.Printf("Save Game took %f seconds.\n", diff.Seconds())
}

//main game loop
func (game *Game) GameTimer() {
    game.setStats()
    last := time.Now()

    i := 1
    ticker := time.NewTicker(60 * time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <- ticker.C:
                game.incrementCaptain()
                
                if (time.Now().Hour() != last.Hour()) {
                    game.setStats()
                    last = time.Now()
                }

                game.setMessage()	                  
                game.SaveGame()

                i++
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()
}

//leaderboard sorting
type Pair struct {
    Key string
    Value int
  }

//leaderboard sorting  
type PairList []Pair

//leaderboard sorting
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

//format leaderboard to string
func (game *Game) printLeaderboard(captains map[string]storage.Captain) (string) {

    pl := make(PairList, len(game.captains))
    i := 0
    for k, v := range game.captains {
        pl[i] = Pair{k, v.Gold}
        i++
    }

    sort.Sort(sort.Reverse(pl))

    var b strings.Builder
    for j, k := range pl {
        fmt.Fprintf(&b, "` %2d ` ` %-27s ` ` %7d `\n", j+1, firstN(captains[k.Key].Name,27), k.Value)
    }
    return b.String()
}

//update the stats struct
func (game *Game) setStats() (){
    captains := game.captains
    stats.Leaderboard = game.printLeaderboard(captains)

    stats.Event = ""
}

//truncate names
func firstN(s string, n int) string {
    if len(s) > n {
         return s[:n]
    }
    return s
}