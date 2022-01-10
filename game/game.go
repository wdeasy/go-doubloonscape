package game

import (
    "fmt"
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

    destinations map[string]storage.Destination
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

    err = game.LoadGame()
    if err != nil {
        return nil, fmt.Errorf("could not load game: %w", err)
    }

    // Timer
    game.GameTimer()

    return &game, nil	
}

//load previous game information from storage
func (game *Game) LoadGame() (error){
    var err error
    game.captains, err = game.storage.LoadCaptains()
    if err != nil {
        return fmt.Errorf("could not load captains: %w", err)
    }

    game.destinations, err = game.storage.LoadDestinations()
    if err != nil {
        return fmt.Errorf("could not load destinations: %w", err)
    }

    currentCaptain, err := game.storage.LoadCurrentCaptain()
    if err != nil {
        fmt.Printf("could not load current captain: %s\n", err)
    }
    game.currentCaptainID = currentCaptain.ID

    return nil
}

//save all info to storage
func (game *Game) SaveGame() {
    // start := time.Now()

    game.storage.SaveCaptains(game.captains)
    game.storage.SaveDestinations(game.destinations)

    // end := time.Now()
    // diff := end.Sub(start)
    // fmt.Printf("Save Game took %f seconds.\n", diff.Seconds())
}

//main game loop
func (game *Game) GameTimer() {
    game.visitDestinations()
    game.setStats()
    last := time.Now()

    i := 1
    ticker := time.NewTicker(time.Second)
    quit := make(chan struct{})
    go func() {
        for {
            select {
            case <- ticker.C:
                if (i % game.timeModifier() == 0) {
                    game.visitDestinations()
                    game.incrementCaptain()
                    game.setMessage()	                  
                    game.SaveGame()                    
                }
                
                if (time.Now().Hour() != last.Hour()) {
                    game.setStats()
                    last = time.Now()
                }

                i++
            case <- quit:
                ticker.Stop()
                return
            }
        }
    }()
}

//time modifier
func (game *Game) timeModifier() (int) {
    if game.destinations["bermuda"].Amount == 0 {
        return 60
    }
        
    return int(60 * (1 + (0.01 * float64(game.destinations["bermuda"].Amount))))
}

//gold modifier
func (game *Game) goldModifier() (float64) {
    if game.destinations["atlantis"].Amount == 0 {
        return 1
    }

    return float64(game.destinations["atlantis"].Amount)
}

//update the stats struct
func (game *Game) setStats() (){
    stats.Leaderboard = game.printLeaderboard()

    stats.Event = ""
}

