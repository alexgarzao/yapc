// YAPC start point.
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "io"
    "time"
    "os"
    "os/signal"
    "syscall"
)

const listenAddr = ":8098"

// Main function. Start http handlers.
func main() {
    fmt.Println("Starting YAPC server!")

    handleSignals()

    http.HandleFunc("/", handler)

    err := http.ListenAndServe(listenAddr, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Finishing YAPC server!")
}

// Handle all signals to the process.
func handleSignals() {
    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel,
        syscall.SIGINT,
        syscall.SIGHUP,
        syscall.SIGTERM,
        syscall.SIGQUIT,
        syscall.SIGABRT,
    )

    go func() {
        for {
            signal := <- signalChannel

            fmt.Println()
            fmt.Println(signal)

            switch signal {
            case syscall.SIGINT:
                // Handle SIGINT
                os.Exit(0)
            case syscall.SIGHUP:
                // Handle SIGHUP
            case syscall.SIGTERM:
                // Handle SIGTERM
            case syscall.SIGQUIT:
                // Handle SIGQUIT
            case syscall.SIGABRT:
                // Handle SIGABRT
            }
        }
    }()
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

    // Set Yapc-Cache-State header.
    w.Header().Set("Yapc-Cache-State", "fetch")

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
