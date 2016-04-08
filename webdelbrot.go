package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

type commandLine struct {
    // Your IP address describes the interface on which we serve
    addr string
    // The port we are to serve upon
    port uint
    // Path to directory containing static files
    static string
    // Number of concurrent render tasks
    jobs uint
}

func parseArguments() commandLine {
    args := commandLine{}
    flag.UintVar(&args.port, "port", 8080, "Port on which to listen")
    flag.StringVar(&args.addr, "bind", "127.0.0.1", "Interface on which to listen")
    flag.StringVar(&args.static, "static", "webdelbrot-static", "Path to static files")
    flag.UintVar(&args.jobs, "jobs", 1, "Concurrent render tasks")
    flag.Parse()
    return args
}

func main() {
    args := parseArguments()

    var input io.Reader = os.Stdin
    baseinfo, readErr := lib.ReadInfo(input)
    if readErr != nil {
        log.Fatal("Error reading info: ", readErr)
    }

    if !foundFiles(args) {
        log.Fatal("Could not find files at: %v", args.static)
    }

    // Handle fractal requests
    const prefix = "/service"
    router := makeWebservice(baseinfo, args.jobs, prefix)
    http.Handle(prefix, router)

    // Handle static files
    handleFiles(args.static)

    // Run webserver
    serveAddr := fmt.Sprintf("%v:%v", args.addr, args.port)
    httpError := http.ListenAndServe(serveAddr, nil)

    if httpError != nil {
        log.Fatal(httpError)
    }
}

func foundFiles(args commandLine) bool {
    _, err := os.Stat(args.static)
    return err == nil
}

func handleFiles(root string) {
    http.HandleFunc("/", makeIndexHandler(root))

    staticFiles := map[string]string {
        "style.css": "text/css",
        "godelbrot.js": "application/javascript",
        "history.js": "application/javascript",
        "mandelbrot.js": "application/javascript",
        "complex.js": "application/javascript",
        "image.js": "application/javascript",
        "zoom.js": "application/javascript",
        "favicon.ico": "image/x-icon",
        "small-logo.png": "image.png",
    }

    for filename, mime := range staticFiles {
        path := filepath.Join(root, filename)
        http.HandleFunc(path, makeFileHandler(path, mime))
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