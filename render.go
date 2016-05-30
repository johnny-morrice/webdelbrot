package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/johnny-morrice/godelbrot/config"
	"github.com/johnny-morrice/godelbrot/rest/protocol"
)

type rendercmd struct {
	parent string // parent is CAS address of prezoomed image
	addr string
	renreq protocol.RenderRequest
}

func newcmd(w, h uint) *rendercmd {
	cmd := &rendercmd{}
	req := &cmd.renreq.Req
	req.ImageWidth, req.ImageHeight = w, h
	return cmd
}

func zoomcmd(w, h uint, parent string, bounds config.ZoomBounds) *rendercmd {
	cmd := newcmd(w, h)

	cmd.renreq.WantZoom = true
	cmd.renreq.Target = bounds
	cmd.parent = parent

	return cmd
}

func (cmd *rendercmd) render() (*img, error) {
	rest := getclient()
	result, err := rest.RenderCycle(cmd.parent, &cmd.renreq)
	if err != nil {
		return nil, err
	}
	cmd.addr = rest.Url(result.Status.ThisUrl)

	b64 := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, b64)
	_, encerr := io.Copy(enc, result.Image)
	if encerr != nil {
		return nil, encerr
	}

	pic := &img{}
	pic.data = b64.Bytes()
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

func (pic *img) cssurl() string {
	if pic.url == "" {
		return fmt.Sprintf("url(\"%v\")", pic.uri())
	} else {
		return pic.url
	}
}