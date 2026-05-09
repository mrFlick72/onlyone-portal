//go:build test

package s3

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestMain(m *testing.M) {
	setupTestBucket()
	code := m.Run()
	teardownTestBucket()
	os.Exit(code)
}

func TestObjectKeyForBuildsDateScopedPath(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-001",
			Date:         testutils.SafeDateFor("15/03/2024"),
		},
	}

	key := repo.ObjectKeyFor(att, "budget-123_EXPENSE")
	assert.Equal(t, "2024/03/15/budget-123_EXPENSE/att-001", key)
}

func TestFileLocationForPrependsBucketName(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	location := repo.FileLocationFor("2024/03/15/budget-123_EXPENSE/att-001")
	assert.Equal(t, TestBucket+"/2024/03/15/budget-123_EXPENSE/att-001", location)
}

func TestSavePutsObjectInS3(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	objectKey := "2024/03/15/budget-123_EXPENSE/att-save-1"
	content := []byte("hello world attachment")

	err := repo.Save(context.Background(), objectKey, "text/plain", content)
	assert.Equal(t, nil, err)

	out, err := s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(TestBucket),
		Key:    aws.String(objectKey),
	})
	assert.Equal(t, nil, err)
	defer out.Body.Close()

	body, err := io.ReadAll(out.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, content, body)
	assert.Equal(t, "text/plain", aws.ToString(out.ContentType))
}

func TestGetContentForReturnsSavedContent(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	objectKey := "2024/04/20/budget-xyz_EXPENSE/att-get-1"
	content := []byte("retrieved content")
	err := repo.Save(context.Background(), objectKey, "text/plain", content)
	assert.Equal(t, nil, err)

	got, err := repo.GetContentFor(context.Background(), objectKey)
	assert.Equal(t, nil, err)
	assert.Equal(t, content, got)
}

func TestGetContentForReturnsErrorWhenObjectMissing(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	_, err := repo.GetContentFor(context.Background(), "does/not/exist")
	assert.NotEqual(t, nil, err)
}

func TestDeleteRemovesObjectFromS3(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	objectKey := "2024/05/01/budget-del_EXPENSE/att-del-1"
	content := []byte("to be deleted")
	assert.Equal(t, nil, repo.Save(context.Background(), objectKey, "text/plain", content))

	err := repo.Delete(context.Background(), objectKey)
	assert.Equal(t, nil, err)

	_, err = s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(TestBucket),
		Key:    aws.String(objectKey),
	})
	assert.NotEqual(t, nil, err)
}

func TestDeleteIsIdempotentWhenObjectMissing(t *testing.T) {
	repo := NewS3AttachmentContentRepository(TestBucket, s3Client)

	err := repo.Delete(context.Background(), "does/not/exist")
	assert.Equal(t, nil, err)
}
