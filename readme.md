# Vortex

Vortex is a super simple key-value store written in Go. It provides a RESTful API to create, read, update, and delete key-value pairs.

The application uses an in-memory map to store the key-value pairs. For restoring the state after a restart, the application writes the map to a file on disk. When the application starts, it reads the file and restores the map.

## Prerequisites

- Go 1.16 or higher
- curl (for testing)

## Running the App

### Configuration

The application uses the following environment variables for configuration:

- `TRANSACTION_LOG_FILE`: The filepath to write the transaction log to. Default is `./transactions.log`.

### Run 

To run the app, navigate to the project directory and use the `go run` command:

```bash
go run .
```

## Usage 

You can interact with the application using curl. The following commands are available:

### Create a new key-value pair

``` bash
curl -X PUT -H "Content-Type: application/json" -d '{"value":"your_value"}' http://localhost:8080/v1/kv/your_key
```

Replace "your_key" and "your_value" with the actual key and value you want to put.


### Get a key-value pair

```bash
curl -X GET http://localhost:8080/v1/kv/your_key
```

Replace "your_key" with the actual key you want to get.

### Delete a key-value pair

```bash
curl -X DELETE http://localhost:8080/v1/kv/your_key
```

Replace "your_key" with the actual key you want to delete.

## Testing

### Run the tests

To run the tests, navigate to the project directory and use the `go test` command:

```bash
go test ./... -race
```

### Test Coverage

To run the tests with coverage, navigate to the project directory and use the `go test` command with the `-coverprofile` flag:

```bash
go test ./... -race -coverprofile=coverage.out
go tool cover -html=coverage.out
```
