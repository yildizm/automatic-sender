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
	"golang.org/x/time/rate"
)

type MessagePayload struct {
    To      string `json:"to"`
    Content string `json:"content"`
}

type Response struct {
    Message   string `json:"message"`
    MessageID string `json:"messageId"`
}

type WorkerPool struct {
    jobs    chan db.Message
    wg      sync.WaitGroup
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    pool := &WorkerPool{
        jobs:   make(chan db.Message, 10),
        ctx:    ctx,
        cancel: cancel,
    }

    for i := 0; i < numWorkers; i++ {
        pool.wg.Add(1)
        go pool.worker()
    }

    return pool
}

func (p *WorkerPool) AddJob(msg db.Message) {
    select {
    case p.jobs <- msg:
    case <-p.ctx.Done():
    }
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for {
        select {
        case msg := <-p.jobs:
            if err := SendMessage(msg); err != nil {
                log.Println("Error sending message:", err)
            }
        case <-p.ctx.Done():
            return
        }
    }
}

func (p *WorkerPool) Shutdown() {
    p.cancel()
    close(p.jobs)
    p.wg.Wait()
}
var rateLimiter = rate.NewLimiter(1, 5)

func SendMessage(msg db.Message) error {
    if err := rateLimiter.Wait(context.Background()); err != nil {
        return fmt.Errorf("rate limiter error: %w", err)
    }
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
    numWorkers := 5
    pool := NewWorkerPool(numWorkers)
    defer pool.Shutdown()

    ticker := time.NewTicker(2 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        var wg sync.WaitGroup

        for {
            messages, err := db.GetUnsentMessages(10) // Fetch messages in batches
            if err != nil {
                log.Println("Error retrieving messages:", err)
                break
            }

            if len(messages) == 0 {
                break
            }

            for _, msg := range messages {
                pool.AddJob(msg)
            }

            wg.Wait()
        }
    }
}
