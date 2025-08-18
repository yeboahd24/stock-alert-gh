package httpclient

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	// Global HTTP client instance
	defaultClient *http.Client
	once          sync.Once
)

// GetDefaultClient returns a singleton HTTP client with optimized connection pooling
func GetDefaultClient() *http.Client {
	once.Do(func() {
		defaultClient = createOptimizedClient()
	})
	return defaultClient
}

// createOptimizedClient creates an HTTP client with optimized settings for external API calls
func createOptimizedClient() *http.Client {
	// Create a custom transport with optimized connection pooling
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        100,              // Maximum number of idle connections across all hosts
		MaxIdleConnsPerHost: 10,               // Maximum number of idle connections per host
		MaxConnsPerHost:     50,               // Maximum number of connections per host
		IdleConnTimeout:     90 * time.Second, // How long an idle connection is kept alive
		
		// TCP connection settings
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second, // TCP keep-alive interval
		}).DialContext,
		
		// TLS and HTTP/2 settings
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		
		// Enable HTTP/2
		ForceAttemptHTTP2: true,
		
		// Disable compression for better performance in some cases
		// DisableCompression: false, // Keep compression enabled by default
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Overall request timeout
	}
}

// CreateClientWithTimeout creates a new HTTP client with a specific timeout
// while still using the optimized transport settings
func CreateClientWithTimeout(timeout time.Duration) *http.Client {
	client := GetDefaultClient()
	// Create a new client with the same transport but different timeout
	return &http.Client{
		Transport: client.Transport,
		Timeout:   timeout,
	}
}