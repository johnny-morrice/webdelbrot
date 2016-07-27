package main

import (
    "github.com/gopherjs/gopherjs/js"
)

type hud struct {
	loadelem *js.Object
}

var __hud *hud

func gethud() *hud {
	if __hud == nil {
		__hud = &hud{}
	}

	return __hud
}

func (h *hud) loading(bounds []uint) {
	document := js.Global.Get("document")
	loadelem := document.Call("createElement", "div")

	loadelem.Call("setAttribute", "id", "loadbox")

	body := document.Get("body")

	body.Call("appendChild", loadelem)

	pix := pixbounds(bounds)
	setbounds(loadelem, pix)

	h.loadelem = loadelem
}

func (h *hud) noloading() {
	if (h.loadelem == nil) {
		return
	}

	parent := h.loadelem.Get("parentElement")

	parent.Call("removeChild", h.loadelem)

	h.loadelem = nil
}