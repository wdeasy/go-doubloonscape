package storage

import (
    "database/sql"
    "fmt"
    "math"
)

type Captain struct {
    DB      *sql.DB       
    ID       string
    Name     string
    Gold     int64
    Captain  bool
    Prestige float64
}

//load captains
func (storage *Storage) LoadCaptains() (map[string]*Captain, error){
    captains := make(map[string]*Captain)

    rows, err := storage.DB.Query(`SELECT * FROM captains ORDER BY gold DESC`)
    if err != nil {
        return captains, fmt.Errorf("could not select captains: %w", err)
    }

    for rows.Next(){
        var captain Captain 
        captain.DB = storage.DB
        err := rows.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain, &captain.Prestige)

        if err != nil {
            fmt.Printf("could not load captain: %s\n", err)
        }
        
        captains[captain.ID] = &captain	
    }

    return captains, nil	
}

//load current captain
func (storage *Storage) LoadCurrentCaptain() (*Captain, error) {
    selectStmt := `SELECT id, name, gold, captain, prestige FROM captains WHERE captain = true`
    row := storage.DB.QueryRow(selectStmt)

    var captain Captain 
    err := row.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain, &captain.Prestige)

    switch err {
    case sql.ErrNoRows:
        return &captain, fmt.Errorf("there is no current captain in the database: %w", err)
    case nil:
        return &captain, nil
    default:
        return &captain, err
    }
}

//save captains
func (storage *Storage) SaveCaptains(captains map[string]*Captain) {
    for _, s := range captains {
        err := s.Save()

        if err != nil {
            fmt.Printf("could not save captain %s: %s\n", s.Name, err)
        }		
    }	
}

//update captain or insert if it doesnt exist
func (captain *Captain) Save() (error){
    updateStmt := `UPDATE captains SET name = $1, gold = $2, captain = $3, prestige = $4 WHERE id = $5`
    _, err := captain.DB.Exec(updateStmt, captain.Name, captain.Gold, captain.Captain, captain.Prestige, captain.ID)
    if err != nil {
        return fmt.Errorf("could not update captain %s: %w", captain.Name, err)
    }

    insertStmt := `INSERT INTO captains (id, name, gold, captain, prestige) SELECT $5, $1, $2, $3, $4 WHERE NOT EXISTS (SELECT 1 FROM captains WHERE id=$5)`
    _, err = captain.DB.Exec(insertStmt, captain.Name, captain.Gold, captain.Captain, captain.Prestige, captain.ID)
    if err != nil {
        return fmt.Errorf("could not insert captain %s: %w", captain.Name, err)
    }

    return nil
}

//remove the current captain
func (captain *Captain) DemoteCaptain() {
    captain.Captain = false
}

//make the current captain
func (captain *Captain) PromoteCaptain() {
    captain.Captain = true
}

//convert captain gold to prestige
func (captain *Captain) AddPrestige(PrestigeConversion float64) {
    captain.Prestige = captain.Prestige + (float64(captain.Gold) * PrestigeConversion)
    captain.Gold = 0
}

//increment captain doubloons by turn or reaction
func (captain *Captain) IncrementDoubloons(goldModifier float64) {
    captain.Gold = captain.Gold + int64(goldModifier * math.Floor(captain.Prestige))
}

//award x amount of doubloons
func (captain *Captain) GiveDoubloons(amount int64) {
    captain.Gold += amount
}

//remove x amount of doubloons
func (captain *Captain) TakeDoubloons(amount int64) {
    captain.Gold -= amount
}