package api

import (
    "encoding/json"
    "net/http"
    "sync"
    "github.com/yildizm/automatic-sender/internal/db"
    "github.com/yildizm/automatic-sender/internal/sender"
    "log"
)

var sending = false
var mu sync.Mutex

// MessageSendingHandler handles starting and stopping the automatic message sending process
// @Summary Start or stop automatic message sending
// @Description Starts or stops the automatic message sending process
// @Tags messages
// @Accept  json
// @Produce  json
// @Param action query string true "Action to perform" Enums(start, stop)
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid action or already in desired state"
// @Router /message-sending [post]
func MessageSendingHandler(w http.ResponseWriter, r *http.Request) {
    action := r.URL.Query().Get("action")
    if action == "" {
        http.Error(w, "Action is required", http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    switch action {
    case "start":
        if sending {
            http.Error(w, "Already sending messages", http.StatusBadRequest)
            return
        }
        go sender.StartSendingMessages()
        sending = true
        log.Println("Started automatic message sending")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    case "stop":
        if !sending {
            http.Error(w, "Not currently sending messages", http.StatusBadRequest)
            return
        }
        // Logic to stop the message sending process (not implemented in this example)
        sending = false
        log.Println("Stopped automatic message sending")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    default:
        http.Error(w, "Invalid action", http.StatusBadRequest)
    }
}

// SentMessagesHandler retrieves a list of sent messages
// @Summary Retrieve sent messages
// @Description Retrieves a list of sent messages
// @Tags messages
// @Accept  json
// @Produce  json
// @Success 200 {array} db.Message
// @Failure 500 {string} string "Error retrieving sent messages"
// @Router /sent-messages [get]
func SentMessagesHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to retrieve sent messages")
    messages, err := db.GetSentMessages()
    if err != nil {
        log.Println("Error retrieving sent messages:", err)
        http.Error(w, "Error retrieving sent messages", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}
