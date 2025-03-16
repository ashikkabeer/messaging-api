# Messaging API

Messaging Systen built with Golang, PostgreSQL, RabbitMQ

## Features
- RESTful API using Gin framework
- PostgreSQL database for data storage
- RabbitMQ for Message queuing 
- Pagination support for messages
- Input validation

## Prerequisites
- Docker and Docker Compose
- Go 1.24.1
- PostgreSQL
- RabbitMQ

## Quick Start

### Setup & Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/ashikkabeer/messaging-api.git
   cd messaging-api
   ```

2. Start services with Docker:
   ```bash
   docker-compose up -d
   ```

3. Run the application:
   ```bash
   go run main.go
   ```
### API Endpoints
1. `POST /messages`: Create a new message

   Request Body
   ```jsx
        {
          "sender_id": "user123",
          "receiver_id": "user456",
          "content": "Hello, how are you?"
        }
   ```
2. `GET /messages?user1=user123&user2=user456` : Retrieve Chat History

   Response 
   ```jsx
   {
      "has_more": true,
      "messages": [
        {
          "message_id": "f5848b69-7935-48dd-aa67-1deb306fb7cb",
          "sender_id": "user456",
          "receiver_id": "user123",
          "content": "it ok aahn",
          "timestamp": "2025-03-16T07:18:57.917897Z",
          "read": false
        },
        {
          "message_id": "18e90bc6-49f5-41b3-9273-faa921b5b5b4",
          "sender_id": "user456",
          "receiver_id": "user123",
          "content": "it ok aahn",
          "timestamp": "2025-03-16T07:18:57.128185Z",
          "read": false
        }
      ]
      "next_cursor": "MjAyNS0wMy0xNlQwNzoxOToxOS43ODk4NTRaLDUzOTRjY2UxLWM3MWItNDZmMC05NjhjLWJlYThlYjQ3OTVhZQ=="
   }
   ```

3. `GET /messages?user1=user123&user2=user456&cursor=MjAyNS0wMy0xNlQwNzoxODo1OC43OTkyOVosMTgzYTIxMTUtNGM1Zi00NzFlLTk5MjgtNzViODM4ZGYzYjZi` : Paginated Message History Retrieval

   Response 
   ```jsx
   {
      "has_more": true,
      "messages": [
        {
          "message_id": "18e90bc6-49f5-41b3-9273-faa921b5b5b4",
            "sender_id": "user456",
            "receiver_id": "user123",
            "content": "it ok aahn",
            "timestamp": "2025-03-16T07:18:57.128185Z",
            "read": false
        },
        {
          "message_id": "11204c3f-ccba-42fc-8a00-b6216e0e8395",
          "sender_id": "user456",
          "receiver_id": "user123",
          "content": "it ok aahn",
          "timestamp": "2025-03-16T07:18:56.328011Z",
          "read": false
        }
      ]
      "next_cursor": "MjAyNS0wMy0xNlQwNzoxOToxOS43ODk4NTRaLDUzOTRjY2UxLWM3MWItNDZmMC05NjhjLWJlYThlYjQ3OTVhZQ=="
   }
   ```


4. `PATCH /messages/{message_id}/read`: Mark a message as read

   Response
   ```jsx
      {  "status": "read"}
   ```

