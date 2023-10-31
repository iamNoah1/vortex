# Vortex

Vortex is a super simple key-value store written in Go.

## Prerequisites

- Go 1.16 or higher
- curl (for testing)

## Running the App

To run the app, navigate to the project directory and use the `go run` command:

```bash
go run .
```

## Usage 

You can interact with the application using curl. The following commands are available:

### Create a new key-value pair

``` bash
curl -X PUT -H "Content-Type: application/json" -d '{"key":"your_key", "value":"your_value"}' http://localhost:8080/v1/kv/your_key
```

Replace "your_key" and "your_value" with the actual key and value you want to put.


## Get a key-value pair

```bash
curl -X GET http://localhost:8080/v1/kv/your_key
```

Replace "your_key" with the actual key you want to get.

