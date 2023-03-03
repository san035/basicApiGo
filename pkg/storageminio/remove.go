package storageminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
)

// RemoveObject удаление файла из minio
func RemoveObject(ctxMinio context.Context, bucket, nameFileMinIO string) (update bool, err error) {
	update, err = MinioClient.BucketExists(ctxMinio, bucket)
	if err != nil {
		err = logger.WrapWithDeep1(&err)
		return
	}

	if !update {
		return
	}

	if nameFileMinIO != NameAllFile {
		err = MinioClient.RemoveObject(ctxMinio, bucket, nameFileMinIO, minio.RemoveObjectOptions{})
		return
	}

	// удаление всех файлов
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		// удаляемые объекты в objectsCh
		for object := range MinioClient.ListObjects(ctxMinio, bucket, minio.ListObjectsOptions{}) {
			if object.Err != nil {
				log.Error().Err(object.Err).Str("bucket", bucket).Msg("Delete bucket -")
				return
			}
			objectsCh <- object
		}
	}()

	// Call RemoveObjects API
	errorCh := MinioClient.RemoveObjects(ctxMinio, bucket, objectsCh, minio.RemoveObjectsOptions{})

	// Ожидание завершения удаления + сохранение последней ошибки в err
	for e := range errorCh {
		err = e.Err
		log.Error().Err(err).Str("bucket", bucket).Str("file", e.ObjectName).Msg("Delete bucket -")
	}
	if err != nil {
		return
	}

	err = MinioClient.RemoveBucket(ctxMinio, bucket)
	return
}
