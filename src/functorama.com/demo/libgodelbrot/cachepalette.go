package libgodelbrot

import (
	"image/color"
)

type Cacher func(iterateLimit uint8, index uint8) color.NRGBA

type CachePalette struct {
	memberColor color.NRGBA
	scale       []color.NRGBA
	limit       uint8
}

func NewCachePalette(iterateLimit uint8, member color.NRGBA, cacher Cacher) CachePalette {
	colors := make([]color.NRGBA, iterateLimit, iterateLimit)
	iLimit := int(iterateLimit)
	for i := 0; i < iLimit; i++ {
		colors[i] = cacher(iterateLimit, uint8(i))
	}
	return CachePalette{
		memberColor: member,
		scale:       colors,
		limit:       iterateLimit,
	}
}

// CachePalette implements Palette
func (palette CachePalette) Color(member MandelbrotMember) color.NRGBA {
	if member.SetMember() {
		return palette.memberColor
	} else {
		return palette.scale[member.InvDivergence]
	}
}
