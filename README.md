# IP to Country Service

A Go service that maps IP addresses to their respective countries and cities.

## Project Structure

- `cmd`: Contains the main application entry point
- `pkg/config`: Configuration loading from environment variables
- `pkg/ip2country`: IP to country lookup implementation
- `pkg/ratelimit`: Rate limiting implementation
- `data`: Contains the IP to country mapping data file

## Configuration

The service can be configured using the following environment variables:

- `IP2COUNTRY_DATA_PATH`: Path to the IP-to-country data file (default: `data/ip2country.csv`)
- `IP2COUNTRY_DB_TYPE`: Type of database to use for IP lookups (default: `csv`)
  - Supported values: `csv` (more types will be added in the future)
- `RATE_LIMIT`: The number of requests per second allowed (default: `100`)
- `PORT`: The port on which the service should listen (default: `8080`)

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

The service is designed to be extensible and support different IP-to-country database formats. Currently, only CSV format is implemented, but it's architected to easily add support for other formats like MaxMind DB (MMDB) or database backends like MySQL.

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

## Running the Service

```
go run cmd/main.go
```

With custom settings:

```
IP2COUNTRY_DATA_PATH=/path/to/data.csv IP2COUNTRY_DB_TYPE=csv RATE_LIMIT=50 PORT=9090 go run cmd/main.go
```

## Development with Hot Reload

The project includes hot reload functionality that automatically recompiles and restarts the server when code changes are detected. This makes development faster and more efficient.

### Prerequisites

You'll need to install the Air tool for hot reloading:

```
go install github.com/air-verse/air@latest
```

### Running in Development Mode

To start the server with hot reload enabled:

```
air
```

This will:

1. Watch for changes in your Go files
2. Automatically rebuild when changes are detected
3. Restart the server with the new binary
4. Show you any build errors in real-time

The hot reload configuration is stored in the `.air.toml` file in the project root. You can modify this file to customize the hot reload behavior if needed.
