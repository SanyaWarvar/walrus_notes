package file

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"os"
	"testing"
	"wn/internal/domain/entity"
)

func TestEncodeFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.png")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	content := []byte("file-content")
	if _, err := part.Write(content); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("ReadForm() error = %v", err)
	}
	fileHeader := form.File["file"][0]

	encoded, err := encodeFile(fileHeader)
	if err != nil {
		t.Fatalf("encodeFile() error = %v", err)
	}

	decoded, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("DecodeString() error = %v", err)
	}
	if string(decoded) != "file-content" {
		t.Errorf("decoded = %q, want file-content", decoded)
	}
}

func TestCreateFile(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("./statics/images", 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	content := []byte("image-bytes")
	encoded := base64.RawStdEncoding.EncodeToString(content)

	err := createFile(&entity.StaticFile{
		Filename:     "test.png",
		FileAsString: encoded,
	})
	if err != nil {
		t.Fatalf("createFile() error = %v", err)
	}

	data, err := os.ReadFile("./statics/images/test.png")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "image-bytes" {
		t.Errorf("file content = %q, want image-bytes", data)
	}
}

func TestCreateFile_EmptyString(t *testing.T) {
	// errors.Wrap(nil, ...) возвращает nil — текущее поведение createFile
	err := createFile(&entity.StaticFile{Filename: "x.png", FileAsString: ""})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
