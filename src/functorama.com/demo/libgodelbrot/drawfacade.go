package libgodelbrot

import (
    "image"
    "functorama.com/demo/draw"
)

type drawFacade struct {
    picture *image.NRGBA
    colors draw.Palette
}

var _ draw.DrawingContext = (*drawFacade)(nil)
var _ draw.ContextProvider = (*drawFacade)(nil)

func (facade *drawFacade) DrawingContext() draw.DrawingContext {
    return facade
}

func (facade *drawFacade) Colors() draw.Palette {
    return facade.colors
}

func (facade *drawFacade) Picture() *image.NRGBA {
    return facade.picture
}

func makeDrawFacade(desc *Info) *drawFacade {
    facade := &drawFacade{}
    facade.colors = createPalette(desc)
    facade.picture = createImage(desc)
    return facade
}

func createImage(desc *Info) *image.NRGBA {
    req := desc.UserRequest
    bounds := image.Rectangle{
        Min: image.ZP,
        Max: image.Point{
            X: int(req.ImageWidth),
            Y: int(req.ImageHeight),
        },
    }
    return image.NewNRGBA(bounds)
}

func createPalette(desc *Info) draw.Palette {
    return draw.NewGrayscalePalette(desc.UserRequest.IterateLimit)
}
