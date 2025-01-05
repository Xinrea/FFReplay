package ui

import (
	"image"
	"sync"

	"github.com/Xinrea/ffreplay/internal/entry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/furex/v2"
)

// Due to the design of furex.View, root view cannot handle any events,
// so we need a global root as wrap.
var root *furex.View = &furex.View{ID: "root", Handler: furex.NewHandler(furex.HandlerOpts{
	Update: func(v *furex.View) {
		realroot := v.NthChild(0)
		if realroot != nil {
			if realroot.Height != v.Height || realroot.Width != v.Width {
				realroot.SetHeight(v.Height)
				realroot.SetWidth(v.Width)
			}
		}
	},
})}

type PlaygroundUI struct {
	view *furex.View
	once sync.Once
}

var _ UI = (*PlaygroundUI)(nil)

func NewPlaygroundUI(ecs *ecs.ECS) *PlaygroundUI {
	ecsInstance = ecs
	handler := &FocusHandler{}
	view := &furex.View{
		ID:        "Playground",
		Direction: furex.Column,
		Justify:   furex.JustifySpaceBetween,
	}
	handler.view = view
	view.Handler = handler
	root.AddChild(view)
	return &PlaygroundUI{
		view: view,
	}
}

type FocusHandler struct {
	view *furex.View
}

// HandleJustPressedMouseButtonLeft implements furex.MouseLeftButtonHandler.
func (f *FocusHandler) HandleJustPressedMouseButtonLeft(frame image.Rectangle, x int, y int) bool {
	for _, c := range f.view.FilterByTagName("input") {
		if fh, ok := c.Handler.(Focusable); ok {
			fh.SetFocus(false)
		}
	}
	entry.GetGlobal(ecsInstance).UIFocus = false
	return false
}

// HandleJustReleasedMouseButtonLeft implements furex.MouseLeftButtonHandler.
func (f *FocusHandler) HandleJustReleasedMouseButtonLeft(frame image.Rectangle, x int, y int) {
}

var _ furex.MouseLeftButtonHandler = (*FocusHandler)(nil)

func (p *PlaygroundUI) Update(w, h int) {
	global := entry.GetGlobal(ecsInstance)
	if global.Loaded.Load() {
		p.once.Do(func() {
			command := CommandView()
			command.MarginBottom = 20
			command.MarginLeft = 20
			partyList := NewPartyList(nil)
			partyList.MarginTop = 40
			partyList.MarginLeft = 20
			p.view.AddChild(partyList)
			p.view.AddChild(command)
		})
	}
	s := ebiten.Monitor().DeviceScaleFactor()
	furex.GlobalScale = s
	root.UpdateWithSize(w, h)
}

func (p *PlaygroundUI) Draw(screen *ebiten.Image) {
	root.Draw(screen)
}