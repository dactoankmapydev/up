package model

import "github.com/aws/aws-sdk-go/service/s3"

type ResponseDataS3 struct {
	StatusCode int64 `json:"status_code"`
	Message   string  `json:"message"`
	UploadID string `json:"upload_id,omitempty"`
	Key  string  `json:"key,omitempty"`
	Bucket string `json:"bucket,omitempty"`
	Etag  string   `json:"etag,omitempty"`
	PartNumber int `json:"part_number,omitempty"`
	Parts []*s3.CompletedPart `json:"parts,omitempty"`
}
