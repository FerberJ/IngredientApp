package db

import (
	"context"
	"gotth/template/backend/configuration"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioProvider struct {
	Client *minio.Client
}

var minioProvider *MinioProvider

func NewMinioProvider(cfg configuration.Configutration) *minio.Client {

	endpoint := cfg.MinioEndpoint
	accessKeyID := cfg.MinioAccessKeyID
	secretAccessKey := cfg.MinioSecretAccessKey
	userSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: userSSL,
	})
	if err != nil {
		panic(err)
	}

	exists, err := minioClient.BucketExists(context.Background(), "images")
	if err != nil {
		// Handle error from BucketExists check
		log.Fatalf("Failed to check if bucket exists: %v", err)
	}

	if !exists {
		err := minioClient.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
		if err != nil {
			// Check if the error is because the bucket was created in the meantime
			if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" {
				log.Printf("Bucket 'images' already exists and is owned by you")
			} else {
				// It's a different error, so panic or handle appropriately
				log.Fatalf("Failed to create bucket: %v", err)
			}
		} else {
			log.Printf("Successfully created bucket 'images'")
		}
	} else {
		log.Printf("Bucket 'images' already exists")
	}

	setMinioProvider(&MinioProvider{Client: minioClient})
	return minioClient
}

func setMinioProvider(provider *MinioProvider) {
	minioProvider = provider
}

func GetMinioProvider() *MinioProvider {
	return minioProvider
}
