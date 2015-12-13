package main

import (
    "crypto/md5"
    "encoding/hex"
)

// Create a hash key based on the URL.
func CreateHash(objectUrl string) (string) {
    hash := md5.New()
    hash.Write([]byte (objectUrl))
    md := hash.Sum(nil)
    return hex.EncodeToString(md[:])
}
