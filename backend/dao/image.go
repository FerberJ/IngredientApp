package dao

import (
	"bytes"
	"context"
	"fmt"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
)

func AddImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	p := db.GetMinioProvider()

	newFilename := uuid.New()

	contentType := handler.Header.Get("Content-Type")

	_, err := p.Client.PutObject(context.Background(), "images", newFilename.String(), file, handler.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return newFilename.String(), nil
}

func DeleteImage(filename string) error {
	p := db.GetMinioProvider()

	err := p.Client.RemoveObject(context.Background(), "images", filename, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func AddImageFromURL(imageURL string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Download the image from URL
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// Get content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("URL does not point to an image (Content-Type: %s)", contentType)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Get minio provider
	p := db.GetMinioProvider()

	// Generate new filename
	newFilename := uuid.New().String()

	// Upload to Minio
	_, err = p.Client.PutObject(
		context.Background(),
		"images",
		newFilename,
		bytes.NewReader(imageData),
		int64(len(imageData)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Minio: %w", err)
	}

	return newFilename, nil
}

func GetImage(filename string, userId string) (io.Reader, error) {
	p := db.GetMinioProvider()
	filter := bson.M{"image": filename}

	recipeRepository := repository.NewRecipeRepository(db.GetMongoProvider())

	res, err := recipeRepository.FindDocument(filter, nil)
	if err != nil {
		return nil, err
	}
	data, _ := bson.Marshal(res)
	var recipe models.Recipe
	bson.Unmarshal(data, &recipe)

	if recipe.Private && userId != recipe.User {
		return nil, fmt.Errorf("Access denied")
	}

	obj, err := p.Client.GetObject(context.Background(), "images", filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}
