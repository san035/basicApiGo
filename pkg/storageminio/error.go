package storageminio

const (
	BucketDoesNotExist = `The specified bucket does not exist` // Не допускать изменение, т.к. сравнение на текст после вызова PresignedGetObject
)
