package db

import (
    "time"
    "database/sql"
)

// Message represents a message to be sent or that has been sent
type Message struct {
    ID        int           `json:"id"`
    Content   string        `json:"content"`
    Recipient string        `json:"recipient"`
    SentAt    sql.NullTime  `json:"sentAt"`
}

func GetUnsentMessages(limit int) ([]Message, error) {
    rows, err := DB.Query("SELECT id, content, recipient FROM messages WHERE sent = FALSE LIMIT $1", limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.Content, &msg.Recipient); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }
    return messages, nil
}

func MarkMessageAsSent(id int, sentAt time.Time) error {
    _, err := DB.Exec("UPDATE messages SET sent = TRUE, sent_at = $1 WHERE id = $2", sentAt, id)
    return err
}

func GetSentMessages() ([]Message, error) {
    rows, err := DB.Query("SELECT id, content, recipient, sent_at FROM messages WHERE sent = TRUE")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var msg Message
        if err := rows.Scan(&msg.ID, &msg.Content, &msg.Recipient, &msg.SentAt); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }
    return messages, nil
}
