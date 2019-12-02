package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"path"
)

const (
	S3Bucket = "fm-demo-stt"
)

type s3Client struct {
	uploader *s3manager.Uploader
}

// initialize the aws s3 bucket upload service
func InitS3Bucket() {
	// create a new aws session
	sess := session.Must(session.NewSession())

	// create the s3 manager uploader
	u := s3manager.NewUploader(sess)

	// create the service client
	var client s3Client
	client.uploader = u

	StorageClient = &client
}

// function that takes a stream and upload to s3, return the location url
func (s *s3Client) Upload(id string, filename string, input io.Reader) (string, error) {
	// generate the s3 key which combines uuid and filename
	key := path.Join(id, filename)

	// upload to s3 bucket
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("private"),
		Body:   input,
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return result.Location, nil
}
