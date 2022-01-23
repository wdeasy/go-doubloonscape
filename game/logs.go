package game

import (
    "fmt"
    "strings"
    "time"

    "github.com/wdeasy/go-doubloonscape/storage"
)

//add to the logs and remove the oldest log
func (game *Game) addToLogs(Text string) {
    log := game.newLog(Text)

    if len(game.logs) == MAX_LOG_LINES {
        _, game.logs = game.logs[0], game.logs[1:]
    }
    game.logs = append(game.logs, &log)
    game.setLogs()
    
    printLog(log.Text)
}

//create a new log
func (game *Game) newLog(text string) (storage.Log) {
    return storage.Log{DB: game.storage.DB, Text: text, Time: time.Now()}
}

func (game *Game) setLogs() {
    game.stats.Log = game.logsString()
}

//generate logs string for embed
func (game *Game) logsString() (*string) {
    var b strings.Builder

    fmt.Fprintf(&b, "%s\n", "```fix")
    for _, v := range game.logs {
        fmt.Fprintf(&b, " %s\n", v.Text)
    }

    if len(game.logs) == 0 {
        fmt.Fprintf(&b, " ℜ𝔢𝔡 𝔖𝔨𝔦𝔢𝔰 𝔞𝔱 𝔫𝔦𝔤𝔥𝔱")
    }    
    fmt.Fprintf(&b, "%s\n", "```")

    String := b.String()
    return &String
}

//print log information to the console
func printLog(log string) {
    fmt.Printf("[%s] %s\n", time.Now().Format("01/02/06 15:04:05"), log)
}