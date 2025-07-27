# Go REST API Starter

## Usage

```sh
go run cmd/api/main.go
```

The server will start on `:8080`.

## Health Check

Visit [http://localhost:8080/health](http://localhost:8080/health) to check the API status.

## Development

- Add new handlers in `/internal/handlers`.
- Register new routes in `/internal/server/server.go`.
- Run tests with `go test ./...`
- Use a `.env` file for configuration (see `.env.example`).

## Docker

To build and run with Docker:

```sh
docker build -t {{ .ProjectName }} .
docker run -p 8080:8080 {{ .ProjectName }}
```
