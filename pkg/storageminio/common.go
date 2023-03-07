package storageminio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/san035/basicApiGo/pkg/common"
	"github.com/san035/basicApiGo/pkg/logger"
	"net/url"
	"strconv"
	"time"
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

// CheckDataExpiresUrl Проверку времени действия ссылки
// Пример url: "http://bpm.dev.itkn.ru:9000/a2fb76a4-44b7-474d-bfaa-82c170d034b7/photo?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=cRLI6CRKLqVtqHxg%2F20230216%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230216T062228Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=bfadc8a9280048551fe97ccf5f7d872d747cfde76d706473cac6493c220e8f3d"
func CheckDataExpiresUrl(urlFile *string) error {
	const secondReserveUrl = int64(60 * 60 * 8) // 8 часов - минимальное кол-во секунд, сколько должен еще действовать url на запас
	if *urlFile == "" {
		return nil
	}

	u, err := url.Parse(*urlFile)
	if err != nil {
		return logger.WrapText(&err, "проверка url photo")
	}

	q := u.Query()
	dateCreate, dateExpires := q.Get("X-Amz-Date"), q.Get("X-Amz-Expires")

	timeCreate, err := time.Parse("20060102T150405Z", dateCreate)
	if err != nil {
		logger.Error(&err).Msg("получение из ссылки X-Amz-Date -")
		return err
	}

	secondExpiries, err := strconv.Atoi(dateExpires)
	if err != nil {
		logger.Error(&err).Msg("получение из ссылки X-Amz-Expires-")
		return err
	}

	if (timeCreate.Unix() + int64(secondExpiries)) < (time.Now().Unix() + secondReserveUrl) {
		return logger.New("ссылка photo устарела, дата создания: " + dateCreate)
	}

	return nil
}
