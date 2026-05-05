package attachment

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func setUpRouter() *gin.Engine {
	return gin.Default()
}

func newMultipartRequest(
	t *testing.T,
	fields map[string]string,
	fileFieldName string,
	fileName string,
	fileContentType string,
	fileContent []byte,
) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("write field %q: %v", k, err)
		}
	}

	if fileFieldName != "" {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="`+fileFieldName+`"; filename="`+fileName+`"`)
		if fileContentType != "" {
			header.Set("Content-Type", fileContentType)
		}
		part, err := writer.CreatePart(header)
		if err != nil {
			t.Fatalf("create part: %v", err)
		}
		if _, err := io.Copy(part, bytes.NewReader(fileContent)); err != nil {
			t.Fatalf("write file content: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req, _ := http.NewRequest("POST", "/api/attachment", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadAttachmentReturns201(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	fileContent := []byte("hello world")
	expected := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-123",
			FineName:    "receipt.png",
			ContentType: "image/png",
		},
		Content: fileContent,
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("AddAttachment", ctx, expected).Return(nil)

	req := newMultipartRequest(t,
		map[string]string{"budgetId": "budget-123"},
		"file", "receipt.png", "image/png", fileContent,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	facade.AssertCalled(t, "AddAttachment", ctx, expected)
}

func TestUploadAttachmentMissingFileReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{"budgetId": "budget-123"},
		"", "", "", nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "AddAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentMissingBudgetIdReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{},
		"file", "receipt.png", "image/png", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "AddAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentFacadeErrorReturns500(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("AddAttachment", ctx, mock.AnythingOfType("*attachment.Attachment")).Return(errors.New("boom"))

	req := newMultipartRequest(t,
		map[string]string{"budgetId": "budget-123"},
		"file", "receipt.png", "image/png", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
