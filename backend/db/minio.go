package db

import (
	"context"
	"gotth/template/backend/configuration"

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

	exists, _ := minioClient.BucketExists(context.Background(), "images")
	if !exists {
		err := minioClient.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
		if err != nil {
			panic(err)
		}
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
