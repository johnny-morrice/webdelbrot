package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

    "github.com/gopherjs/gopherjs/js"
)

type zoomstate uint8

const (
	BEGIN = zoomstate(iota)
	PROGRESS
)

type fractal struct {
	fractal *js.Object
	toolbar *js.Object
	window *js.Object

	state zoomstate
	xmouse uint
	ymouse uint
	begin time.Time
	zoombox *js.Object
	tick *time.Ticker

	lock sync.Mutex
}

var __fractal *fractal

func getfractal() *fractal {
	const fractalId = "fractal"
	const toolbarId = "toolbar"
    if __fractal == nil {
    	__fractal = &fractal{}
    	__fractal.fractal = getelementbyid(fractalId)
    	__fractal.toolbar = getelementbyid(toolbarId)
    	__fractal.window = js.Global.Get("window")
    }

    return __fractal
}

func (fr *fractal) zoom(x, y uint) {
	fr.lock.Lock()
	defer fr.lock.Unlock()
	switch fr.state {
	case BEGIN:
		fr.startzoom(x, y)
		fr.state = PROGRESS
	default:
		panic(fmt.Sprintf("Invalid zoom state: %v", fr.state))
	}
}

// Finish zoom
func (fr *fractal) mark() {
	fr.lock.Lock()
	defer fr.lock.Unlock()
	switch fr.state {
	case BEGIN:
		panic(fmt.Sprintf("Invalid zoom state: %v", fr.state))
	case PROGRESS:
		fr.markzoom()
		fr.state = BEGIN
	default:
		panic(fmt.Sprintf("Unknown zoom state: %v", fr.state))
	}
}

func (fr *fractal) cancel() {
	fr.lock.Lock()
	defer fr.lock.Unlock()
	switch fr.state {
	case BEGIN:
		log.Printf("Warning: invalid zoom state: %v", fr.state)
	case PROGRESS:
		log.Println("Cancel zoom")
		fr.cleanzoom()
		fr.state = BEGIN
	}
}

func (fr *fractal) markzoom() {
	log.Println("Mark zoom")
	
	// Nothing yet

	fr.cleanzoom()
}

func (fr *fractal) zoomin() {
	fr.lock.Lock()
	defer fr.lock.Unlock()

	if fr.state != PROGRESS {
		return
	}

	elapsed := time.Now().Sub(fr.begin)

	exp := float64(elapsed) / 1000000000.0
	shrink := math.Pow(__SHRINK_RATE, exp)

	w, h := fr.dims()

	fw, fh := float64(w), float64(h)
	fxmouse, fymouse := float64(fr.xmouse), float64(fr.ymouse)

	botoffset := float64(fr.toolbarheight())

	xbnd := shrink * math.Max(fxmouse, fw - fxmouse)
	ybnd := shrink * math.Max(fymouse, fh - fymouse)

	fbounds := []float64{
		fxmouse - xbnd, 
		fw - fxmouse - xbnd, 
		fymouse - ybnd, 
		botoffset + fh - fymouse - ybnd,
	}

	fbounds = preserveAspect(fxmouse, fymouse, fw, fh, fbounds)

	bounds := make([]string, len(fbounds))

	for i, fb := range fbounds {
		if fb < 0.0 {
			fb = 0.0
		}
		ib := uint(math.Floor(fb))
		bounds[i] = fmt.Sprintf("%vpx", ib)
	}

	props := []string{
		"left",
		"right",
		"top",
		"bottom",
	}

	if __DEBUG {
		log.Printf("Dims are: %v %v", fw, fh)
		log.Printf("Zoom time: %v", elapsed)	
		log.Printf("Exp is: %v", exp)
		log.Printf("Shrink factor is %v", shrink)
		log.Printf("Mouse at %v %v", fxmouse, fymouse)

		genfbounds := make([]interface{}, len(fbounds))
		for i, b := range fbounds {
			genfbounds[i] = b
		}

		genbounds := make([]interface{}, len(bounds))
		for i, b := range bounds {
			genbounds[i] = b
		}

		log.Printf("Fbounds are %v %v %v %v", genfbounds...)
		log.Printf("Bounds are %v %v %v %v", genbounds...)
	}


	for i, p := range props {
		style := fr.zoombox.Get("style")
		style.Set(p, bounds[i])
	}
}

func preserveAspect(x, y, w, h float64, bounds []float64) []float64 {
	xmin := bounds[0]
	xmax := bounds[1]
	ymin := bounds[2]
	ymax := bounds[3]
	bw := xmax - xmin
	bh := ymax - ymin

	aspect := w / h
	baspect := bw / bh

	resize := make([]float64, len(bounds))
	for i, b := range bounds {
		resize[i] = b
	}

	// Midpoints
	// x = (xmax + xmin) / 2
	// y = (ymax + ymin) / 2

	// Rule: always shrink
	if aspect < baspect {
		// Too tall, make thinner
		thinpart := (bh * aspect) / 2
		resize[0] = x - thinpart
		resize[1] = x + thinpart
		if __DEBUG {
			log.Printf("Made thinner; %v to %v", bw, resize[1] - resize[0])	
		}
	} else if aspect > baspect {
		// Too thin; make shorter
		shortpart := (bw / aspect) / 2
		resize[2] = y - shortpart
		resize[3] = y + shortpart
		if __DEBUG {
			log.Printf("Made shorter; %v to %v", bh, resize[3] - resize[2])	
		}
	}

	if __DEBUG {
		log.Printf("Screen aspect is %v", aspect)
		log.Printf("Preadjusted zoom aspect was %v", baspect)		
		log.Printf("Adjusted zoom aspect is %v", (resize[1] - resize[0]) / (resize[3] - resize[2]))
	}

	return resize
}

// Restart zoom tick
func (fr *fractal) zoomproc() {

	if fr.tick != nil {
		fr.tick.Stop()
	}

	tick := time.NewTicker(__ZOOM_MS * time.Millisecond)
	fr.tick = tick

	go func() {
		for range tick.C {
			fr.zoomin()
		}
	}()
}

func (fr *fractal) cleanzoom() {
	if __DEBUG {
		log.Println("Removing zoombox")
	}
	fr.zoombox.Get("parentNode").Call("removeChild", fr.zoombox)
}

func (fr *fractal) startzoom(x, y uint) {
	if __DEBUG {
		log.Println("Begin zoom")
	}
	fr.xmouse = x
	fr.ymouse = y

	document := js.Global.Get("document")

	fr.zoombox = document.Call("createElement", "div")
	fr.zoombox.Set("id", "zoombox")

	// Attach the global event listeners to the zoombox
	events := []string{
		"mousemove",
		"mousedown",
		"mouseup",
		"contextmenu",
	}

	handlers := []func(*js.Object) bool {
		getgodel().Fractal_mousemove,
		getgodel().Fractal_mousedown,
		getgodel().Fractal_mouseup,
		getgodel().Fractal_contextmenu,
	}

	for i, e := range events {
		h := handlers[i]
		fr.zoombox.Call("addEventListener", e, h)
	}

	fr.begin = time.Now()

	document.Get("body").Call("appendChild", fr.zoombox)

	fr.zoomproc()
}

// React to mouse motion
func (fr *fractal) inspect(x, y uint) {
	// NOP
}

func (fr *fractal) dims() (uint, uint) {
	w := fr.window.Get("innerWidth").Uint64()
	wh := fr.window.Get("innerHeight").Uint64()
	toolh := fr.toolbarheight()
	h := wh - toolh
	return cropu64(w), cropu64(h)
}

func (fr *fractal) toolbarheight() uint64 {
	return fr.toolbar.Call("getAttribute", "height").Uint64()
}

func (fr *fractal) defaultrendercmd() *rendercmd {
	cmd := &rendercmd{}
	req := &cmd.renreq.Req
	req.ImageWidth, req.ImageHeight = fr.dims()
	return cmd
}

func (fr *fractal) replace(pic *img) {
	w, h := fr.dims()
	fr.fractal.Call("setAttribute", "src", pic.uri())
	style := fr.fractal.Get("style")
	style.Set("width", fmt.Sprintf("%vpx", w))
	style.Set("height", fmt.Sprintf("%vpx", h))
}

func getelementbyid(id string) *js.Object {
	return js.Global.Get("document").Call("getElementById", id)
}

func cropu64(n uint64) uint {
	const maxuint = uint64(^uint(0))
	if n > maxuint {
		panic("Unsigned integer overflow")
	}
	return uint(n)
}