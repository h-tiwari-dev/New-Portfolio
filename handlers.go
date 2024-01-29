package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"app/helpers"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func GenerateUrl() string {
	var url string

	urlExist := true
	for urlExist {
		url = helpers.GenerateRandomString(7)
		//Verify is the url is unique
		//list := models.List{}
		//db.DB.Where("url = ?", url).First(&list)
		//if list.ID == 0 {
		//	urlExist = false
		//}
	}

	return url
}

func Home(title string, c *gin.Context) {
	render(c, gin.H{"title": title}, "home.html")
}

func Routes(c *gin.Context) {
	var result []gin.H
	for _, item := range navItems {
		result = append(result, gin.H{"name": item.Name, "url": item.URL})
	}
	render(c, gin.H{"navItems": result}, "navitems.html")
}

func DownloadResume(c *gin.Context) {
	DownloadFile(c, "Resume.pdf")
}

func DownloadDSA(c *gin.Context) {
	DownloadFile(c, "Data Structure And Algorithms specialization Certificate.pdf")
}

func DownloadDLS(c *gin.Context) {
	DownloadFile(c, "Deep Learning specialization Certificate.pdf")
}

func DownloadFile(c *gin.Context, fileName string) {
	filePath := filepath.Join("assets/data", fileName)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set the appropriate headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprint(fileSize(filePath)))

	// Open and read the file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer file.Close()

	// Send the file as the response
	c.FileAttachment(filePath, filePath)
}

// Helper function to get the file size
func fileSize(filePath string) int64 {
	file, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return file.Size()
}

type Blog struct {
	Heading string
	Content string
}

func LookingForWork(c *gin.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		print(err)
		return
	}
	LOOKING_FOR_WORK := os.Getenv("LOOKING_FOR_WORK")
	i, _ := strconv.Atoi(LOOKING_FOR_WORK)

	if i == 1 {
		render(c, gin.H{}, "lookingForWork.html")
	} else {
		return
	}
}

func Blogs(title string, c *gin.Context) {
	blogs := []Blog{
		{Heading: "Blog 1", Content: "Hello World"},
		{Heading: "Blog 2", Content: "HTMX and Golang"},
		{Heading: "Blog 2", Content: "HTMX and Golang"},
		{Heading: "Blog 2", Content: "HTMX and Golang"},
		{Heading: "Blog 2", Content: "HTMX and Golang"},
		{Heading: "Blog 2", Content: "HTMX and Golang"},
	}

	render(c, gin.H{"title": title, "blogs": blogs}, "blogs.html")
}
