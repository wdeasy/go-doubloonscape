package storage

import (
    "database/sql"
    "fmt"
    "time"
)

type Event struct {
    DB      *sql.DB    
    Name     string
    Last     time.Time
    Up       bool
}

//update event or insert if it doesnt exist
func (event *Event) Save() (error){
    updateStmt := `UPDATE events SET last = $1, up = $2 WHERE name = $3`
    _, err := event.DB.Exec(updateStmt, event.Last, event.Up, event.Name)
    if err != nil {
        return fmt.Errorf("could not update event %s: %w", event.Name, err)
    }

    insertStmt := `INSERT INTO events (name, last, up) SELECT CAST($1 AS VARCHAR), $2, $3 WHERE NOT EXISTS (SELECT 1 FROM events WHERE name = $1)`
    _, err = event.DB.Exec(insertStmt, event.Name, event.Last, event.Up)
    if err != nil {
        return fmt.Errorf("could not insert event %s: %w", event.Name, err)
    }

    return nil
}

//load events
func (storage *Storage) LoadEvents() (map[string]*Event, error){
    events := storage.CreateEventMap()

    rows, err := storage.DB.Query(`SELECT * FROM events ORDER BY name ASC`)
    if err != nil {
        return events, fmt.Errorf("could not select events: %w", err)
    }

    for rows.Next(){
        var event Event
        event.DB = storage.DB      
          
        err := rows.Scan(&event.Name, &event.Last, &event.Up)
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

//see if the event is off cooldown
func (event *Event) Ready(cooldown int) (bool){
    return time.Now().After(event.Last.Add(time.Minute * time.Duration(cooldown)))
}

//create an empty event map
func (storage *Storage) CreateEventMap() (map[string]*Event) {
    return make(map[string]*Event)
}