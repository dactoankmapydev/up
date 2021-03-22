package model

import "github.com/aws/aws-sdk-go/service/s3"

type DataUploadResp struct {
	Parts *s3.CompletedPart `json:"parts,omitempty"`
}
