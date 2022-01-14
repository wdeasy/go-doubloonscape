package storage

import (
    "database/sql"
    "fmt"
    "time"
)

type Log struct {
    DB      *sql.DB    
    Text    string
    Time    time.Time
}

//insert log
func (log *Log) Save() (error){
    insertStmt := `INSERT INTO logs(text, time) values($1, $2)`
    _, err := log.DB.Exec(insertStmt, log.Text, log.Time)
    if err != nil {
        return fmt.Errorf("could not insert log %s: %w", log.Text, err)
    }

    return nil
}

//load logs
func (storage *Storage) LoadLogs(limit int) ([]*Log, error){
    var logs []*Log

    selectStmt := `SELECT * FROM logs ORDER BY time ASC LIMIT $1`
    rows, err := storage.DB.Query(selectStmt, limit)
    if err != nil {
        return logs, fmt.Errorf("could not select logs: %w", err)
    }

    for rows.Next(){
        var log Log
        log.DB = storage.DB      
          
        err := rows.Scan(&log.Text, &log.Time)
        if err != nil {
            fmt.Printf("could not load log: %s\n", err)
        }
        
        logs = append(logs, &log)
    }

    return logs, nil	
}

//save logs
func (storage *Storage) SaveLogs(logs []*Log) {

    storage.ClearLogs()

    for _, s := range logs {
        err := s.Save()

        if err != nil {
            fmt.Printf("could not save log %s: %s\n", s.Text, err)
        }		
    }	
}

//truncate logs
func (storage *Storage) ClearLogs() (error){
    truncateStmt := `TRUNCATE logs`
    _, err := storage.DB.Exec(truncateStmt)
    if err != nil {
        return fmt.Errorf("could not truncate logs: %w", err)
    }

    return nil
}