package libgodelbrot

type NativeMandelbrotThunk struct {
	BaseThunk
	NativeMandelbrotMember
}

func NewNativeMandelbrotThunkReals(r float64, i float64) *NativeMandelbrotThunk {
	return NewNativeMandelbrotThunk(complex(r, i))
}

func NewNativeMandelbrotThunk(c complex128) *NativeMandelbrotThunk {
	return &NativeMandelbrotThunk{
		NativeMandelbrotMember: NativeMandelbrotMember{
			C: c
		}
	}
}

type NativeMandelbrotThunkHeap struct {
	zone  []NativeMandelbrotThunk
	index int
}

func NewNativeMandelbrotThunkHeap(size uint) *NativeMandelbrotThunkHeap {
	return &NativeMandelbrotThunkHeap{
		zone:  make([]NativeMandelbrotThunk, 0, int(size)),
		index: 0,
	}
}

func (heap *NativeMandelbrotThunkHeap) NativeMandelbrotThunk(x float64, y float64) *NativeMandelbrotThunk {
	heap.zone = append(heap.zone, NativeMandelbrotThunk{evaluated: false, c: complex(x, y)})
	index := heap.index
	heap.index++
	return &heap.zone[index]
}
