package main

import (
	"github.com/gin-gonic/gin"
)

type NavItem struct {
	Name     string
	URL      string
	Renderer func(string, *gin.Context)
}

var navItems = []NavItem{
	{Name: "/home", URL: "/", Renderer: Home},
	// {Name: "/blogs", URL: "/blogs", Renderer: Blogs},
}
