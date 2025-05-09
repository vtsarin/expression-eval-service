# Expression Evaluation Service

A Go service that evaluates mathematical expressions with support for batch processing.

## Features

- Single expression evaluation
- Batch expression evaluation
- Expression history with pagination
- Concurrent processing
- Error handling and logging
- Rate limiting
- CORS support
- Health check endpoint

## Prerequisites

- Go 1.21 or later
- Docker (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/expression-eval-service.git
cd expression-eval-service
```

2. Install dependencies:
```bash
go mod download
```

## Running the Service

### Local Development

```bash
go run main.go
```

The service will start on port 8080.

### Using Docker

Build the Docker image:
```bash
docker build -t expression-eval-service .
```

Run the container:
```bash
docker run -p 8080:8080 expression-eval-service
```

## API Endpoints

### Single Expression Evaluation

```bash
curl -X POST http://localhost:8080/api/evaluate/single \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "3 + 4"
  }'
```

### Batch Expression Evaluation

```bash
curl -X POST http://localhost:8080/api/evaluate/batch \
  -H "Content-Type: application/json" \
  -d '{
    "expressions": [
      "3 + 4",
      "10 / 2",
      "(2 + 3) * 5",
      "3 / 0",
      "4 + * 2"
    ]
  }'
```

### Get History

```bash
curl "http://localhost:8080/api/evaluate/history?page=1&pageSize=10"
```

## Configuration

The service can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `RATE_LIMIT`: Requests per second (default: 100)
- `ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins (default: "*")

## Error Handling

The service provides detailed error messages for:
- Invalid expressions
- Division by zero
- Syntax errors
- Rate limit exceeded
- Invalid requests

## Logging

The service uses structured logging with Zap, including:
- Request/response logging
- Error logging
- Performance metrics
- Correlation IDs

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 