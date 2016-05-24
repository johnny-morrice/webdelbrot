package main

import (
	"log"
)

type history struct {
	frames []*rendercmd
}

var __history *history

func gethistory() *history {
    if __history == nil {
    	__history = &history{}
    	__history.frames = []*rendercmd{getfractal().defaultrendercmd()}
    }

    return __history
}

func (h *history) back() {
	count := len(h.frames)
	if count > 1 {
		h.frames = h.frames[:count - 1]
		h.render()
	}
}

func (h *history) last() *rendercmd {
	return h.frames[len(h.frames) - 1]
}

func (h *history) restart() {
	h.frames = h.frames[:1]
	h.render()
}

func (h *history) render() {
	go func() {
		pic, err := h.last().render()
		if err != nil {
			log.Printf("Render error: %v", err)
			return
		}
		getfractal().replace(pic)
	}()
}

func (h *history) zoom(bounds []uint) {

}