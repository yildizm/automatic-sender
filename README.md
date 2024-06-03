# Automatic Message Sending System

## Overview
This project is a Go-based system that automatically sends messages retrieved from a database to specified recipients every 2 minutes.

## Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/dl/) (version 1.20 or later)

## Setup
1. **Clone the Repository**:
   ```sh
   git clone https://github.com/yourusername/automatic-sender.git
   cd automatic-sender

1. **Set Up Environment Variables:**:
   ```
   DATABASE_URL=postgres://user:password@db:5432/messagesdb?sslmode=disable
    REDIS_URL=redis://redis:6379
   ```

1. **Run Docker Compose**:
   ```docker-compose up --build```

### Docker
To run the application using Docker, use the following commands:

```sh
docker-compose up --build

Usage
API Endpoints
Start Automatic Message Sending:

URL: POST /start
Description: Starts the automatic message sending process.
Response: 200 OK or 400 Already sending messages
Stop Automatic Message Sending:

URL: POST /stop
Description: Stops the automatic message sending process.
Response: 200 OK or 400 Not currently sending messages
Retrieve Sent Messages:

URL: GET /sent-messages
Description: Retrieves a list of sent messages.
Response: 200 OK with a list of messages or 500 Internal Server Error

Accessing Swagger UI
http://localhost:8080/swagger/index.html
