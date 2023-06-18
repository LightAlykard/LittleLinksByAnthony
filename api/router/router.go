package router

import "github.com/LightAlykard/LittleLinksByAnthony/api/handler"

type Router struct {
	hs *handler.Handlers
}

func NewRouter(hs *handler.Handlers) *Router {
	return &Router{
		hs: hs,
	}
}
