package nativebase

import (
    "functorama.com/demo/base"
)

type MockRenderApplication struct {
    base.MockRenderApplication
    MockNativeCoordProvider
}

type MockNativeCoordProvider struct {
    TNativeUserCoords bool

    PlaneMin complex128
    PlaneMax complex128
    FixAspect bool
}

func (mock *MockNativeCoordProvider) NativeUserCoords() (complex128, complex128) {
    mock.TNativeUserCoords = true
    return mock.PlaneMin, mock.PlaneMax
}