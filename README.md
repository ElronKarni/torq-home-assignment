# IP to Country Service

A Go service that maps IP addresses to their respective countries and cities.

![System Data Flow Diagram](/data-flow-diagram.png)

## Installation and Running the Service

### Prerequisites

- Go 1.24 or higher
- Optional: Docker and Docker Compose for containerized deployment
- Optional: Air for hot reloading during development

### Setup

1. Clone the repository:

   ```
   git clone https://github.com/ElronKarni/torq-home-assignment.git
   ```

2. Install dependencies:

   ```
   go mod download
   ```

### Running the Service

There are several ways to run the application:

#### 1. Using Go directly:

```
go run cmd/main.go
```

#### 2. Using Make:

```
make run
```

#### 3. Using Docker Compose:

```
make docker-up
```

Or directly:

```
docker-compose up --build -d
```

#### 5. Development with Hot Reload:

To start the server with hot reload enabled (requires Air to be installed):

```
air
```

## Testing

Run all tests:

```
make test
```

Run tests with coverage report:

```
make test-coverage
```

Test the rate limiting functionality:

```
make rate-limit-test
```

## Project Structure

- `cmd`: Contains the main application entry point
- `internal/config`: Configuration loading from environment variables
- `internal/ip2country`: IP to country lookup implementation
- `internal/middleware`: HTTP middleware implementations
- `internal/handlers`: HTTP request handlers
- `internal/routes`: API route definitions
- `internal/utils`: Utility functions
- `pkg/ratelimit`: Rate limiting implementation
- `data`: Contains the IP to country mapping data file

## Configuration

The service can be configured using the following environment variables:

- `IP2COUNTRY_DB_TYPE`: Type of database to use for IP lookups (default: `csv`)
  - Supported values: `csv` (more types will be added in the future)
- `RATE_LIMIT`: The number of requests per second allowed (default: `50`)
- `PORT`: The port on which the service should listen (default: `8080`)
- `CSV_DATA_PATH`: Path to the CSV data file when using CSV database type (default: `data/ip2country.csv`)
- `MONGO_URI`: MongoDB connection URI when using MongoDB database type (default: `mongodb://localhost:27017`)
- `REDIS_ADDR`: Redis server address when using Redis database type (default: `localhost:6379`)
- `ALLOWED_ORIGINS`: Comma-separated list of allowed origins for CORS (default: `http://localhost:3000`)

## Data File Format

The IP to country data file is a CSV file with the following format:

```
ip,city,country
```

Example:

```
1.1.1.1,Sydney,Australia
8.8.8.8,Mountain View,United States
```

## Extensibility

The service is designed to be extensible and support different IP-to-country database formats. Currently, only CSV format is implemented, but it's architected to easily add support for other formats like Redis or MongoDB database and more...

To use a different database type, simply set the `IP2COUNTRY_DB_TYPE` environment variable to the desired type. New types can be added by implementing the `ip2country.Service` interface.

## API Endpoints

### GET /v1/find-country

Returns the country and city for a given IP address.

**Query Parameters**:

- `ip`: The IP address to look up

**Example Request**:

```
GET /v1/find-country?ip=1.1.1.1
```

**Example Success Response (200 OK)**:

```json
{
  "country": "Australia",
  "city": "Sydney"
}
```

**Error Responses**:

- 400 Bad Request - Missing or invalid IP address

```json
{
  "error": "Invalid IP address"
}
```

- 404 Not Found - IP address not found in database

```json
{
  "error": "IP address not found"
}
```

- 429 Too Many Requests - Rate limit exceeded

```json
{
  "error": "Too many requests"
}
```

## Rate Limiting

The service implements a rate limiter that restricts the number of requests per second based on the `RATE_LIMIT` environment variable. If the rate limit is exceeded, the service returns a 429 HTTP status code with the following response:

```json
{
  "error": "Too many requests"
}
```
