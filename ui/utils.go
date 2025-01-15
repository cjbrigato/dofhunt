package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

func FramelessWindowMoveWidget(widget g.Widget, isMovingFrame *bool, wnd *g.MasterWindow) *g.CustomWidget {
	return g.Custom(func() {
		if *isMovingFrame && !g.IsMouseDown(g.MouseButtonLeft) {
			*isMovingFrame = false
			return
		}

		widget.Build()

		if g.IsItemHovered() {
			if g.IsMouseDown(g.MouseButtonLeft) {
				*isMovingFrame = true
			}
		}

		if *isMovingFrame {
			delta := imgui.CurrentIO().MouseDelta()
			dx := int(delta.X)
			dy := int(delta.Y)
			if dx != 0 || dy != 0 {
				ox, oy := wnd.GetPos()
				wnd.SetPos(ox+dx, oy+dy)
			}
		}
	})
}
