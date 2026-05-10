package s3

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/portal/i18n/domain/i18n"
	"gopkg.in/yaml.v3"
)

type S3BundleRepository struct {
	Bucket string
	Client *awss3.Client
	logger *logging.Logger
}

func NewS3BundleRepository(bucket string, client *awss3.Client) *S3BundleRepository {
	return &S3BundleRepository{
		Bucket: bucket,
		Client: client,
		logger: logging.GetLoggerInstanceForComponentByType(&S3BundleRepository{}),
	}
}

func ObjectKeyFor(key i18n.BundleKey) string {
	return fmt.Sprintf("%s/%s/message_bundle_%s.yaml", key.Application, key.Page, key.Language)
}

func (r *S3BundleRepository) GetBundle(ctx context.Context, key i18n.BundleKey) (*i18n.MessageBundle, error) {
	objectKey := ObjectKeyFor(key)
	out, err := r.Client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			return nil, i18n.ErrBundleNotFound
		}
		r.logger.LogErrorfFor("error fetching message bundle from s3 key %s: %v", objectKey, err)
		return nil, err
	}
	defer out.Body.Close()

	content, err := io.ReadAll(out.Body)
	if err != nil {
		r.logger.LogErrorfFor("error reading message bundle from s3 key %s: %v", objectKey, err)
		return nil, err
	}

	messages := map[string]any{}
	if err := yaml.Unmarshal(content, &messages); err != nil {
		r.logger.LogErrorfFor("error parsing message bundle from s3 key %s: %v", objectKey, err)
		return nil, err
	}

	return &i18n.MessageBundle{
		Application: key.Application,
		Page:        key.Page,
		Language:    key.Language,
		Messages:    messages,
	}, nil
}
