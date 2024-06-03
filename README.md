# Automatic Message Sending System

## Overview
This project is a Go-based system that automatically sends messages retrieved from a database to specified recipients every 2 minutes.

## Features

- Automatically sends messages retrieved from the database every 2 minutes.
- Ensures that each message is sent only once.
- Caches sent message IDs and timestamps in Redis.
- Provides API endpoints to start and stop the automatic message sending process.
- Provides an API endpoint to retrieve a list of sent messages.
- API documentation using Swagger.


## Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/dl/) (version 1.20 or later)

## Setup
1. **Clone the Repository**:
   ```sh
   git clone https://github.com/yourusername/automatic-sender.git
   cd automatic-sender

2. **Set Up Environment Variables:**:
   ```
   DATABASE_URL=postgres://user:password@db:5432/messagesdb?sslmode=disable
    REDIS_URL=redis://redis:6379
   ```
### Docker
To run the application using Docker, use the following commands:

```sh
docker-compose up --build

## Usage

### API Endpoints

- **Start/Stop Automatic Message Sending**:
  - **URL**: `POST /message-sending?action=start`
  - **URL**: `POST /message-sending?action=stop`
  - **Description**: Starts or stops the automatic message sending process.
  - **Response**: 
    - `200 OK` - Action was successful.
    - `400 Invalid action or already in desired state` - If the action is invalid or the system is already in the desired state.

- **Retrieve Sent Messages**:
  - **URL**: `GET /sent-messages`
  - **Description**: Retrieves a list of sent messages.
  - **Response**: 
    - `200 OK` - Returns a list of sent messages.
    - `500 Internal Server Error` - If there is an error retrieving the messages.

### Accessing Swagger UI

Swagger UI is available at:
```
http://localhost:8080/swagger/index.html
