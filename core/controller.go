package core

import "go-web-dev-2/views"

func NewStaticController() *StaticController {
	return &StaticController{
		HomeView: views.NewView("base", "core/home"),
		ContactView: views.NewView("base", "core/contact"),
	}
}

type StaticController struct {
	HomeView *views.View
	ContactView *views.View
}