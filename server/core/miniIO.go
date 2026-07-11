package core

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinIO(cfg Config) (*minio.Client, error) {
	client, err := minio.New(
		cfg.MinioEndpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.MinioAccessKey,
				cfg.MinioSecretKey,
				"",
			),
			Secure: false,
		},
	)
	return client, err
}
