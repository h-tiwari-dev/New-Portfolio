package main

import (
	"github.com/gin-gonic/gin"

	"app/handlers"
)

type navItem struct {
	Name     string
	URL      string
	Renderer func(string, *gin.Context)
}
type Nav struct {
	handlers *handlers.Handlers
	items    []navItem
	router   *gin.Engine
}

func NewNav(_handlers *handlers.Handlers) *Nav {
	return &Nav{
		handlers: _handlers,
	}
}

func (nv *Nav) getNavItems() []navItem {
	return []navItem{
		{Name: "/home", URL: "/", Renderer: nv.handlers.Home},
		{Name: "/blogs", URL: "/blogs-page", Renderer: nv.handlers.BlogsPage},
	}
}

func (nv *Nav) registerRoutes() {
	for _, item := range nv.getNavItems() {
		localItem := item

		router.GET(localItem.URL, func(ctx *gin.Context) {
			localItem.Renderer(localItem.Name, ctx)
		})
	}
}
