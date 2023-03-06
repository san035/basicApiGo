package storageminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/san035/basicApiGo/pkg/common"
)

// CreatBucketIfNotExist создание бакета
func CreatBucketIfNotExist(ctx context.Context, bucketName *string) error {
	exists, err := MinioClient.BucketExists(ctx, *bucketName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	// Создание нового бакета
	err = MinioClient.MakeBucket(ctx, *bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

// TrueTypeFile - возвращает признак, что тип файла верный
func TrueTypeFile(typeFile string) bool {
	return common.InArray(ListNameFolder, typeFile)
}
