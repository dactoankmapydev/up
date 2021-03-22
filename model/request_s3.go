package model

import "github.com/aws/aws-sdk-go/service/s3"

type ReqUploadS3 struct {
	Key string `json:"key,omitempty"`
	UploadID string `json:"upload_id,omitempty"`
	Bucket string `json:"bucket,omitempty"`
	PartNum string `json:"part_num,omitempty"`
	BytePart string `json:"byte_part,omitempty"`
	Parts []*s3.CompletedPart `json:"parts,omitempty"`
}
