package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/wdeasy/go-doubloonscape/storage"
    "github.com/wdeasy/go-doubloonscape/game"
)

func main() {

    //initialize the storage
    storage, err := storage.InitStorage()
    if err != nil {
        fmt.Println(err)
        return
    }

    //initialize the game
    game, err := game.InitGame(storage)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("DoubloonScape is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    //stop the game
    game.SaveGame()
    game.CloseBot()
    storage.CloseStorage()
    fmt.Println("DoubloonScape is now stopped.")
}
