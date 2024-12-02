# Task Processor Application

This is a task processor application that consumes tasks from a RabbitMQ message queue, processes them, and stores the results in MongoDB. The tasks can be of various types, such as health checks or email processing. The application tracks the status of each task, including retries, failures, and timestamps.

## Features

- **Task Queue Processing**: Consumes tasks from a RabbitMQ queue.
- **Task Types**: Supports different types of tasks like health checks and email processing.
- **Task Status Tracking**: Updates task status (e.g., in-progress, completed, failed) and stores results in MongoDB.
- **Retry Logic**: Retries task execution up to a configurable number of retries.
- **Graceful Shutdown**: Handles system shutdown gracefully by stopping the workers and cleaning up resources.
- **Task Result Storage**: Stores task results, including success or failure status, messages, and timestamps in MongoDB.

## Architecture

- **RabbitMQ**: Used as the message broker to queue tasks for processing.
- **MongoDB**: Used as the database to store task status and results.
- **Go**: The application is written in Go, utilizing the `streadway/amqp` library to interface with RabbitMQ and the `mongo-driver` for MongoDB operations.

## Requirements

- Docker and Docker Compose (for running RabbitMQ, MongoDB, and the app together)
- Go 1.18+ for running the application locally
- MongoDB running (either locally or via Docker)
- RabbitMQ running (either locally or via Docker)

## Environment Variables

The application requires the following environment variables to function properly:

- `RABBITMQ_URL`: URL for RabbitMQ connection (default: `amqp://guest:guest@rabbitmq:5672/`).
- `MONGODB_URL`: MongoDB connection URL (default: `mongodb://mongodb:27017/`).
- `MAX_RETRIES`: Maximum number of retries for a failed task (default: `3`).
- `API_PORT`: The port the producer will receive task request on (default: `8080`).

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/millroy094/task-processor.git
cd task-processor
