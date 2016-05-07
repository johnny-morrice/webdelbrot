package main

import (
    "fmt"
    "io"

    "github/johnny-morrice/godelbrot/config"
    "github/johnny-morrice/godelbrot/restclient"

    "github.com/gopherjs/gopherjs/js"
)

func main() {
    bindings := map[string]interface{}
    bindings["window_resize"] = js_window_resize
    bindings["fractal_click_zoom"] = js_fractal_click_zoom
    bindings["fractal_click_cancel"] = js_fractal_click_cancel
    bindings["fractal_mousemove"] = js_fractal_mousemove
    bindings["toolbar_restart_click"] = js_toolbar_restart_click
    bindings["toolbar_back_click"] = js_toolbar_back_click

    js.Global.Set("godel", bindings)
}

func js_fractal_click_zoom(event *js.Object) bool {
    x, y := mousepos(event)
    getfractal().zoom(x, y)
    return false
}

func js_fractal_mousemove(event *js.Object) bool {
    x, y := mousepos(event)
    getfractal().inspect(x, y)
    return false
}

func js_fractal_click_cancel(event *js.Object) bool {
    getfractal().cancel()
    return false
}

func mousepos(event *js.Object) (uint, uint) {
    keys := []string{"x", "y"}
    vals := make([]int, len(keys))
    for i, k := range keys {
        vals[i] = event.Key(k).Uint()
    }
    return vals[0], vals[1]
}

func js_window_resize() {
    gethistory().last().resize()
}

func js_toolbar_restart_click() bool {
    gethistory().restart()
    return false
}

func js_toolbar_back_click() bool {
    gethistory().back()
    return false
}