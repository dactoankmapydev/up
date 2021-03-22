package repository

type UploadRepo interface {
	InsertEtag(string) error
}
