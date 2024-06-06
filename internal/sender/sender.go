package sender

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/yildizm/automatic-sender/internal/db"
    "github.com/yildizm/automatic-sender/internal/redis"
    "fmt"
)

type MessagePayload struct {
    To      string `json:"to"`
    Content string `json:"content"`
}

type Response struct {
    Message   string `json:"message"`
    MessageID string `json:"messageId"`
}

func SendMessage(msg db.Message) error {
    payload := MessagePayload{
        To:      msg.Recipient,
        Content: msg.Content,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", "https://webhook.site/6d59e78a-2d8a-4ae3-88d4-971485e87f8f", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-ins-auth-key", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return fmt.Errorf("failed to send message: %s", resp.Status)
    }

    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return err
    }

    sentAt := time.Now()
    if err := db.MarkMessageAsSent(msg.ID, sentAt); err != nil {
        return err
    }

    if err := redis.CacheMessageID(msg.ID, response.MessageID, sentAt); err != nil {
        return err
    }

    return nil
}

func StartSendingMessages() {
    ticker := time.NewTicker(2 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        messages, err := db.GetUnsentMessages(2)
        if err != nil {
            log.Println("Error retrieving messages:", err)
            continue
        }

        for _, msg := range messages {
            if err := SendMessage(msg); err != nil {
                log.Println("Error sending message:", err)
            }
        }
    }
}
