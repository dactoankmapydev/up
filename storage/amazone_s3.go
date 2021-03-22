package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func GetS3Session() *session.Session {
	cred := credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), "")
	_, err := cred.Get()
	if err != nil {
		log.Printf("bad credentials: %s\n", err)
	}

	return session.Must(session.NewSession(
		&aws.Config{
			Credentials: cred,
			Endpoint: aws.String(os.Getenv("ENDPOINT")),
			Region: aws.String(os.Getenv("REGION")),
		}))
}
func GetS3Service() *s3.S3 {
	return s3.New(GetS3Session())
}

func CreateBucketS3(bucketName string) *s3.CreateBucketOutput {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(os.Getenv("LOCATION")),
		},
	}
	bucket, err := GetS3Service().CreateBucket(input)
	if err != nil {
		log.Println(err)
	}
	return bucket
}

func InitUploadS3(bucket, key string) (string, string, string) {
	createResp, err := GetS3Service().CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
	})
	if err != nil {
		log.Println(err)
	}
	return *createResp.UploadId, *createResp.Key, *createResp.Bucket
}

func UploadPartS3(bucket, key, uploadID string, fileBytes []byte, partNum int) (completedPart *s3.CompletedPart, err error) {
	var try int
	for try <= 10 {
		uploadResp, err := GetS3Service().UploadPart(&s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        &bucket,
			Key:           &key,
			PartNumber:    aws.Int64(int64(partNum)),
			UploadId:      &uploadID,
			ContentLength: aws.Int64(int64(len(fileBytes))),
		})
		// Upload failed
		if err != nil {
			log.Println("request error", err)
			// Max retries reached! Quitting
			if try == 10 {
				return nil, err
			} else {
				// Retrying
				log.Println("request retrying")
				try++
			}
		} else {
			// Upload is done!
			return &s3.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: aws.Int64(int64(partNum)),
			}, nil
		}
	}
	return nil, nil
}

func UploadCompleteS3(bucket, key, uploadID string, parts []*s3.CompletedPart) string {
	resp, err := GetS3Service().CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: &uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		log.Println(err)
	}
	return *resp.ETag
}

func UploadPartAbortS3(bucket, key, uploadID string) {
	_, err := GetS3Service().AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   &bucket,
		Key: &key,
		UploadId: &uploadID,
	})
	if err != nil {
		log.Println(err)
	}
}



