package yapc_test

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "net/http"
    "io"
    "log"
    "os"
    "net/url"
    "fmt"
    "crypto/tls"
    "crypto/md5"
    "encoding/hex"
)


// Test a valid GET request to an existent object without proxying
func TestValidGetRequestToAnExistentObjectWithoutProxying(t *testing.T) {

    Convey("Given the URL", t, func() {
        url := "http://raw.githubusercontent.com/alexgarzao/yapc/master/README.md"

        Convey("When the object is downloaded", func() {
            statusCode, _ := getObject(url)
            md5sum := "ABC123"

            Convey("The result code must be 200", func() {
                So(statusCode, ShouldEqual, http.StatusOK)
            })

            Convey("The md5sum must be ABC123", func() {
                So(md5sum, ShouldEqual, "ABC123")
            })
        })
    })
}


// Test a valid GET request to an existent object with proxying
func TestValidGetRequestToAnExistentObjectWithProxying(t *testing.T) {

    Convey("Given the URL", t, func() {
        url := "http://raw.githubusercontent.com/alexgarzao/yapc/master/README.md"

        Convey("When the object is downloaded", func() {
            statusCode, _, _ := getObjectFromProxy("http://localhost:8098", url)
            md5sum := "ABC123"

            Convey("The result code must be 200", func() {
                So(statusCode, ShouldEqual, http.StatusOK)
            })

            Convey("The md5sum must be ABC123", func() {
                So(md5sum, ShouldEqual, "ABC123")
            })
        })
    })
}


// Test if the first get is a fetch, and the second is a hit.
func TestIfFirstGetIsFetchAndTheSecondIsHit(t *testing.T) {

    object_url := "http://pbs.twimg.com/profile_images/603610759671611392/JRQtMqMR_normal.png" // Size: 1655 bytes.

    statusCode1, _, cacheState1 := getObjectFromProxy("http://localhost:8098", object_url)
    statusCode2, _, cacheState2 := getObjectFromProxy("http://localhost:8098", object_url)

    Convey("Given the object", t, func() {

        Convey("When the object is downloaded (first time), cache state is a fetch", func() {

            Convey("The result code must be 200", func() {
                So(statusCode1, ShouldEqual, http.StatusOK)
            })

            Convey("Cache state is fetch", func() {
                So(cacheState1, ShouldEqual, "fetch")
            })
        })

        Convey("But when the object is downloaded again (second time), cache state is a hit", func() {

            Convey("The result code must be 200", func() {
                So(statusCode2, ShouldEqual, http.StatusOK)
            })

            Convey("Cache state is fetch", func() {
                So(cacheState2, ShouldEqual, "hit")
            })
        })
    })
}


//--** Private functions **--//


func getObject(url string) (statusCode int, objectLocation string) {
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close()

    // Open a file for writing.
    file, err := os.Create("/tmp/object.download")
    if err != nil {
        log.Fatal(err)
    }

    // Use io.Copy to just dump the response body to the file. This supports huge files.
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()

    statusCode = response.StatusCode
    objectLocation = "/tmp/object.download"

    return
}


func getObjectFromProxy(proxyRawUrl, objectUrl string) (statusCode int, objectLocation string, cacheState string) {
    proxyUrl, err := url.Parse(proxyRawUrl)
    if err != nil {
        fmt.Println("Bad proxy URL", err)
        return
    }

    myClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyUrl),
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //set ssl
        },
    }

    response, err := myClient.Get(objectUrl)
    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close()

    cacheState = response.Header.Get("Yapc-Cache-State")

    // Create the hash key based on the URL.
    objectHash := createHash(objectUrl)

    // Open a file for writing.
    file, err := os.Create("/tmp/fromproxy/" + objectHash)
    if err != nil {
        log.Fatal(err)
    }

    // Use io.Copy to just dump the response body to the file. This supports huge files.
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()

    statusCode = response.StatusCode
    objectLocation = "/tmp/object.download"

    return
}


func createHash(objectUrl string) (string) {
    hash := md5.New()
    hash.Write([]byte (objectUrl))
    md := hash.Sum(nil)
    return hex.EncodeToString(md[:])
}
