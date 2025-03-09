package dao

import (
	"context"
	"fmt"
	"gotth/template/backend/db"
	"gotth/template/backend/models"
	"gotth/template/backend/repository"
	"io"
	"mime/multipart"

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
