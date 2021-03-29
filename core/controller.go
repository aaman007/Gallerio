package core

import "go-web-dev-2/views"

func NewController() *Controller {
	return &Controller{
		HomeView: views.NewView("base", "core/home"),
		ContactView: views.NewView("base", "core/contact"),
	}
}

type Controller struct {
	HomeView *views.View
	ContactView *views.View
}