package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
    var err error
    DB, err = sql.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatal(err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal(err)
    }

    createTable()
}

func createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY,
        content TEXT NOT NULL CHECK (char_length(content) <= 160),
        recipient VARCHAR(20) NOT NULL,
        sent BOOLEAN DEFAULT FALSE,
        sent_at TIMESTAMP
    );`
    _, err := DB.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}
