package main

import (
	"fmt"

    "github.com/gopherjs/gopherjs/js"
)

type fractal struct {
	fractal *js.Object
	toolbar *js.Object
	window *js.Object
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

}

// Finish zoom
func (fr *fractal) mark() {

}

func (fr *fractal) inspect(x, y uint) {

}

func (fr *fractal) cancel() {

}

func (fr *fractal) dims() (uint, uint) {
	w := fr.window.Get("innerWidth").Uint64()
	wh := fr.window.Get("innerHeight").Uint64()
	toolh := fr.toolbar.Call("getAttribute", "height").Uint64()
	h := wh - toolh
	return cropu64(w), cropu64(h)
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