package main

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes() {
	for _, item := range navItems {
		localItem := item

		router.GET(localItem.URL, func(ctx *gin.Context) {
			localItem.Renderer(localItem.Name, ctx)
		})
	}

	router.GET("/routes", Routes)
	router.GET("/download/resume", DownloadResume)
	router.GET("/download/DSA", DownloadDSA)
	router.GET("/download/DLS", DownloadDLS)
	router.GET("/checkForWork", LookingForWork)
}
