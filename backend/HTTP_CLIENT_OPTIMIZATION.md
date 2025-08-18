# HTTP Client Optimization

## Overview
This document describes the HTTP client optimizations implemented to improve external API call performance through connection pooling and reuse.

## Key Improvements

### 1. Singleton HTTP Client Pattern
- Created a shared HTTP client instance using the singleton pattern
- Eliminates the overhead of creating new clients for each request
- Ensures consistent configuration across the application

### 2. Optimized Connection Pooling
The HTTP transport is configured with the following optimizations:

```go
MaxIdleConns:        100   // Maximum idle connections across all hosts
MaxIdleConnsPerHost: 10    // Maximum idle connections per host
MaxConnsPerHost:     50    // Maximum connections per host
IdleConnTimeout:     90s   // How long idle connections are kept alive
```

### 3. TCP Connection Optimization
```go
DialContext: (&net.Dialer{
    Timeout:   10s,  // Connection timeout
    KeepAlive: 30s,  // TCP keep-alive interval
}).DialContext
```

### 4. TLS and HTTP/2 Optimization
```go
TLSHandshakeTimeout:   10s
ResponseHeaderTimeout: 10s
ExpectContinueTimeout: 1s
ForceAttemptHTTP2:     true
```

## Performance Benefits

### Before Optimization
- New HTTP client created for each request
- New TCP connections established for each request
- DNS lookups performed for each request
- TLS handshakes performed for each request

### After Optimization
- Single HTTP client instance reused across all requests
- TCP connections pooled and reused
- DNS lookups cached and reused
- TLS sessions reused when possible
- HTTP/2 multiplexing enabled for better performance

## Usage

### Default Client
```go
import "shares-alert-backend/internal/httpclient"

client := httpclient.GetDefaultClient()
resp, err := client.Get("https://api.example.com/data")
```

### Custom Timeout Client
```go
client := httpclient.CreateClientWithTimeout(20 * time.Second)
resp, err := client.Get("https://api.example.com/data")
```

## Implementation Details

### Files Modified
- `backend/internal/httpclient/client.go` - New shared HTTP client package
- `backend/internal/services/stock_service.go` - Updated to use shared client

### Connection Reuse
The optimized client will:
1. Reuse existing connections when making requests to the same host
2. Keep connections alive for 90 seconds after use
3. Maintain up to 10 idle connections per host
4. Support up to 50 concurrent connections per host

## Monitoring
To monitor the effectiveness of connection pooling, you can:
1. Check connection metrics in production
2. Monitor response times for external API calls
3. Observe reduced DNS lookup frequency
4. Track TLS handshake reduction

## Expected Performance Improvements
- **Reduced latency**: 50-200ms savings per request due to connection reuse
- **Lower CPU usage**: Fewer TLS handshakes and connection establishments
- **Better throughput**: HTTP/2 multiplexing allows multiple requests over single connection
- **Reduced network overhead**: Fewer TCP connection establishments