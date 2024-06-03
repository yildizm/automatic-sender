package main

import (
    "log"
    "net/http"
    "os"
    "time"
    "github.com/yildizm/automatic-sender/internal/api"
    "github.com/yildizm/automatic-sender/internal/db"
    "github.com/yildizm/automatic-sender/internal/redis"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/yildizm/automatic-sender/docs" // This is required for swag to find your docs
)

// @title Automatic Message Sending System API
// @version 1.0
// @description This is an API for an automatic message sending system.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-ins-auth-key
func main() {
    log.Println("Starting application")

    // Database connection
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }
    db.InitDB(dbURL)
    log.Println("Database connected")

    // Redis connection
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        log.Fatal("REDIS_URL is not set")
    }
    redis.InitRedis(redisURL)
    log.Println("Redis connected")

    // Router setup
    router := mux.NewRouter()

    // API routes
    router.HandleFunc("/message-sending", api.MessageSendingHandler).Methods("POST")
    router.HandleFunc("/sent-messages", api.SentMessagesHandler).Methods("GET")

    // Swagger setup
    router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    // Middleware for CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "x-ins-auth-key"},
        AllowCredentials: true,
        MaxAge:           300,
    })

    // Start HTTP server
    srv := &http.Server{
        Handler:      c.Handler(router),
        Addr:         "0.0.0.0:8080",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Println("Starting server on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
