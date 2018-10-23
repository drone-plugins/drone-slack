package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func loadUserMapFromS3(bucket, key string) (map[string]string, error) {
	if bucket == "" || key == "" {
		return nil, nil
	}

	config := &aws.Config{
		Region: aws.String("us-west-2"),
	}

	sess := session.Must(session.NewSession(config))
	downloader := s3manager.NewDownloader(sess)

	requestInput := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	if _, err := downloader.Download(buf, &requestInput); err != nil {
		return nil, err
	}

	users := map[string]string{}
	if err := json.Unmarshal(buf.Bytes(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func translateOrReturn(key string, dict map[string]string) string {
	// nil dict will always return exists = false
	if value, exists := dict[key]; exists {
		return value
	}

	return key
}
