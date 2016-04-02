package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    "github.com/johnny-morrice/godelbrot/process"
    lib "github.com/johnny-morrice/godelbrot/libgodelbrot"
)

type RenderRequest struct {
    Req lib.Request
    Target lib.ZoomTarget
    WantZoom bool
}

func (rr *RenderRequest) validate() error {
    if rr.Req.ImageWidth < 1 || rr.Req.ImageHeight < 1 {
        return errors.New("Invalid Req")
    }

    validerr := rr.Target.Validate();
    if rr.WantZoom && validerr != nil {
        return errors.New("Invalid Target")
    }

    if !rr.WantZoom && validerr == nil {
        return errors.New("False WantZoom yet valid Target")
    }


    return nil
}

type RenderResponse struct {
    NextReq lib.Request
    ImageURL string
}

type PictureRequest struct {
    Code string
}

const httpheader string = "X-Godelbrot-Packet"
const formkey string = "godelbrotPacket"

type session struct {
    w http.ResponseWriter
    req *http.Request
}

func (s session) httpError(msg string, code int) error {
    http.Error(s.w, msg, code)
    return errors.New(fmt.Sprintf("(%v) %v", s.req.RemoteAddr, msg))
}

func (s session) internalError() error {
    return s.httpError("Internal error", 500)
}

type sem chan bool

func semaphor(concurrent uint) sem {
    return sem(make(chan bool, concurrent))
}

func (s sem) acquire(n uint) {
    for i := uint(0); i < n; i++ {
        s<- true
    }
}

func (s sem) release(n uint) {
    for i := uint(0); i < n; i++ {
        <-s
    }
}

type renderBuffers struct {
    png bytes.Buffer
    info bytes.Buffer
    nextinfo bytes.Buffer
    report bytes.Buffer
}

func (rb renderBuffers) logReport() {
    sc := bufio.NewScanner(&rb.report)
    for sc.Scan() {
        err := sc.Err()
        if err != nil {
            log.Printf("Error while printing error (omg!): %v", err)
        }
        log.Println(sc.Text())
    }
}

func (rb renderBuffers) input(info *lib.Info) error {
    return lib.WriteInfo(&rb.info, info)
}

// renderService renders fractals
type renderService struct {
    s sem
}

// makeRenderService creates a render service that allows at most `concurrent` concurrent tasks.
func makeRenderService(concurrent uint) renderService {
    rs := renderService{}
    rs.s = semaphor(concurrent)
    return rs
}

// render a fractal into the renderBuffers
func (rs renderService) render(rbuf renderBuffers, zoomArgs []string) error {
    rs.s.acquire(1)
    var err error
    if zoomArgs == nil || len(zoomArgs) == 0 {
        next, zerr := process.ZoomRender(&rbuf.info, &rbuf.png, &rbuf.report, zoomArgs)
        err = zerr
        if zerr != nil {
            _, err = io.Copy(&rbuf.nextinfo, next)
        }
    } else {
        tee := io.TeeReader(&rbuf.info, &rbuf.nextinfo)
        err = process.Render(tee, &rbuf.png, &rbuf.report)
    }
    rs.s.release(1)
    return err
}

type rendering struct {
    unixTime int64
    info *lib.Info
    png []byte
    code string
}

func (r rendering) hashcode() string {
    if r.code == "" {
        r.code = ""
    }
    return r.code
}

type webservice struct {
    baseinfo lib.Info
    rs renderService
    cache chan map[string]rendering
}

func makeWebservice(baseinfo *lib.Info, concurrent uint) *webservice {
    ws := &webservice{}
    ws.baseinfo = *baseinfo
    ws.rs = makeRenderService(concurrent)
    return ws
}

func (ws *webservice) pictureHandler(w http.ResponseWriter, req *http.Request) {
    sess := session{}
    sess.w = w
    sess.req = req
    err := ws.serveFractal(sess)
    if err != nil {
        log.Println(err)
    }
}

func (ws *webservice) renderHandler(w http.ResponseWriter, req *http.Request) {
    sess := session{}
    sess.w = w
    sess.req = req
    resp, err := ws.renderFractal(sess)
    if err != nil {
        log.Println(err)
        return
    }
    err = ws.serveInfo(sess, resp)
    if err != nil {
        log.Println(err)
    }
}

func (ws *webservice) renderFractal(s session) (*RenderResponse, error) {
    jsonPacket := s.req.FormValue(formkey)

    if len(jsonPacket) == 0 {
        err := s.httpError(fmt.Sprintf("No data found in parameter '%v'", formkey), 400)
        return nil, err
    }

    dec := json.NewDecoder(strings.NewReader(jsonPacket))
    renreq := &RenderRequest{}
    jsonerr := dec.Decode(renreq)

    if jsonerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid JSON"), 400)
        log.Println(err)
        return nil, jsonerr
    }

    validerr := renreq.validate()
    if validerr != nil {
        err := s.httpError(fmt.Sprintf("Invalid render request: %v", validerr), 400)
        log.Println(err)
        return nil, validerr
    }

    ws.safeTarget(renreq)
    info := ws.mergeInfo(renreq)

    buffs := renderBuffers{}
    bufferr := buffs.input(info)
    if bufferr != nil {
        err := s.internalError()
        log.Println(err)
        return nil, bufferr
    }

    var zoomArgs []string
    if renreq.WantZoom {
        zoomArgs = process.ZoomArgs(renreq.Target)
    }
    renderErr := ws.rs.render(buffs, zoomArgs)

    // Copy any stderr messages
    buffs.logReport()

    if renderErr != nil {
        err := s.internalError()
        log.Println(err)
        return nil, renderErr
    }

    nextinfo, infoerr := lib.ReadInfo(&buffs.nextinfo)
    if infoerr != nil {
        err := s.internalError()
        log.Println(err)
        return nil, infoerr
    }

    rndr := rendering{}
    rndr.png = buffs.png.Bytes()
    rndr.info = info
    rndr.unixTime = time.Now().Unix()
    code := rndr.hashcode()

    imgcache := <-ws.cache
    imgcache[code] = rndr
    ws.cache<- imgcache

    resp := &RenderResponse{}
    resp.ImageURL = fmt.Sprintf("/image/%v", code)
    resp.NextReq = nextinfo.UserRequest

    return resp, nil
}

func (ws *webservice) serveInfo(s session, resp *RenderResponse) error {
    s.w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(s.w)
    jsonErr := enc.Encode(resp)
    if jsonErr != nil {
        err := s.internalError()
        log.Println(err)
        return jsonErr
    }
    return nil
}

func (ws *webservice) serveFractal(s session) error {
    formval := s.req.FormValue(formkey)
    picreq := &PictureRequest{}
    dec := json.NewDecoder(strings.NewReader(formval))
    jsonerr := dec.Decode(picreq)

    if jsonerr != nil {
        err := s.internalError()
        log.Println(err)
        return jsonerr
    }

    if picreq.Code == "" {
        err := s.httpError("No Code", 400)
        return err
    }

    imgcache := <-ws.cache
    ws.cache<- imgcache
    rndr, ok := imgcache[picreq.Code]

    if !ok {
        err := s.httpError(fmt.Sprintf("Invalid Code: %v", picreq.Code), 400)
        return err
    }

    buff := bytes.NewBuffer(rndr.png)
    // Write image buffer as http response
    s.w.Header().Set("Content-Type", "image/png")
    _, cpyerr := io.Copy(s.w, buff)
    if cpyerr != nil {
        err := s.internalError()
        log.Println(err)
        return cpyerr
    }

    return nil
}

// Only allow zoom reconfiguration if autodetection is enabled throughout the base info.
func (ws *webservice) safeTarget(renreq *RenderRequest) {
    req := ws.baseinfo.UserRequest
    dyn := req.Renderer == lib.AutoDetectRenderMode
    dyn = dyn && req.Numerics == lib.AutoDetectNumericsMode

    renreq.Target.UpPrec = dyn
    renreq.Target.Reconfigure = dyn
    renreq.Target.Frames = 1
}

func (ws *webservice) mergeInfo(renreq *RenderRequest) *lib.Info {
    req := ws.baseinfo.UserRequest

    req.ImageWidth = renreq.Req.ImageWidth
    req.ImageHeight = renreq.Req.ImageHeight

    inf := new(lib.Info)
    *inf = ws.baseinfo
    inf.UserRequest = req

    return inf
}

func format(f float64) string {
    return strconv.FormatFloat(f, 'e', -1, 64)
}