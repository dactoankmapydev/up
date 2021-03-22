package repo_impl

import (
	"backup/conn"
	"backup/repository"
)

type UploadRepoImpl struct {
	client *conn.RedisDB
}

func NewBackupRepo(client *conn.RedisDB) repository.UploadRepo {
	return &UploadRepoImpl{
		client: client,
	}
}

func (upload *UploadRepoImpl) InsertEtag(etag string) error  {
	return nil
}
