package main

import (
    "io"

    "net/http"

    "github.com/johnny-morrice/godelbrot/restclient"
)

func defaultconfig() restclient.Config {
    config := restclient.Config{}
    config.Addr = __ADDR
    config.Port = __PORT
    config.Prefix = __PREFIX
    config.Ticktime = __TICKTIME
    config.Debug = __DEBUG
    config.Http = (*goHttp)(&http.Client{})
    return config
}

var __restclient *restclient.Client
func getclient() *restclient.Client {
    if __restclient == nil {
        config := defaultconfig()
        __restclient = restclient.New(config)
    }

    return __restclient
}

// Interface implementation below for restclient access to gopherjs-friendly net/http.

type goHttp http.Client
var _ restclient.HttpClient = (*goHttp)(nil)

func (goh *goHttp) Get(url string) (restclient.HttpResponse, error) {
    r, err := (*http.Client)(goh).Get(url)
    return (*goHttpResponse)(r), err
}

func (goh *goHttp) Post(url, ctype string, content io.Reader) (restclient.HttpResponse, error) {
    r, err := (*http.Client)(goh).Post(url, ctype, content)
    return (*goHttpResponse)(r), err
}

type goHttpResponse http.Response
var _ restclient.HttpResponse = (*goHttpResponse)(nil)

func (r *goHttpResponse) GetBody() io.ReadCloser {
    return (*http.Response)(r).Body
}
func (r *goHttpResponse) GetStatusCode() int {
    return (*http.Response)(r).StatusCode
}
func (r *goHttpResponse) GetStatus() string {
    return (*http.Response)(r).Status
}
func (r *goHttpResponse) GetHeader() map[string][]string {
    return (*http.Response)(r).Header
}
func (r *goHttpResponse) Write(w io.Writer) error {
    return (*http.Response)(r).Write(w)
}