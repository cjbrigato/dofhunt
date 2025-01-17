package ui

import (
	"image/color"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/settings"
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

func WithUIStyle(fn func(), settings *settings.AppSettings) {
	imgui.PushStyleVarVec2(imgui.StyleVarCellPadding, imgui.Vec2{1.0, 1.0})
	imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{1.0, 1.0})
	imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 0.0})
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 6.0)
	imgui.PushStyleVarFloat(imgui.StyleVarChildBorderSize, 0)
	imgui.PushStyleColorVec4(imgui.ColChildBg, g.ToVec4Color(color.RGBA{50, 50, 70, 0}))
	imgui.PushStyleColorVec4(imgui.ColButton, g.ToVec4Color(color.RGBA{50, 50, 70, 130}))
	g.PushColorWindowBg(settings.WindowColor)
	g.PushColorFrameBg(settings.FrameColor)

	fn()

	g.PopStyleColorV(4)
	imgui.PopStyleVarV(6)
}
