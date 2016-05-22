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
		log.Println("Begin zoom")
		fr.startzoom(x, y)
		fr.state = PROGRESS
	case PROGRESS:
		fr.zoomin()
	default:
		panic(fmt.Sprintf("Unknown zoom state: %v", fr.state))
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
		log.Println("Mark zoom")
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
		panic(fmt.Sprintf("Invalid zoom state: %v", fr.state))
	case PROGRESS:
		log.Println("Cancel zoom")
		fr.cleanzoom()
		fr.state = BEGIN
	}
}

func (fr *fractal) markzoom() {
	// Nothing yet
}

func (fr *fractal) zoomin() {
	elapsed := time.Now().Sub(fr.begin) * time.Millisecond
	log.Printf("Milliseconds since zoom start: %v", elapsed)

	exp := float64(elapsed) / 1000.0
	shrink := math.Pow(__SHRINK_RATE, exp)

	w, h := fr.dims()

	fw, fh := float64(w), float64(h)
	fxmouse, fymouse := float64(fr.xmouse), float64(fr.ymouse)

	zwsect := (fw * shrink) / 2.0
	zhsect := (fh * shrink) / 2.0

	botoffset := float64(fr.toolbarheight())

	fbounds := []float64{
		fxmouse - zwsect,
		fw - (fxmouse + zwsect),
		fymouse - zhsect,
		(fh - (fymouse + zhsect)) + botoffset,
	}

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

	for i, p := range props {
		style := fr.zoombox.Get("style")
		style.Set(p, bounds[i])
	}
}

func (fr *fractal) cleanzoom() {
	fr.zoombox.Get("parentNode").Call("removeChild", fr.zoombox)
}

func (fr *fractal) startzoom(x, y uint) {
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