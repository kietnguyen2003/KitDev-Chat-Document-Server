package infracstructure

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinIOStorageDocument struct {
	client     *minio.Client
	bucketName string
}

func NewMinIOStorageDocument(client *minio.Client, bucketName string) *MinIOStorageDocument {
	return &MinIOStorageDocument{
		client:     client,
		bucketName: bucketName,
	}
}

func (ds *MinIOStorageDocument) CreateBucket(ctx context.Context) error {
	exists, err := ds.client.BucketExists(ctx, ds.bucketName)
	if err != nil {
		fmt.Println("BucketExists error:", err)
		return err
	}

	if exists {
		fmt.Println("Bucket exists:", ds.bucketName)
		return nil
	}

	err = ds.client.MakeBucket(ctx, ds.bucketName, minio.MakeBucketOptions{})
	if err != nil {
		fmt.Println("MakeBucket error:", err)
		return err
	}

	fmt.Println("Create bucket success:", ds.bucketName)
	return nil
}

func (ds *MinIOStorageDocument) UploadReader(
	ctx context.Context,
	bucket string,
	object string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	info, err := ds.client.PutObject(
		ctx,
		bucket,
		object,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func (ds *MinIOStorageDocument) DeleteDocument(ctx context.Context, object string) error {
	return ds.client.RemoveObject(ctx, ds.bucketName, object, minio.RemoveObjectOptions{})
}
