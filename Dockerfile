# --- Build Stage ---
# Use a specific Go version for reproducible builds
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy module files and download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application, creating a static binary
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal image like alpine
RUN CGO_ENABLED=0 go build -o main ./cmd

# --- Final Stage ---
# Use a minimal, non-root image for the final stage
FROM alpine:latest

# It's good practice to run as a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser/app

# Copy only the compiled executable from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Run the executable
CMD ["./main"]
