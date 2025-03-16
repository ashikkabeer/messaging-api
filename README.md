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

## Resources Used
- [https://medium.easyread.co/how-to-do-pagination-in-postgres-with-golang-in-4-common-ways-12365b9fb528](https://medium.easyread.co/how-to-do-pagination-in-postgres-with-golang-in-4-common-ways-12365b9fb528)
- [https://www.rabbitmq.com/tutorials/tutorial-one-go](https://www.rabbitmq.com/tutorials/tutorial-one-go)