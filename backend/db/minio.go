package db

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioProvider struct {
	Client *minio.Client
}

var minioProvider *MinioProvider

func NewMinioProvider() *minio.Client {

	endpoint := "localhost:9000"
	accessKeyID := "BH7ksy5eGsnwm2sfFWLG"
	secretAccessKey := "HAAz9oz4sqzS3dznYbsg1r0gZijkJFAnFlGJZNPE"
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
