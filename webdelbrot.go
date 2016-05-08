package main

import (
    "log"

    "github.com/gopherjs/gopherjs/js"
)

func main() {
    godel := newgodel()
    js.Global.Set("godel", js.MakeWrapper(godel))
}

func newgodel() *Godel {
    const batchlength = 300
    godel := &Godel{}
    godel.debounce = newdebounce(batchlength)
    return godel
}

type Godel struct {
    debounce *debouncer
}

func (godel *Godel) Redraw() {
    gethistory().restart()
}

func (godel *Godel) Fractal_mousedown(event *js.Object) bool {
    x, y := mousepos(event)
    getfractal().zoom(x, y)
    return false
}

func (godel *Godel) Fractal_mousemove(event *js.Object) bool {
    x, y := mousepos(event)
    getfractal().inspect(x, y)
    return false
}

func (godel *Godel) Fractal_contextmenu(event *js.Object) bool {
    return false
}

func (godel *Godel) Fractal_mouseup(event *js.Object) bool {
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
    keys := []string{"x", "y"}
    vals := make([]uint, len(keys))
    for i, k := range keys {
        val64 := event.Get(k).Uint64()
        vals[i] = cropu64(val64)
    }
    return vals[0], vals[1]
}


func (godel *Godel) Window_resize() {
    godel.debounce.do(func () {
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