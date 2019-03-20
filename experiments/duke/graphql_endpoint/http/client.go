package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	loghttp "github.com/motemen/go-loghttp"
	"github.com/motemen/go-nuts/roundtime"
)

// stole code from here:
// https://medium.com/@clasikas/golang-and-elasticsearch-how-to-get-raw-request-body-json-and-a-little-more-9ee345a9149b
var LoggingClient = &http.Client{
	Transport: &loghttp.Transport{
		LogRequest: func(req *http.Request) {

			var bodyBuffer []byte
			if req.Body != nil {
				bodyBuffer, _ = ioutil.ReadAll(req.Body) // after this operation body will equal 0
				// Restore the io.ReadCloser to request
				req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
			}

			fmt.Println("--------- Elasticsearch ---------")
			fmt.Println("Request URL : ", req.URL)
			fmt.Println("Request Method : ", req.Method)
			fmt.Println("Request Body : ", string(bodyBuffer))
		},
		LogResponse: func(resp *http.Response) {
			ctx := resp.Request.Context()
			if start, ok := ctx.Value(loghttp.ContextKeyRequestStart).(time.Time); ok {
				fmt.Println("Response Status : ", resp.StatusCode)
				fmt.Println("Response Duration : ", roundtime.Duration(time.Now().Sub(start), 2))
			} else {
				fmt.Println("Response Status : ", resp.StatusCode)
			}
			fmt.Println("--------------------------------")
		},
	},
}
