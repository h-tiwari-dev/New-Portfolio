package tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func TestCreateAdmin(t *testing.T) {
	jsonData := map[string]string{
		"username": os.Getenv("TEST_USERNAME"),
		"email":    os.Getenv("TEST_EMAIL"),
		"password": os.Getenv("TEST_PASSWORD"),
	}
	jsonValue, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", "http://localhost:8080/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Admin-Key", "Julia@1984")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func TestLogin(t *testing.T) {
	username := os.Getenv("TEST_USERNAME")
	password := os.Getenv("TEST_PASSWORD")

	jsonData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", "http://localhost:8080/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Parse the response body to extract the token
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	token := result["token"] // Assuming the token is returned with the key "token"

	os.Setenv("BEARER_TOKEN", token)
}

func TestCreateBlog(t *testing.T) {
	content, err := os.ReadFile("../blogs/md/pydantic.md")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Read image file
	imageContent, err := os.ReadFile("../blogs/md/assets/pydantic/img/blog_img.jpg")
	if err != nil {
		t.Fatalf("Failed to read image file: %v", err)
	}

	// Encode image data to base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageContent)

	// Create JSON data
	jsonData := map[string]interface{}{
		"title":       "Elevate Your LLM’s Game with Pydantic! Part I",
		"description": "",
		"content_md":  string(content),
		"images": []map[string]interface{}{
			{
				"filename":  "blog_img.jpg", // Replace with the actual filename
				"data":      imageBase64,
				"top_image": 1,
			},
		},
	}

	jsonValue, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", "http://localhost:8080/auth/blogs", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	// db.DB.Where("username = ?", os.Getenv("TEST_USERNAME")).Delete(&models.User{})
}

func TestGetBlogs(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080/blogs", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}
