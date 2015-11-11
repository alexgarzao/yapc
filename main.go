package main

import (
    "fmt"
    "log"
    "net/http"
)

const listenAddr = "localhost:8080"

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
    fmt.Fprint(w, "Hello, web")
}
