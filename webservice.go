package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "image/png"
    "log"
    "net/http"
    "strconv"
    "strings"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

type WebCommand string
const (
    render = WebCommand("render")
    displayImage = WebCommand("displayImage")
)

type WebRenderParams struct {
    ImageWidth uint
    ImageHeight uint
    RealMin float64
    RealMax float64
    ImagMin float64
    ImagMax float64
}

type GodelbrotPacket struct {
    Command WebCommand
    Render WebRenderParams
}

const godelbrotHeader string = "X-Godelbrot-Packet"
const godelbrotGetParam string = "godelbrotPacket"

type queueCommand uint

const (
    queueRender = queueCommand(iota)
    queueStop = queueCommand(iota)
)

type renderQueueItem struct {
    command queueCommand
    w http.ResponseWriter
    req *http.Request
    complete chan<- bool
}

func launchRenderService(desc *lib.Info) (func(http.ResponseWriter, *http.Request), chan<- renderQueueItem) {
    input := make(chan renderQueueItem)

    go handleRenderRequests(desc, input)

    return httpChanWriter(input), input
}

func httpChanWriter(input chan<- renderQueueItem) func(http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, req *http.Request) {
        done := make(chan bool)
        input <- renderQueueItem{
            command: queueRender,
            w: w,
            req: req,
            complete: done,
        }
        // Block until rendering is complete
        <- done
    }
}

type webservice struct {
    queueItem renderQueueItem
    desc *lib.Info
}

func handleRenderRequests(desc *lib.Info, input <-chan renderQueueItem) {
    serv := webservice{}
    serv.desc = desc

    for queue := range input {
        switch queue.command {
        case queueRender:
            serv.queueItem = queue
            serv.render()
        case queueStop:
            break
        default:
            panic(fmt.Sprintf("Unknown queueCommand: %v", queue.command))
        }
    }
}

func (serv webservice) render() {
    jsonPacket := serv.queueItem.req.URL.Query().Get(godelbrotGetParam)

    if len(jsonPacket) == 0 {
        http.Error(serv.queueItem.w, fmt.Sprintf("No data found in parameter '%v'", godelbrotHeader), 400)
        serv.queueItem.complete <- true
        return
    }

    dec := json.NewDecoder(strings.NewReader(jsonPacket))
    userPacket := GodelbrotPacket{}
    jsonError := dec.Decode(&userPacket)

    if jsonError != nil {
        http.Error(serv.queueItem.w, fmt.Sprintf("Invalid JSON packet: %v", jsonError), 400)
        serv.queueItem.complete <- true
        return
    }

    args := userPacket.Render

    if args.ImageWidth == 0 || args.ImageHeight == 0 {
        http.Error(serv.queueItem.w, "ImageHeight and ImageWidth cannot be 0", 422)
        serv.queueItem.complete <- true
        return
    }

    desc, cerr := serv.configure(args)

    if cerr != nil {
        msg := fmt.Sprintf("Error in configuration: %v", cerr)
        http.Error(serv.queueItem.w, msg, 400)
        serv.queueItem.complete <- true
        return
    }

    pic, renderError := lib.Render(desc)

    if renderError != nil {
        log.Fatal(fmt.Sprintf("Render error: %v", renderError))
    }

    buff := bytes.Buffer{}
    pngError := png.Encode(&buff, pic)

    if pngError != nil {
        http.Error(serv.queueItem.w, fmt.Sprintf("Error encoding PNG: %v", pngError), 500)
        serv.queueItem.complete <- true
        return
    }

    // Craft response
    responsePacket := GodelbrotPacket{
        Command: displayImage,
    }

    responseHeaderPacket, marshalError := json.Marshal(responsePacket)

    if marshalError != nil {
        http.Error(serv.queueItem.w, fmt.Sprintf("Error marshalling response header: %v", marshalError), 500)
        serv.queueItem.complete <- true
        return
    }

    // Respond to the request
    serv.queueItem.w.Header().Set("Content-Type", "image/png")
    serv.queueItem.w.Header().Set(godelbrotHeader, string(responseHeaderPacket))

    // Write image buffer as http response
    serv.queueItem.w.Write(buff.Bytes())

    // Notify that rendering is complete
    serv.queueItem.complete <- true
}

func (serv webservice) configure(args WebRenderParams) (*lib.Info, error) {
    req := &serv.desc.UserRequest
    req.RealMin = format(args.RealMin)
    req.RealMax = format(args.RealMax)
    req.ImagMin = format(args.ImagMin)
    req.ImagMax = format(args.ImagMax)
    req.ImageWidth = args.ImageWidth
    req.ImageHeight = args.ImageHeight

    userDesc, cerr := lib.Configure(req)

    if cerr != nil {
        return nil, cerr
    }

    userDesc.NativeInfo = serv.desc.NativeInfo
    userDesc.UserRequest = *req

    return userDesc, nil
}

func format(f float64) string {
    return strconv.FormatFloat(f, 'e', -1, 64)
}