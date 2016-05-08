package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/johnny-morrice/godelbrot/rest/protocol"
)

type rendercmd struct {
	parent string // parent is CAS address of prezoomed image
	renreq protocol.RenderRequest
}

func defaultrendercmd() *rendercmd {
	cmd := &rendercmd{}
	w, h := getfractal().dims()
	cmd.renreq.Req.ImageWidth = cropu64(w)
	cmd.renreq.Req.ImageHeight = cropu64(h)
	return cmd
}

func (cmd *rendercmd) render() (*img, error) {
	rest := getclient()
	datach := make(chan io.Reader)
	errch := make(chan error)
	go func() {
		data, err := rest.Cycle(cmd.parent, &cmd.renreq)
		errch <- err
		datach<- data
	}()
	if err := <-errch; err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(<-datach)
	if err != nil {
		return nil, err
	}

	pic := &img{}
	pic.data = bytes
	return pic, nil
}

type img struct {
	url string
	data []byte
}

func (pic *img) uri() string {
	if pic.url == "" {
		return fmt.Sprintf("data:image/png;base64,%v", string(pic.data))
	} else {
		return pic.url
	}
}

func cropu64(n uint64) uint {
	const maxuint = uint64(^uint(0))
	if n > maxuint {
		panic("Unsigned integer overflow")
	}
	return uint(n)
}