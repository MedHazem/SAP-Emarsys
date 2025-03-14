FROM golang:1.21 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o due-date-server .

# Create a small final image
FROM gcr.io/distroless/base-debian12

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/due-date-server .

# Expose the application's port
EXPOSE 8080

# Command to run the executable
CMD ["./due-date-server"]
