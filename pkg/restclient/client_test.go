package restclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"wn/pkg/applogger"
)

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

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "payload" {
			t.Errorf("body = %q, want payload", body)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := NewRestClient(noopLogger{}, false, false)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/test", io.NopCloser(stringReader("payload")))
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}

	body, status, err := client.MakeRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("MakeRequest() error = %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("status = %d, want 201", status)
	}
	if string(body) != `{"ok":true}` {
		t.Errorf("body = %q", body)
	}
}

type stringReader string

func (s stringReader) Read(p []byte) (int, error) {
	copy(p, s)
	return len(s), io.EOF
}

func TestGetURLFromUrl(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/api/v1/users", nil)
	got := getURLFromUrl(req)
	want := "https://example.com/api/v1/users"
	if got != want {
		t.Errorf("getURLFromUrl() = %q, want %q", got, want)
	}
}
