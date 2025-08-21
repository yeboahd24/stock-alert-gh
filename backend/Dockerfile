FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies including Chrome
RUN apk --no-cache add ca-certificates sqlite chromium nss freetype freetype-dev harfbuzz ttf-freefont

# Set Chrome environment variables
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Create data directory for SQLite
RUN mkdir -p data

# Expose port
EXPOSE 10000

# Run the application
CMD ["./server"]