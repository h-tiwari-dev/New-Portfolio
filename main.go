package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"app/models"
)

var router *gin.Engine

func main() {
	// IP RATE LIMITER
	rate, err := limiter.NewRateFromFormatted("240-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance_iprate := limiter.New(store, rate)
	middleware_iprate := mgin.NewMiddleware(instance_iprate)
	// END IP RATE LIMITER

	// DB MIGRATIONS
	models.Migrate()
	// END DB MIGRATIONS

	router = gin.Default()

	router.ForwardedByClientIP = true
	router.Use(middleware_iprate)

	// templ := template.Must(template.New("").ParseFS(embeddedFiles, "templates/*"))
	// router.SetHTMLTemplate(templ)
	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("templates/*")

	InitializeRoutes()

	router.Run()
}
