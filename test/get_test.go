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
)


// Test a valid GET request to an existent object without proxying
func TestValidGetRequestToAnExistentObjectWithoutProxying(t *testing.T) {
    Convey("Given the URL", t, func() {
        url := "http://github.com/alexgarzao/AWS/raw/master/test/test.ico"

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
        url := "http://github.com/alexgarzao/AWS/raw/master/test/test.ico"

        Convey("When the object is downloaded", func() {
            statusCode, _ := getObjectFromProxy("http://localhost:8080", url)
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


//--** Private functions **--//


func getObject(url string) (statusCode int, objectLocation string) {
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    defer response.Body.Close()

    // Open a file for writing.
    file, err := os.Create("/tmp/test.ico")
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
    objectLocation = "/tmp/test.ico"

    return
}


func getObjectFromProxy(proxyRawUrl, objectUrl string) (statusCode int, objectLocation string) {
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

    // Open a file for writing.
    file, err := os.Create("/tmp/test.ico")
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
    objectLocation = "/tmp/test.ico"

    return
}
