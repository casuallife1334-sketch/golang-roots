package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOFileStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOFileStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIOFileStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(accessKey, secretKey, ""), Secure: useSSL})
	if err != nil {
		return nil, err
	}
	storage := &MinIOFileStorage{client: client, bucket: bucket}
	exists, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return storage, nil
}

func (s *MinIOFileStorage) Put(ctx context.Context, key string, file io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, key, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *MinIOFileStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
}

func (s *MinIOFileStorage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}
