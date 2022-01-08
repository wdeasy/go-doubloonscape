package storage

import (
	"database/sql"
	"fmt"
)

type Captain struct {
    ID string
    Name string
    Gold int
    Captain bool
    Prestige float32
}

//update captain or insert if it doesnt exist
func (storage *Storage) SaveCaptain(captain Captain) (error){
    updateStmt := `UPDATE captains SET name = $1, gold = $2, captain = $3, prestige = $4 WHERE id = $5`
    _, err := storage.DB.Exec(updateStmt, captain.Name, captain.Gold, captain.Captain, captain.Prestige, captain.ID)
	if err != nil {
		return fmt.Errorf("could not update captain %s: %w", captain.Name, err)
	}

    insertStmt := `INSERT INTO captains (id, name, gold, captain, prestige) SELECT $5, $1, $2, $3, $4 WHERE NOT EXISTS (SELECT 1 FROM captains WHERE id=$5)`
    _, err = storage.DB.Exec(insertStmt, captain.Name, captain.Gold, captain.Captain, captain.Prestige, captain.ID)
	if err != nil {
		return fmt.Errorf("could not insert captain %s: %w", captain.Name, err)
	}

	return nil
}

//load captains
func (storage *Storage) LoadCaptains() (map[string]Captain, error){
	captains := make(map[string]Captain)

    rows, err := storage.DB.Query(`SELECT * FROM captains ORDER BY gold DESC`)
	if err != nil {
		return captains, fmt.Errorf("could not select captains: %w", err)
	}

    for rows.Next(){
        var captain Captain 
        err := rows.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain, &captain.Prestige)

        if err != nil {
            fmt.Printf("could not load captain: %s\n", err)
        }
        
		captains[captain.ID] = captain	
    }

    return captains, nil	
}

//save captains
func (storage *Storage) SaveCaptains(captains map[string]Captain) {
    for _, s := range captains {
        err := storage.SaveCaptain(s)

		if err != nil {
			fmt.Printf("could not save captain %s: %s\n", s.Name, err)
		}		
    }	
}

//load current captain
func (storage *Storage) LoadCurrentCaptain() (Captain, error) {
	selectStmt := `SELECT id, name, gold, captain, prestige FROM captains WHERE captain = true`
	row := storage.DB.QueryRow(selectStmt)

	var captain Captain 
	err := row.Scan(&captain.ID, &captain.Name, &captain.Gold, &captain.Captain, &captain.Prestige)

	switch err {
	case sql.ErrNoRows:
		return captain, fmt.Errorf("there is no current captain in the database: %w", err)
	case nil:
		return captain, nil
	default:
		return captain, err
    }
}