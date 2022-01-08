package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)
var (
    DatabaseURL string = os.Getenv("DATABASE_URL")
)

type Storage struct {
  DB *sql.DB
}

//open and test the database connection
func InitStorage() (*Storage, error) {
    if DatabaseURL == "" {
      return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
    }

    db, err := sql.Open("postgres", DatabaseURL)
    if err != nil {
         return nil, fmt.Errorf("could not open sql: %w", err)
    }

    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("could not ping DB: %w", err)
    }

    return &Storage{
        DB: db,
      }, nil
}

//close the database connection
func (storage *Storage) CloseStorage() {
  err := storage.DB.Close()
  if err != nil {
    fmt.Printf("error while closing database: %s", err)
  }  
}