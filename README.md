# Go Expression Evaluation Service

A robust Go service built with Gin framework that provides expression evaluation capabilities with history tracking. The service features structured logging, error handling, and middleware for request tracking.

## Features

- Expression evaluation with support for:
  - Basic arithmetic operations (+, -, *, /)
  - Parentheses for grouping
  - Floating-point numbers
- Evaluation history tracking
- Thread-safe operations
- Structured logging with Zap
- Request correlation ID tracking
- Request latency monitoring
- Health check endpoint
- Docker support

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

## Project Structure

```
.
├── controllers/         # HTTP request handlers
├── errors/             # Custom error definitions
├── evaluator/          # Expression parsing and evaluation
├── interfaces/         # Interface definitions
├── middlewares/        # HTTP middleware components
├── services/           # Business logic
├── utils/              # Utility functions
├── Dockerfile         # Multi-stage Docker build
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── main.go            # Application entry point
└── README.md          # This file
```

## Setup and Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-service
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the service:
   ```bash
   go run main.go
   ```

The service will start on port 8080.

## Docker Deployment

Build and run using Docker:

```bash
# Build the image
docker build -t go-service .

# Run the container
docker run -p 8080:8080 go-service
```

## API Documentation

### 1. Evaluate Expression

Evaluates a mathematical expression and stores the result in history.

**Endpoint:** `POST /api/evaluate`

**Request Body:**
```json
{
    "expression": "3 + (2 * 4)"
}
```

**Success Response (200 OK):**
```json
{
    "status": "success",
    "message": "Expression evaluated successfully",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "result": 11
    }
}
```

**Error Response (400 Bad Request):**
```json
{
    "status": "error",
    "message": "Invalid expression",
    "error": {
        "code": "E4007701",
        "message": "Invalid expression",
        "details": "unexpected token: )"
    }
}
```

### 2. Get Evaluation History

Retrieves the history of all expression evaluations.

**Endpoint:** `GET /api/history`

**Success Response (200 OK):**
```json
{
    "status": "success",
    "message": "History retrieved successfully",
    "data": [
        {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "expression": "3 + (2 * 4)",
            "result": 11,
            "timestamp": "2024-03-14T12:00:00Z"
        }
    ]
}
```

### 3. Health Check

Checks the health status of the service.

**Endpoint:** `GET /api/health`

**Success Response (200 OK):**
```json
{
    "status": "success",
    "message": "Service is healthy",
    "data": {
        "status": "UP",
        "timestamp": "2024-03-14T12:00:00Z"
    }
}
```

## Testing the API

Using curl:

```bash
# Evaluate an expression
curl -X POST http://localhost:8080/api/evaluate \
  -H "Content-Type: application/json" \
  -d '{"expression": "3 + (2 * 4)"}'

# Get evaluation history
curl http://localhost:8080/api/history

# Check service health
curl http://localhost:8080/api/health
```

## Logging

The service uses structured logging with Zap. Logs include:
- Request correlation IDs
- Request latency
- Operation timing
- Error details
- Evaluation history tracking

## Error Handling

The service implements a custom error handling system with:
- Structured error responses
- Error codes for different scenarios
- Detailed error messages
- Error tracking in evaluation history

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 