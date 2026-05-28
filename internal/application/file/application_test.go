package file

import (
	"context"
	"mime/multipart"
	"testing"
	"wn/internal/domain/dto"
	"wn/pkg/applogger"

	"github.com/google/uuid"
)

type mockFileService struct {
	filename string
	err      error
}

func (m *mockFileService) NewFile(_ context.Context, _ *multipart.FileHeader) (string, error) {
	return m.filename, m.err
}

type noopLogger struct{}

func (noopLogger) IsDebugLevel() bool                         { return false }
func (noopLogger) IsInfoLevel() bool                          { return false }
func (noopLogger) Debug(string)                               {}
func (noopLogger) Info(string)                                {}
func (noopLogger) Warn(string)                                {}
func (noopLogger) Error(string)                               {}
func (noopLogger) Warnf(string, ...any)                       {}
func (noopLogger) Errorf(string, ...any)                      {}
func (noopLogger) Debugf(string, ...any)                      {}
func (noopLogger) Infof(string, ...any)                       {}
func (l noopLogger) WithCtx(context.Context) applogger.Logger { return l }

type noopTx struct{}

func (noopTx) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestUploadFile(t *testing.T) {
	srv := NewService(noopTx{}, noopLogger{}, &mockFileService{filename: "abc.png"})

	resp, err := srv.UploadFile(context.Background(), uuid.New(), dto.UploadFileRequest{}, "https://host")
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
	want := "https://host/statics/images/abc.png"
	if resp.ImgUrl != want {
		t.Errorf("ImgUrl = %q, want %q", resp.ImgUrl, want)
	}
}
