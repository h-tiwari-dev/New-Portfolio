package handlers

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"app/db"
	"app/helpers"
	"app/models"
)

func (h *Handlers) BlogsPage(title string, c *gin.Context) {
	var blogs []models.Blog
	// Assuming you have a GORM DB instance named "db" initialized somewhere in your code

	// Fetch blogs from the database
	if err := db.DB.Find(&blogs).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve blogs from the database"},
		)
		return
	}
	helpers.Render(c, gin.H{"title": title, "blogs": blogs}, "blogs.html")
}

func (h *Handlers) BlogPage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	if err := db.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	contentTemplate, err := template.New("content").Parse(blog.ContentHTML)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML content"})
		return
	}

	// Execute the template to get the rendered HTML
	var renderedContent strings.Builder
	err = contentTemplate.Execute(&renderedContent, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render HTML content"})
		return
	}

	var image models.Image
	result := db.DB.Joins(
		"JOIN blog_images ON blog_images.image_id = images.id",
	).Where(
		"blog_images.blog_id = ? AND images.top_image= ?",
		blog.ID,
		1,
	).First(&image)

	if result == nil {
		h.logger.WithFields(logrus.Fields{
			"Blog ID": blog.ID,
		}).Info("Blog Image not found")
	}
	//var images []models.Image
	//result = db.DB.Find(&images, "blog_id = ?", blog.ID)

	helpers.Render(c,
		gin.H{
			"title":   blog.Title,
			"content": template.HTML(renderedContent.String()),
			"image":   base64.StdEncoding.EncodeToString(image.Data),
			//"imagesv
		},
		"blog.html",
	)
}

func (h *Handlers) CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the blog without images first
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	// Add all the images in the reques []models.BlogImage

	var blogImages []models.Image

	for _, blogImageData := range blog.Images {
		blogImage := models.Image{
			Filename: blogImageData.Filename,
			Data:     blogImageData.Data,
			TopImage: blogImageData.TopImage,
		}
		if err := db.DB.Create(&blogImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save blog images"})
			return
		} else {
			blogImages = append(blogImages, blogImage)
		}
	}

	blog.Images = blogImages

	if blog.ContentMD == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Content is Required."})
		return
	}

	image_ids, htmlBytes := helpers.ConvertMdToHTML([]byte(blog.ContentMD))
	blog.ContentHTML = string(htmlBytes)

	h.logger.WithFields(logrus.Fields{
		"image_ids": image_ids,
	}).Info("Images Saved")
	// Add a relation to blog

	// Save the blog to the database
	if err := db.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}

func (h *Handlers) GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	if err := db.DB.Find(&blogs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

func (h *Handlers) GetBlogById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	if err := db.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (h *Handlers) UpdateBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var existingBlog models.Blog
	if err := db.DB.First(&existingBlog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	var updatedBlog models.Blog
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingBlog.Title = updatedBlog.Title
	existingBlog.ContentMD = updatedBlog.ContentMD
	existingBlog.UpdatedAt = time.Now()

	if err := db.DB.Save(&existingBlog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully", "blog": existingBlog})
}

func (h *Handlers) DeleteBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	if err := db.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := db.DB.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}
