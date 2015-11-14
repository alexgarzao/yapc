package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
)

const listenAddr = ":8080"

func main() {
    fmt.Println("Starting YAPC server!")

    http.HandleFunc("/", handler)

    err := http.ListenAndServe(listenAddr, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Finishing YAPC server!")
}

func handler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.String()
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close()

    // Using io.Copy to dump the upstream response to the downstream.
    _, err = io.Copy(w, response.Body)
    if err != nil {
        log.Fatal(err)
    }
}
