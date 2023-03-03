package storageminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/san035/basicApiGo/pkg/common"
	"github.com/san035/basicApiGo/pkg/logger"
	"time"
)

var MinioClient *minio.Client

const (
	NameProfileFolder  = `profile` // Название папки в хранилище
	NameDiagramFolder  = `diagram`
	NameDocFolder      = `doc`
	NameAllFile        = `*.*`
	NameAvatar         = `avatar`
	ExpiresPhotoSecond = time.Second * 60 * 60 * 24 * 7 // срок действия ссылки 1 год
)

func Init() (err error) {
	err = common.LoadConfig(&Config)
	if err != nil {
		return
	}
	setDefaultConfig()

	// создание соединения minio
	err = createConnect()
	return
}

// createConnect создание соединения minio
func createConnect() (err error) {
	var lastErr *error
	for _, confgMinIO := range Config.ListMinIO {
		MinioClient, err = MinioConnection(&confgMinIO)
		if err == nil {
			lastErr = nil
			break
		}
		lastErr = &err
	}

	if lastErr != nil {
		return logger.Wrap(lastErr)
	}

	// тест связи
	_, err = MinioClient.BucketExists(context.Background(), "test")
	if err != nil {
		return logger.Wrap(&err)
	}
	return
}

// MinioConnection func for opening minio connection.
func MinioConnection(confgMinIO *MinIOInstance) (*minio.Client, error) {
	const useSSL = false
	// Initialize minio client object.
	minioClient, err := minio.New(confgMinIO.Uri, &minio.Options{
		Creds:  credentials.NewStaticV4(confgMinIO.AccessKeyId, confgMinIO.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, err
}
