package game

import (
    "fmt"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/wdeasy/go-doubloonscape/storage"
)

type Game struct {
    storage *storage.Storage
    dg *discordgo.Session

    captains map[string]*storage.Captain
    currentCaptainID string
    currentMessageID string

    destinations map[string]*storage.Destination
    treasure *storage.Treasure
    events map[string]*storage.Event
    logs []*storage.Log

    stats Stats
}

type Stats struct {
    Leaderboard *string
    Destinations *string
    Log *string
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
        printLog(fmt.Sprintf("could not load current captain: %s\n", err))
    }
    game.currentCaptainID = currentCaptain.ID

    game.treasure, err = game.storage.LoadTreasure()
    if err != nil {
        return fmt.Errorf("could not load treasure: %w", err)
    }    

    game.events, err = game.storage.LoadEvents()
    if err != nil {
        return fmt.Errorf("could not load events: %w", err)
    }

    game.logs, err = game.storage.LoadLogs(MAX_LOG_LENGTH)
    if err != nil {
        return fmt.Errorf("could not load logs: %w", err)
    }    

    return nil
}

//save all info to storage
func (game *Game) SaveGame() {
    // start := time.Now()

    game.storage.SaveCaptains(game.captains)
    game.storage.SaveDestinations(game.destinations)
    game.treasure.Save()
    game.storage.SaveEvents(game.events)
    game.storage.SaveLogs(game.logs)

    // end := time.Now()
    // diff := end.Sub(start)
    // printLog(fmt.Sprintf("Save Game took %f seconds.\n", diff.Seconds()))
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
                    game.checkTreasure()
                    game.checkEvents()
                    game.setMessage()	                  
                    game.SaveGame()                    
                }
                
                if (time.Now().Day() != last.Day()) {
                    game.logTreasure()
                }

                if (time.Now().Hour() != last.Hour()) {
                    game.setLeaderboard()
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
    if game.destinations[BERMUDA_NAME].Amount == 0 {
        return 60
    }
        
    return int(60 * (1 + (0.01 * float64(game.destinations[BERMUDA_NAME].Amount))))
}

//gold modifier
func (game *Game) goldModifier() (float64) {
    if game.destinations[ATLANTIS_NAME].Amount == 0 {
        return 1
    }

    return float64(game.destinations[ATLANTIS_NAME].Amount)
}

//update the stats struct
func (game *Game) setStats() {
    game.setLeaderboard()
    game.setDestinations()    
    game.setLogs()
}