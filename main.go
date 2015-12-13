// YAPC start point.
package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "time"
    "os"
    "os/signal"
    "syscall"
    "sync"
    "io/ioutil"
)

type CacheMemory struct {
    sync.RWMutex
    objectList map[string][]byte
}

var (
    cacheMemory = CacheMemory {
        objectList: make(map[string][]byte),
    }
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

    startTime := time.Now()

    hashKey := CreateHash(r.URL.String())

    // Verify if is a cache hit.
    isHit, statusCode, responseLength, err := verifyIsHit(hashKey, w)
    if err != nil {
        log.Fatal(err)
    }

    if isHit == false {
        // If is a miss, proxy to the upstream.
        statusCode, responseLength, err = proxyUpstreamRequest(w, r, hashKey)
        if err != nil {
            log.Fatal(err)
        }
    }

    elapsedTime := time.Since(startTime)

    logRequest(r, statusCode, responseLength, elapsedTime, isHit)
}

func verifyIsHit(hashKey string, w http.ResponseWriter) (exist bool, statusCode int, upstreamResponseLength int, err error) {

    cacheMemory.RLock()
    body, exist := cacheMemory.objectList[hashKey]
    cacheMemory.RUnlock()

    if !exist {
        return
    }

    // Set Yapc-Cache-State header.
    w.Header().Set("Yapc-Cache-State", "hit")

    // Using io.Write to dump the data object to the downstream.
    _, err = w.Write(body)
    if err != nil {
        return
    }

    statusCode = http.StatusOK
    upstreamResponseLength = len(body)

    return
}

// Proxy a request to the Upstream.
func proxyUpstreamRequest(w http.ResponseWriter, r *http.Request, hashKey string) (statusCode int, upstreamResponseLength int, err error) {

    upstreamResponse, err := http.Get(r.URL.String())
    if err != nil {
        return
    }

    defer upstreamResponse.Body.Close()

    // Set Yapc-Cache-State header.
    w.Header().Set("Yapc-Cache-State", "fetch")

    // Assumption: read the whole object and then send to downstream is a good approach only for small objects.
    // Probably, in another version, the objects will be read in chunks (maybe 4KB).
    body, err := ioutil.ReadAll(upstreamResponse.Body)

    cacheMemory.Lock()
    cacheMemory.objectList[hashKey] = body
    cacheMemory.Unlock()

    _, err = w.Write(body)
    if err != nil {
        return
    }

    statusCode = upstreamResponse.StatusCode
    upstreamResponseLength = len(body)

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
func logRequest(r *http.Request, upstreamStatusCode int, upstreamResponseLength int, elapsedUpstreamGetTime time.Duration, isHit bool) {
    // Log format: [Timestamp] RemoteAddr "URL" Status BytesSent "HttpUserAgent" UpstreamResponseTime
    // Example   : [2006-01-02T15:04:05Z07:00] 87.19.231.27 "http://www.teste.com.br/xyz" 200 12345 "CURL" 123

    var cacheState string
    if isHit {
        cacheState = "HIT"
    } else {
        cacheState = "MISS"
    }

    log.Printf(
        "[%s] %s %s %q %d %d %q %v",
        time.Now().Format(time.RFC3339),
        getIP(r),
        cacheState,
        r.URL.String(),
        upstreamStatusCode,
        upstreamResponseLength,
        r.UserAgent(),
        elapsedUpstreamGetTime,
    )
}
