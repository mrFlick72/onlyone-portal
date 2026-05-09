package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

type S3AttachmentContentRepository struct {
	Bucket string
	Client *s3.Client
	logger *logging.Logger
}

func NewS3AttachmentContentRepository(bucket string, client *s3.Client) *S3AttachmentContentRepository {
	return &S3AttachmentContentRepository{
		Bucket: bucket,
		Client: client,
		logger: logging.GetLoggerInstanceForComponentByType(&S3AttachmentContentRepository{}),
	}
}

func (repository *S3AttachmentContentRepository) ObjectKeyFor(att *attachment.Attachment, partitionKey string) string {
	return fmt.Sprintf("%s/%s/%s", datePathFor(att.Date), partitionKey, att.AttachmentId)
}

func (repository *S3AttachmentContentRepository) FileLocationFor(objectKey string) string {
	return fmt.Sprintf("%s/%s", repository.Bucket, objectKey)
}

func (repository *S3AttachmentContentRepository) Save(ctx context.Context, objectKey string, contentType string, content []byte) error {
	_, err := repository.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(repository.Bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		repository.logger.LogErrorfFor("error uploading attachment content to s3 key %s: %v", objectKey, err)
	}
	return err
}

func (repository *S3AttachmentContentRepository) GetContentFor(ctx context.Context, objectKey string) ([]byte, error) {
	object, err := repository.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(repository.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		repository.logger.LogErrorfFor("error fetching attachment content from s3 key %s: %v", objectKey, err)
		return nil, err
	}
	defer object.Body.Close()

	content, err := io.ReadAll(object.Body)
	if err != nil {
		repository.logger.LogErrorfFor("error reading attachment content from s3 key %s: %v", objectKey, err)
		return nil, err
	}
	return content, nil
}

func datePathFor(d date.Date) string {
	return d.GetTime().Format("2006/01/02")
}
