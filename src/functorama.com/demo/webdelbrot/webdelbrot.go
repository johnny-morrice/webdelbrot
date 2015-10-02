package main

import (
    "flag"
    "log"
    "net/http"
    "path/filepath"
    "fmt"
)

type commandLine struct {
    // Your IP address describes the interface on which we serve
    addr string
    // The port we are to serve upon
    port uint
    // Path to directory containing static files
    static string
}

func parseArguments() commandLine {
    args := commandLine{}
    flag.UintVar(&args.port, "port", 8080, "Port on which to listen")
    flag.StringVar(&args.addr, "addr", "127.0.0.1", "Interface on which to listen")
    flag.StringVar(&args.static, "static", "static", "Path to static files")
    flag.Parse()
    return args
}

func main() {
    args := parseArguments()

    handlers := map[string]func(http.ResponseWriter, *http.Request) {
        "/":                makeIndexHandler(args.static),
        "/style.css":       makeStyleHandler(args.static),
        "/godelbrot.js":    makeJavascriptHandler(args.static),
        "/service":         makeWebserviceHandler(),
    }

    for patt, h := range handlers {
        http.HandleFunc(patt, h)
    }

    serveAddr := fmt.Sprintf("%v:%v", args.addr, args.port)
    httpError := http.ListenAndServe(serveAddr, nil)

    if httpError != nil {
        log.Fatal(httpError)
    }
}

func makeFileHandler(path string, mime string) func(http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", mime)
        http.ServeFile(w, req, path)
    }
}

func makeIndexHandler(static string) func(http.ResponseWriter, *http.Request) {
    return makeFileHandler(filepath.Join(static, "index.html"), "text/html")
}

func makeStyleHandler(static string) func(http.ResponseWriter, *http.Request) {
    return makeFileHandler(filepath.Join(static, "style.css"), "text/css")
}

func makeJavascriptHandler(static string) func(http.ResponseWriter, *http.Request) {
    return makeFileHandler(filepath.Join(static, "godelbrot.js"), "application/javascript")
}