package main

import (
    "github.com/gopherjs/gopherjs/js"
)

type fractal struct {
	element *js.Object
}

var __fractal *fractal

func getfractal() *fractal {
	const elementID = "fractal"
    if __fractal == nil {
    	__fractal = &fractal{}
    	__fractal.element = js.Global.Get("document").Call("getElementById", elementID);
    }

    return __fractal
}

func (fr *fractal) zoom(x, y uint) {

}

func (fr *fractal) inspect(x, y uint) {

}

func (fr *fractal) cancel() {

}

func (fr *fractal) dims() (uint64, uint64) {
	w := fr.element.Call("getAttribute", "width")
	h := fr.element.Call("getAttribute", "height")
	return w.Uint64(), h.Uint64()
}

func (fr *fractal) replace(pic *img) {
	fr.element.Call("setAttribute", pic.uri())
}