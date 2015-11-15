// YAPC start point.
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "io"
    "time"
)

const listenAddr = ":8080"

// Main function. Start http handlers.
func main() {
    fmt.Println("Starting YAPC server!")

    http.HandleFunc("/", handler)

    err := http.ListenAndServe(listenAddr, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Finishing YAPC server!")
}

// Main handler. Receive a request and proxy them to the upstream.
func handler(w http.ResponseWriter, r *http.Request) {
    upstreamStatusCode, upstreamResponseLength, elapsedUpstreamGetTime, err := proxyUpstreamRequest(w, r)
    if err != nil {
        log.Fatal(err)
    }

    logRequest(r, upstreamStatusCode, upstreamResponseLength, elapsedUpstreamGetTime)
}

// Proxy a request to the Upstream.
func proxyUpstreamRequest(w http.ResponseWriter, r *http.Request) (statusCode int, upstreamResponseLength int64, elapsedTime time.Duration, err error) {
    startTime := time.Now()

    upstreamResponse, err := http.Get(r.URL.String())
    if err != nil {
        return
    }

    defer upstreamResponse.Body.Close()

    // Using io.Copy to dump the upstreamResponse to the downstream.
    upstreamResponseLength, err = io.Copy(w, upstreamResponse.Body)
    if err != nil {
        log.Fatal(err)
    }

    elapsedTime = time.Since(startTime)
    statusCode = upstreamResponse.StatusCode

    return
}

// Get IP from a http request.
func getIP(r *http.Request) string {
    if ipProxy := r.Header.Get("X-Forwarded-For"); len(ipProxy) > 0 {
        return ipProxy
    }

    ip, _, _ := net.SplitHostPort(r.RemoteAddr)

    return ip
}

// Log a http request.
func logRequest(r *http.Request, upstreamStatusCode int, upstreamResponseLength int64, elapsedUpstreamGetTime time.Duration) {
    // Log format: [Timestamp] RemoteAddr "URL" Status BytesSent "HttpUserAgent" UpstreamResponseTime
    // Example   : [2006-01-02T15:04:05Z07:00] 87.19.231.27 "http://www.teste.com.br/xyz" 200 12345 "CURL" 123
    log.Printf(
        "[%s] %s %q %d %d %q %v",
        time.Now().Format(time.RFC3339),
        getIP(r),
        r.URL.String(),
        upstreamStatusCode,
        upstreamResponseLength,
        r.UserAgent(),
        elapsedUpstreamGetTime,
    )
}
