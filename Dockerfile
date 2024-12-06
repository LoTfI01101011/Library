# Stage 1: Build
FROM golang:1.22.2 AS build

WORKDIR /app

# Copy dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application files
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8000

# CMD to run migration first and then start the application
CMD ["sh", "-c", "go run migrate/migrate.go migrate && ./main"]
