package dpi

import (
	"syscall"
)

const (
	DPIAwarenessContextUnaware           = 16
	DPIAwarenessContextSystemAware       = 17
	DPIAwarenessContextPerMonitorAware   = 18
	DPIAwarenessContextPerMonitorAwareV2 = 34
	DPIUnawarePixels                     = 96
)

var (
	scalingFactor                  = 1.0
	user32                         = syscall.NewLazyDLL("user32.dll")
	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	pGetActiveWindow               = user32.NewProc("GetActiveWindow")
	pGetDpiForWindow               = user32.NewProc("GetDpiForWindow")
)

func SetDpiAware() {
	pSetProcessDpiAwarenessContext.Call(DPIAwarenessContextPerMonitorAwareV2)
}

func InitScalingFactor() {
	hwnd, _, _ := pGetActiveWindow.Call()
	r, _, _ := pGetDpiForWindow.Call(hwnd)
	scalingFactor = float64(r) / float64(DPIUnawarePixels)
}

func Scaledi(i int) int {
	return int(float64(i) * scalingFactor)
}

func Scaledf32(f float32) float32 {
	return float32(float64(f) * scalingFactor)
}

func Scaledf(f float64) float64 {
	return f * scalingFactor
}
