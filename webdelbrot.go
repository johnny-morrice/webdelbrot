package main

import (
    "log"

    "github.com/gopherjs/gopherjs/js"
)

func main() {
    godel := getgodel()
    js.Global.Set("godel", js.MakeWrapper(godel))
}

var __godel *Godel
func getgodel() *Godel {
    if __godel == nil {
        __godel = &Godel{}
        __godel.debounce = map[string]*debouncer {}
        __godel.debounce["resize"] = newdebounce(__RESIZE_MS)
    }

    return __godel
}

type Godel struct {
    debounce map[string]*debouncer
}

func (godel *Godel) Redraw() {
    gethistory().restart()
}

func (godel *Godel) Fractal_mousedown(event *js.Object) bool {
    stopPropagation(event)
    x, y := mousepos(event)
    getfractal().zoom(x, y)    

    return false
}

func (godel *Godel) Fractal_touchstart(event *js.Object) bool {
    stopPropagation(event)
    x, y := touchpos(event)
    getfractal().zoom(x, y)    

    return false
}

func (godel *Godel) Fractal_mousemove(event *js.Object) bool {
    stopPropagation(event)
    x, y := mousepos(event)
    getfractal().inspect(x, y)
    return false
}

func (godel *Godel) Fractal_touchend(event *js.Object) bool {
    stopPropagation(event)
    getfractal().mark()
    return false
}

func (godel *Godel) Fractal_contextmenu(event *js.Object) bool {
    stopPropagation(event)
    getfractal().cancel()
    return false
}

func (godel *Godel) Fractal_mouseup(event *js.Object) bool {
    stopPropagation(event)
    switch button := event.Get("button").Uint64(); button {
    case 0:
        getfractal().mark()
    case 2:
        getfractal().cancel()        
    default:
        if __DEBUG {
            log.Printf("Not handling mouse button %v", button)    
        }
    }
    return false
}

func mousepos(event *js.Object) (uint, uint) {
    keys := []string{"clientX", "clientY"}
    vals := make([]uint, len(keys))
    for i, k := range keys {
        val64 := event.Get(k).Uint64()
        vals[i] = cropu64(val64)
    }
    return vals[0], vals[1]
}

func touchpos(event *js.Object) (uint, uint) {
    touches := event.Get("touches")

    return mousepos(touches.Index(0))
}


func (godel *Godel) Window_resize() {
    godel.debounce["resize"].do(func () {
        gethistory().render()
    })
}

func (godel *Godel) Toolbar_restart_click() bool {
    gethistory().restart()
    return false
}

func (godel *Godel) Toolbar_back_click() bool {
    gethistory().back()
    return false
}

func stopPropagation(event *js.Object) {
    event.Call("stopPropagation")
    event.Call("preventDefault")
}