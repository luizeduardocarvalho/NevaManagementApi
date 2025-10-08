FROM golang:1.23-alpine

# Install required packages for development
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o labflux-server cmd/main.go

EXPOSE 8080

# Run the server
CMD ["./labflux-server"]