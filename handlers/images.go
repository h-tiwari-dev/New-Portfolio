package handlers

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"app/db"
	"app/helpers"
	"app/models"
)

func (h *Handlers) GetImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	var image models.Image
	if err := db.DB.First(&image, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	helpers.Render(
		c,
		gin.H{"image": base64.StdEncoding.EncodeToString(image.Data)},
		"image.html",
	)
	//c.JSON(http.StatusOK, gin.H{
	//	"image": base64.StdEncoding.EncodeToString(image.Data),
	//})
}
