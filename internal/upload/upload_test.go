package upload

import (
	"bytes"
	"testing"
	"net/http"
	"mime/multipart"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsImage(t *testing.T) {
	tests := []struct {
		name	string
		fileContent	[]byte
		contentType	string
		isImage	bool
	}{
		{
			name: "valid jpeg image",
			fileContent: []byte("\xFF\xD8\xFF"),
			contentType: "image/jpeg",
			isImage: true,
		},
		{
			name: "valid png image",
			fileContent: []byte("\x89\x50\x4E\x47\x0D\x0A\x1A\x0A"),
			contentType: "image/png",
			isImage: true,
		},
		{
			name: "valid gif image",
			fileContent: []byte("GIF87a"),
			contentType: "image/gif",
			isImage: true,
		},
		{
			name: "valid gif image",
			fileContent: []byte("GIF89a"),
			contentType: "image/gif",
			isImage: true,
		},
		{
			name: "invalid plain text",
			fileContent: []byte("invalid"),
			contentType: "text/plain; charset=utf-8",
			isImage: false,
		},
		{
			name: "invalid mp4 video",
			fileContent: []byte("\x66\x74\x79\x70"),
			contentType: "video/mp4",
			isImage: false,
		},
		{
			name: "invalid doc file",
			fileContent: []byte("\x0D\x44\x4F\x43"),
			contentType: "application/msword",
			isImage: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			bw := multipart.NewWriter(&b)

			fw, err := bw.CreateFormFile("image", "test.jpg")
			require.NoError(t, err)
			_, err = fw.Write(tt.fileContent)
			require.NoError(t, err)
			bw.Close()

			req, err := http.NewRequest(http.MethodPost, "/frame", &b)
			require.NoError(t, err)
			req.Header.Set("Content-Type", bw.FormDataContentType())

			file, _, err := req.FormFile("image")
			require.NoError(t, err)

			assert.Equal(t, tt.isImage, isImage(file))
		})
	}
}
