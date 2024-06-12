package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"AnimeFrameBot/internal/frame"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestGetEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		createImageDir bool
		wantStatus     int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "random frame normal",
			endpoint:       "/frame/random/3",
			createImageDir: true,
			wantStatus:     http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var frames []frame.Frame
				err := json.Unmarshal(body, &frames)
				assert.NoError(t, err)
				assert.Equal(t, 3, len(frames))
			},
		},
		{
			name:           "random frame bad count type",
			endpoint:       "/frame/random/asdf",
			createImageDir: true,
			wantStatus:     http.StatusBadRequest,
		},
		{
			name:           "random frame bad count value",
			endpoint:       "/frame/random/-1",
			createImageDir: true,
			wantStatus:     http.StatusBadRequest,
		},
		{
			name:           "random bad imageDir",
			endpoint:       "/frame/random/3",
			createImageDir: false,
			wantStatus:     http.StatusInternalServerError,
		},
		{
			name:           "fuzzy frame normal",
			endpoint:       "/frame/fuzzy/some/3",
			createImageDir: true,
			wantStatus:     http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var frames []frame.Frame
				err := json.Unmarshal(body, &frames)
				assert.NoError(t, err)
				assert.LessOrEqual(t, len(frames), 3)
			},
		},
		{
			name:           "fuzzy frame bad count type",
			endpoint:       "/frame/fuzzy/asdf/hjkl",
			createImageDir: true,
			wantStatus:     http.StatusBadRequest,
		},
		{
			name:           "fuzzy frame bad count value",
			endpoint:       "/frame/fuzzy/asdf/-1",
			createImageDir: true,
			wantStatus:     http.StatusBadRequest,
		},
		{
			name:           "fuzzy bad imageDir",
			endpoint:       "/frame/fuzzy/some/3",
			createImageDir: false,
			wantStatus:     http.StatusInternalServerError,
		},
		{
			name:           "exact frame normal",
			endpoint:       "/frame/exact/some/3",
			createImageDir: true,
			wantStatus:     http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var frames []frame.Frame
				err := json.Unmarshal(body, &frames)
				assert.NoError(t, err)
				assert.LessOrEqual(t, len(frames), 3)
			},
		},
		{
			name:           "exact frame bad count type",
			endpoint:       "/frame/exact/asdf/hjkl",
			wantStatus:     http.StatusBadRequest,
			createImageDir: true,
		},
		{
			name:           "exact frame bad count value",
			endpoint:       "/frame/exact/asdf/-1",
			wantStatus:     http.StatusBadRequest,
			createImageDir: true,
		},
		{
			name:           "exact bad imageDir",
			endpoint:       "/frame/exact/some/3",
			wantStatus:     http.StatusInternalServerError,
			createImageDir: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageDir := t.TempDir()
			if !tt.createImageDir {
				imageDir = filepath.Join(basepath, "NotExist")
			} else {
				for i := 0; i < 10; i++ {
					_, err := os.Create(filepath.Join(imageDir, strconv.Itoa(i)+".jpg"))
					require.NoError(t, err)
				}
			}

			server := NewServer(imageDir)

			req, err := http.NewRequest(http.MethodGet, tt.endpoint, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestRestPostUploadEndpoint(t *testing.T) {
	err := gofakeit.Seed(0)
	require.NoError(t, err)

	tests := []struct {
		name        string
		fileContent []byte
		fieldname   string
		filename    string
		wantStatus  int
		fileExists  bool
	}{
		{
			name:        "valid jpeg image",
			fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
			fieldname:   "image",
			filename:    "test.jpg",
			wantStatus:  http.StatusCreated,
			fileExists:  true,
		},
		{
			name:        "valid jpeg image",
			fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
			fieldname:   "image",
			filename:    "test.jpeg",
			wantStatus:  http.StatusCreated,
			fileExists:  true,
		},
		{
			name:        "valid png image",
			fileContent: gofakeit.ImagePng(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
			fieldname:   "image",
			filename:    "test.png",
			wantStatus:  http.StatusCreated,
			fileExists:  true,
		},
		{
			name:        "file larger than 10MB",
			fileContent: make([]byte, 10*1024*1024+1),
			fieldname:   "image",
			filename:    "test.jpg",
			wantStatus:  http.StatusBadRequest,
			fileExists:  false,
		},
		{
			name:        "fieldname not image",
			fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
			fieldname:   "file",
			filename:    "test.jpg",
			wantStatus:  http.StatusBadRequest,
			fileExists:  false,
		},
		{
			name:        "file not an image",
			fileContent: []byte("not an image"),
			fieldname:   "image",
			filename:    "test.jpg",
			wantStatus:  http.StatusBadRequest,
			fileExists:  false,
		},
		{
			name:        "change permission of imageDir",
			fileContent: []byte("\xFF\xD8\xFF"),
			fieldname:   "image",
			filename:    "test.jpg",
			wantStatus:  http.StatusInternalServerError,
			fileExists:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageDir := t.TempDir()
			if tt.name == "change permission of imageDir" {
				err := os.Chmod(imageDir, 0o000)
				require.NoError(t, err)
			}

			server := NewServer(imageDir)

			var b bytes.Buffer
			bw := multipart.NewWriter(&b)

			fw, err := bw.CreateFormFile(tt.fieldname, tt.filename)
			require.NoError(t, err)
			_, err = fw.Write(tt.fileContent)
			require.NoError(t, err)
			bw.Close()

			req, err := http.NewRequest(http.MethodPost, "/frame", &b)
			require.NoError(t, err)
			req.Header.Set("Content-Type", bw.FormDataContentType())
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestRestDownloadEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		wantStatus int
	}{
		{
			name:       "file exists",
			filename:   "test.jpg",
			wantStatus: http.StatusOK,
		},
		{
			name:       "file not exists",
			filename:   "notexist.jpg",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageDir := t.TempDir()
			if tt.name == "file exists" {
				err := os.WriteFile(filepath.Join(imageDir, tt.filename), []byte("test"), 0o644)
				require.NoError(t, err)
			}

			server := NewServer(imageDir)

			req, err := http.NewRequest(http.MethodGet, "/frame/"+tt.filename, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
