package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

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

func (cmd *rendercmd) render() (*img, error) {
	rest := getclient()
	png, err := rest.Cycle(cmd.parent, &cmd.renreq)
	if err != nil {
		return nil, err
	}

	b64 := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, b64)
	_, encerr := io.Copy(enc, png)
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