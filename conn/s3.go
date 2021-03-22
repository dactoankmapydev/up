package conn

/*type S3Storage struct {
	s3session     *s3.S3
	AccessKey     string
	SecretKey     string
	Endpoint      string
	Region        string
	BucketName    string
	Location      string
}

func (storage *S3Storage) NewS3() {
	cred := credentials.NewStaticCredentials(storage.AccessKey, storage.SecretKey, "")
	_, err := cred.Get()
	if err != nil {
		logrus.Printf("bad credentials: %s", err)
	}
	storage.s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Credentials: cred,
		Endpoint: aws.String(storage.Endpoint),
		Region: aws.String(storage.Region),
	})))

	logrus.Info("successful connection to s3")
}*/
