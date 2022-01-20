package storage

import (
    "database/sql"
    "fmt"
    "math"
)

type Treasure struct {
    DB      *sql.DB
    Amount   int64	
    Up       bool
    Prestige float64    
}

//load treasure
func (storage *Storage) LoadTreasure() (*Treasure, error){
    var treasure Treasure
    treasure.DB = storage.DB
    treasure.Up = false
    treasure.Amount = 1
    treasure.Prestige = 1


    rows, err := storage.DB.Query(`SELECT * FROM treasure LIMIT 1`)
    if err != nil {
        return &treasure, fmt.Errorf("could not select treasure: %w", err)
    }

    for rows.Next(){
        err := rows.Scan(&treasure.Amount, &treasure.Up, &treasure.Prestige)
        if err != nil {
            fmt.Printf("could not load treasure: %s\n", err)
        }
    }

    return &treasure, nil	
}

//update treasure or insert if it doesnt exist
func (treasure *Treasure) Save() (error){
    updateStmt := `UPDATE treasure SET amount = $1, up = $2, prestige = $3`
    _, err := treasure.DB.Exec(updateStmt, treasure.Amount, treasure.Up, treasure.Prestige)
    if err != nil {
        return fmt.Errorf("could not update treasure: %w", err)
    }

    insertStmt := `INSERT INTO treasure (amount, up, prestige) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM treasure)`
    _, err = treasure.DB.Exec(insertStmt, treasure.Amount, treasure.Up, treasure.Prestige)
    if err != nil {
        return fmt.Errorf("could not insert treasure: %w", err)
    }

    return nil
}

//increment treasure by 1 each minute
func (treasure *Treasure) Increment(goldModifier float64) {
    treasure.Amount += int64(goldModifier * math.Floor(treasure.Prestige))
}

//reset treasure after it is claimed
func (treasure *Treasure) Reset() {
    treasure.Up = false
    treasure.Amount = 1
}

//convert captain gold to prestige
func (treasure *Treasure) SetPrestige(newPrestige float64) {
    treasure.Prestige = newPrestige
}