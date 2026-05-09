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
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func mustDate(t *testing.T, raw string) date.Date {
	t.Helper()
	d, err := date.DateFor(raw)
	if err != nil {
		t.Fatalf("parse date %q: %v", raw, err)
	}
	return *d
}

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

func TestUploadAttachmentReturns204(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	fileContent := []byte("hello world")
	expected := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-123",
			BudgetType:  "expense",
			Date:        mustDate(t, "15/03/2024"),
			FineName:    "receipt.txt",
			ContentType: "text/plain",
		},
		Content: fileContent,
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("SaveAttachment", ctx, expected).Return(nil)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":   "budget-123",
			"budgetType": "expense",
			"date":       "15/03/2024",
		},
		"file", "receipt.txt", "text/plain", fileContent,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "SaveAttachment", ctx, expected)
}

func TestUploadAttachmentWithAttachmentIdUpdatesExistingAttachment(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	fileContent := []byte("updated content")
	expected := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "attachment-789",
			BudgetId:     "budget-123",
			BudgetType:   "expense",
			Date:         mustDate(t, "15/03/2024"),
			FineName:     "receipt.txt",
			ContentType:  "text/plain",
		},
		Content: fileContent,
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("SaveAttachment", ctx, expected).Return(nil)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":     "budget-123",
			"budgetType":   "expense",
			"date":         "15/03/2024",
			"attachmentId": "attachment-789",
		},
		"file", "receipt.txt", "text/plain", fileContent,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "SaveAttachment", ctx, expected)
}

func TestUploadAttachmentMissingFileReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":   "budget-123",
			"budgetType": "expense",
			"date":       "15/03/2024",
		},
		"", "", "", nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "SaveAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentMissingBudgetIdReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetType": "expense",
			"date":       "15/03/2024",
		},
		"file", "receipt.txt", "text/plain", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "SaveAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentMissingBudgetTypeReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId": "budget-123",
			"date":     "15/03/2024",
		},
		"file", "receipt.txt", "text/plain", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "SaveAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentMissingDateReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":   "budget-123",
			"budgetType": "expense",
		},
		"file", "receipt.txt", "text/plain", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "SaveAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentInvalidDateReturns400(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":   "budget-123",
			"budgetType": "expense",
			"date":       "2024-03-15",
		},
		"file", "receipt.txt", "text/plain", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "SaveAttachment", mock.Anything, mock.Anything)
}

func TestUploadAttachmentFacadeErrorReturns500(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("SaveAttachment", ctx, mock.AnythingOfType("*attachment.Attachment")).Return(errors.New("boom"))

	req := newMultipartRequest(t,
		map[string]string{
			"budgetId":   "budget-123",
			"budgetType": "expense",
			"date":       "15/03/2024",
		},
		"file", "receipt.txt", "text/plain", []byte("hello"),
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAttachmentMetadataReturns200WithRepresentation(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("GetAttachments", ctx, "budget-123", attachment.Expense).Return([]attachment.AttachmentMetadata{
		{AttachmentId: "a-1", BudgetId: "budget-123", BudgetType: "expense", Owner: "USER", FineName: "r.pdf"},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/attachment/metadata/expense/budget-123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte(`"attachmentId":"a-1"`)) || !bytes.Contains([]byte(body), []byte(`"fileName":"r.pdf"`)) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestGetAttachmentMetadataFacadeErrorReturns500(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("GetAttachments", ctx, "budget-x", attachment.Revenue).Return(nil, errors.New("db down"))

	req, _ := http.NewRequest(http.MethodGet, "/api/attachment/metadata/revenue/budget-x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAttachmentContentReturns200WithContent(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	expected := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-1",
			FineName:     "doc.pdf",
			ContentType:  "application/pdf",
		},
		Content: []byte("PDFBYTES"),
	}
	facade.On("GetAttachmentBy", ctx, "att-1").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/attachment/att-1/content", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Equal(t, `attachment; filename="doc.pdf"`, w.Header().Get("Content-Disposition"))
	assert.Equal(t, "PDFBYTES", w.Body.String())
}

func TestGetAttachmentContentDefaultsToOctetStreamWhenNoContentType(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	expected := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{AttachmentId: "att-2", FineName: "blob"},
		Content:            []byte("X"),
	}
	facade.On("GetAttachmentBy", ctx, "att-2").Return(expected, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/attachment/att-2/content", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
}

func TestGetAttachmentContentFacadeErrorReturns500(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("GetAttachmentBy", ctx, "missing").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/attachment/missing/content", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteAttachmentReturns204(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("DeleteAttachment", ctx, "att-1").Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/attachment/att-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "DeleteAttachment", ctx, "att-1")
}

func TestDeleteAttachmentFacadeErrorReturns500(t *testing.T) {
	r := setUpRouter()
	facade := new(AttachmentActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterAttachmentEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("DeleteAttachment", ctx, "missing").Return(errors.New("not found"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/attachment/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
