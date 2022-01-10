package storage

import (
	"fmt"
	"time"
)

type Destination struct {
    Name string
    End time.Time
	Amount int
}

//update destination or insert if it doesnt exist
func (storage *Storage) SaveDestination(destination Destination) (error){
    updateStmt := `UPDATE destinations SET end_time = $1, amount = $2 WHERE name = $3`
    _, err := storage.DB.Exec(updateStmt, destination.End, destination.Amount, destination.Name)
    if err != nil {
        return fmt.Errorf("could not update destination %s: %w", destination.Name, err)
    }

    insertStmt := `INSERT INTO destinations (name, end_time, amount) SELECT CAST($3 AS VARCHAR), $1, $2 WHERE NOT EXISTS (SELECT 1 FROM destinations WHERE name = $3)`
    _, err = storage.DB.Exec(insertStmt, destination.End, destination.Amount, destination.Name)
    if err != nil {
        return fmt.Errorf("could not insert destination %s: %w", destination.Name, err)
    }

    return nil
}

//load destinations
func (storage *Storage) LoadDestinations() (map[string]Destination, error){
    destinations := make(map[string]Destination)

    rows, err := storage.DB.Query(`SELECT * FROM destinations ORDER BY name ASC`)
    if err != nil {
        return destinations, fmt.Errorf("could not select destinations: %w", err)
    }

    for rows.Next(){
        var destination Destination
        err := rows.Scan(&destination.Name, &destination.End, &destination.Amount)

        if err != nil {
            fmt.Printf("could not load destination: %s\n", err)
        }
        
        destinations[destination.Name] = destination	
    }

    return destinations, nil	
}

//save destinations
func (storage *Storage) SaveDestinations(destinations map[string]Destination) {
    for _, s := range destinations {
        err := storage.SaveDestination(s)

        if err != nil {
            fmt.Printf("could not save destination %s: %s\n", s.Name, err)
        }		
    }	
}