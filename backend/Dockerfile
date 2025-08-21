# Stage 1: Build Go binary
FROM golang:1.22-bullseye as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# Stage 2: Runtime
FROM debian:bullseye-slim

# Install Chromium and dependencies for rod
RUN apt-get update && \
    apt-get install -y \
    chromium \
    chromium-common \
    ca-certificates \
    fonts-liberation \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libgbm1 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libxss1 \
    libxtst6 \
    libdrm2 \
    libxkbcommon0 \
    libasound2 \
    xdg-utils && \
    rm -rf /var/lib/apt/lists/*

# Set environment variables for rod
ENV ROD_LAUNCHER_BIN=/usr/bin/chromium
ENV DISPLAY=:99
ENV CHROME_BIN=/usr/bin/chromium
ENV CHROMIUM_FLAGS="--no-sandbox --disable-gpu --disable-dev-shm-usage --disable-setuid-sandbox --no-first-run --no-zygote --single-process --disable-background-timer-throttling --disable-backgrounding-occluded-windows --disable-renderer-backgrounding"

WORKDIR /app
COPY --from=builder /app/main .

# Expose Render PORT
ENV PORT=10000
EXPOSE 10000

CMD ["./main"]
