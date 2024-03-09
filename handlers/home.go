package handlers

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

func (h *Handlers) GenerateUrl() string {
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

func (h *Handlers) Home(title string, c *gin.Context) {
	helpers.Render(c, gin.H{"title": title}, "home.html")
}

func (h *Handlers) DownloadResume(c *gin.Context) {
	h.DownloadFile(c, "Resume.pdf")
}

func (h *Handlers) DownloadDSA(c *gin.Context) {
	h.DownloadFile(c, "Data Structure And Algorithms specialization Certificate.pdf")
}

func (h *Handlers) DownloadDLS(c *gin.Context) {
	h.DownloadFile(c, "Deep Learning specialization Certificate.pdf")
}

func (h *Handlers) DownloadFile(c *gin.Context, fileName string) {
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
	c.Header("Content-Length", fmt.Sprint(helpers.FileSize(filePath)))

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

func (h *Handlers) LookingForWork(c *gin.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		print(err)
		return
	}
	LOOKING_FOR_WORK := os.Getenv("LOOKING_FOR_WORK")
	i, _ := strconv.Atoi(LOOKING_FOR_WORK)

	if i == 1 {
		helpers.Render(c, gin.H{}, "lookingForWork.html")
	} else {
		return
	}
}
