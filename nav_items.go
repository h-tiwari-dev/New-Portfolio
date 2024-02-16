package main

import (
	"github.com/gin-gonic/gin"

	"app/handlers"
)

type NavItem struct {
	Name     string
	URL      string
	Renderer func(string, *gin.Context)
}

var navItems = []NavItem{
	{Name: "/home", URL: "/", Renderer: handlers.Home},
	{Name: "/blogs", URL: "/blogs", Renderer: handlers.Blogs},
}
