package upload

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsImage(t *testing.T) {
	tests := []struct {
		input   string
		isImage bool
	}{
		{input: "image.jpg", isImage: true},
		{input: "image.jpeg", isImage: true},
		{input: "image.png", isImage: true},
		{input: "image.gif", isImage: true},
		{input: "image.txt", isImage: false},
		{input: "image.mp4", isImage: false},
		{input: "image.doc", isImage: false},
		{input: "image.sh", isImage: false},
	}

	for _, tt := range tests {
		handler := &multipart.FileHeader{
			Filename: tt.input,
		}
		assert.Equal(t, tt.isImage, isImage(handler))
	}
}
