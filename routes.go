package main

import (
	"github.com/gin-gonic/gin"

	"app/handlers"
	"app/helpers"
	middleware "app/middlewares"
)

func routes(nav *Nav, c *gin.Context) {
	var result []gin.H
	for _, item := range nav.getNavItems() {
		result = append(result, gin.H{"name": item.Name, "url": item.URL})
	}
	helpers.Render(c, gin.H{"navItems": result}, "navitems.html")
}

func InitializeRoutes() {
	handlers := handlers.NewHandlers()
	nav := NewNav(handlers)

	nav.registerRoutes()
	router.GET("/routes", func(ctx *gin.Context) {
		routes(nav, ctx)
	})
	router.GET("/download/resume", handlers.DownloadResume)
	router.GET("/download/DSA", handlers.DownloadDSA)
	router.GET("/download/DLS", handlers.DownloadDLS)
	router.GET("/checkForWork", handlers.LookingForWork)

	authMiddleware := middleware.Authenticator(false)

	// Public routes
	router.POST("/signup", handlers.SignupHandler)
	router.POST("/login", handlers.LoginHandler)

	// Authenticated routes
	authGroup := router.Group("/")
	authGroup.Use(authMiddleware)
	{
		// Your existing authenticated route(s)
		// authGroup.GET("/profile", profileHandler)

		// Protected routes for blogs accessible only to authenticated users
		blogGroup := authGroup.Group("/blogs")
		{
			blogGroup.POST("/", handlers.CreateBlog)
			blogGroup.GET("/", handlers.GetBlogs) // New route for getting all blogs
			blogGroup.GET("/:id", handlers.GetBlogById)
			blogGroup.PUT("/:id", handlers.UpdateBlog)
			blogGroup.DELETE("/:id", handlers.DeleteBlog)
		}
	}
}
