package storageminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"mime/multipart"
)

type FileOpen interface {
	Open()
}

// SaveFile сохранение файла в Minio
func SaveFile(ctxMinio context.Context, file *multipart.FileHeader, bucketName *string, folderMinIO string) error {
	// Get Buffer from file
	buffer, err := file.Open()
	if err != nil {
		return logger.WrapWithDeep1(&err)
	}
	defer func() {
		errBuffer := buffer.Close()
		if errBuffer != nil {
			logger.Error(&errBuffer).Msg("buffer.Close()-")
		}
	}()

	// Создание бакета
	err = CreatBucketIfNotExist(ctxMinio, bucketName)
	if err != nil {
		return logger.WrapWithDeep1(&err)
	}

	// Сохранение файла в minio
	folderMinIO = CreatFullNameFileMinio(folderMinIO, file.Filename)
	info, err := MinioClient.PutObject(ctxMinio, *bucketName, folderMinIO, buffer, file.Size, minio.PutObjectOptions{ContentType: file.Header["Content-Type"][0]})
	if err != nil {
		return logger.WrapWithDeep1(&err)
	}
	log.Debug().Str("Bucket", *bucketName).Str("file", file.Filename).Int64("size", info.Size).Msg("Uploaded +")
	return nil
}

func CreatFullNameFileMinio(folder, filename string) string {
	switch folder {
	case NameDiagramFolder:
		return folder + `/` + filename
	case NameProfileFolder:
		return folder + `/` + NameAvatar
	}
	return NameDocFolder + `/` + filename
}
