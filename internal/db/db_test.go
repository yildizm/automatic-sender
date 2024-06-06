package db

import (
    "database/sql"
    "testing"
    _ "github.com/lib/pq"
    "os"
)

// TestDBConnection tests the database connection
func TestDBConnection(t *testing.T) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        t.Fatal("DATABASE_URL is not set")
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        t.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        t.Fatalf("Could not ping the database: %v", err)
    }
}

// TestInsertMessages tests the insertion of messages
func TestInsertMessages(t *testing.T) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        t.Fatal("DATABASE_URL is not set")
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        t.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    // Clean up table before testing
    _, err = db.Exec("DELETE FROM messages")
    if err != nil {
        t.Fatalf("Could not clean messages table: %v", err)
    }

    // Insert test messages
    insertMessagesSQL := `INSERT INTO messages (content, recipient, sent_at) VALUES
        ('Test Message 1', '+905551111111', NULL),
        ('Test Message 2', '+905552222222', NULL),
        ('Test Message 3', '+905553333333', NULL);`
    
    _, err = db.Exec(insertMessagesSQL)
    if err != nil {
        t.Fatalf("Could not insert test messages: %v", err)
    }

    // Verify insertion
    rows, err := db.Query("SELECT id, content, recipient, sent_at FROM messages")
    if err != nil {
        t.Fatalf("Could not query messages table: %v", err)
    }
    defer rows.Close()

    count := 0
    for rows.Next() {
        count++
    }

    if count != 3 {
        t.Fatalf("Expected 3 messages, got %d", count)
    }
}
