package storage

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type Event struct {
    DB      *sql.DB    
    Name     string
    Last     time.Time
    Up       bool
    Amount   int64
    Chance   int
    Max      int
    Cooldown int
}

//update event or insert if it doesnt exist
func (event *Event) Save() (error){
    updateStmt := `UPDATE events SET last = $1, up = $2, amount = $3, chance = $4, max = $5, cooldown = $6 WHERE name = $7`
    _, err := event.DB.Exec(updateStmt, event.Last, event.Up, event.Amount, event.Chance, event.Max, event.Cooldown, event.Name)
    if err != nil {
        return fmt.Errorf("could not update event %s: %w", event.Name, err)
    }

    insertStmt := `INSERT INTO events (name, last, up, amount, chance, max, cooldown) SELECT CAST($1 AS VARCHAR), $2, $3, $4, $5, $6, $7 WHERE NOT EXISTS (SELECT 1 FROM events WHERE name = $1)`
    _, err = event.DB.Exec(insertStmt, event.Name, event.Last, event.Up, event.Amount, event.Chance, event.Max, event.Cooldown)
    if err != nil {
        return fmt.Errorf("could not insert event %s: %w", event.Name, err)
    }

    return nil
}

//load events
func (storage *Storage) LoadEvents() (map[string]*Event, error){
    events := make(map[string]*Event)

    rows, err := storage.DB.Query(`SELECT * FROM events ORDER BY name ASC`)
    if err != nil {
        return events, fmt.Errorf("could not select events: %w", err)
    }

    for rows.Next(){
        var event Event
        event.DB = storage.DB      
          
        err := rows.Scan(&event.Name, &event.Last, &event.Up, &event.Amount, &event.Chance, &event.Max, &event.Cooldown)
        if err != nil {
            fmt.Printf("could not load event: %s\n", err)
        }
        
        events[event.Name] = &event	
    }

    return events, nil	
}

//save events
func (storage *Storage) SaveEvents(events map[string]*Event) {
    for _, s := range events {
        err := s.Save()

        if err != nil {
            fmt.Printf("could not save event %s: %s\n", s.Name, err)
        }		
    }	
}

func (event *Event) Ready() (bool) {
    //if event already up
    if event.Up {
        return true
    }

    // if cooldown not up
    if !time.Now().After(event.Last.Add(time.Minute * time.Duration(event.Cooldown))) {
        return false
    } 

    // if roll fails
    if !event.Roll() {
        return false
    }

    event.Up = true  

    return true
}

func (event *Event) Roll() (bool) {
    return rand.Intn(event.Max) <= event.Chance
}

//reset the event
func (event *Event) Reset() {
    event.Last = time.Now()
    event.Up = false
}
