package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/yildizm/automatic-sender/internal/db"
	"github.com/yildizm/automatic-sender/internal/redis"
)

type MessagePayload struct {
    To      string `json:"to"`
    Content string `json:"content"`
}

type Response struct {
    Message   string `json:"message"`
    MessageID string `json:"messageId"`
}

func SendMessage(ctx context.Context, msg db.Message) error {
    payload := MessagePayload{
        To:      msg.Recipient,
        Content: msg.Content,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "https://webhook.site/6d59e78a-2d8a-4ae3-88d4-971485e87f8f", bytes.NewBuffer(payloadBytes))
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
    numWorkers := 5
    jobs := make(chan db.Message, 10)
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(jobs, &wg)
    }

    ticker := time.NewTicker(2 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        messages, err := db.GetUnsentMessages(2)
        if err != nil {
            log.Println("Error retrieving messages:", err)
            continue
        }

        for _, msg := range messages {
            jobs <- msg
        }
    }

    close(jobs)
    wg.Wait()
}

func worker(jobs <-chan db.Message, wg *sync.WaitGroup) {
    defer wg.Done()

    for msg := range jobs {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := SendMessage(ctx, msg); err != nil {
            log.Println("Error sending message:", err)
        }
    }
}
