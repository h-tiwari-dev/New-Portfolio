package main

import (
	"github.com/gin-gonic/gin"

	"app/handlers"
	"app/helpers"
)

func routes(c *gin.Context) {
	var result []gin.H
	for _, item := range navItems {
		result = append(result, gin.H{"name": item.Name, "url": item.URL})
	}
	helpers.Render(c, gin.H{"navItems": result}, "navitems.html")
}

func InitializeRoutes() {
	for _, item := range navItems {
		localItem := item

		router.GET(localItem.URL, func(ctx *gin.Context) {
			localItem.Renderer(localItem.Name, ctx)
		})
	}

	router.GET("/routes", routes)
	router.GET("/download/resume", handlers.DownloadResume)
	router.GET("/download/DSA", handlers.DownloadDSA)
	router.GET("/download/DLS", handlers.DownloadDLS)
	router.GET("/checkForWork", handlers.LookingForWork)

	authMiddleware := middleware.Authenticator()

	// Public routes
	router.POST("/signup", signupHandler)
	router.POST("/login", loginHandler)

	// Authenticated routes
	authGroup := router.Group("/auth")
	authGroup.Use(authMiddleware)
	{
		// Your existing authenticated route(s)
		authGroup.GET("/profile", profileHandler)

		// Protected routes for blogs accessible only to authenticated users
		blogGroup := authGroup.Group("/blogs")
		{
			blogGroup.POST("/", createBlogHandler)
			blogGroup.GET("/", getAllBlogsHandler) // New route for getting all blogs
			blogGroup.GET("/:id", getBlogHandler)
			blogGroup.PUT("/:id", updateBlogHandler)
			blogGroup.DELETE("/:id", deleteBlogHandler)
		}
	}
}
