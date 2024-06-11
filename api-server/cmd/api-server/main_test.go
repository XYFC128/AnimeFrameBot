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
	type args struct {
		endpoint       string
		method         string
		wantStatus     int
		createImageDir bool
		checkResponse  func(t *testing.T, body []byte)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "random frame normal",
			args: args{
				endpoint:       "/frame/random/3",
				method:         http.MethodGet,
				wantStatus:     http.StatusOK,
				createImageDir: true,
				checkResponse: func(t *testing.T, body []byte) {
					var frames []frame.Frame
					err := json.Unmarshal(body, &frames)
					require.NoError(t, err)
					assert.Equal(t, 3, len(frames))
				},
			},
		},
		{
			name: "random frame bad count type",
			args: args{
				endpoint:       "/frame/random/asdf",
				method:         http.MethodGet,
				wantStatus:     http.StatusBadRequest,
				createImageDir: true,
			},
		},
		{
			name: "random frame bad count value",
			args: args{
				endpoint:       "/frame/random/-1",
				method:         http.MethodGet,
				wantStatus:     http.StatusBadRequest,
				createImageDir: true,
			},
		},
		{
			name: "random bad imageDir",
			args: args{
				endpoint:       "/frame/random/3",
				method:         http.MethodGet,
				wantStatus:     http.StatusInternalServerError,
				createImageDir: false,
			},
		},
		{
			name: "fuzzy frame normal",
			args: args{
				endpoint:       "/frame/fuzzy/some/3",
				method:         http.MethodGet,
				wantStatus:     http.StatusOK,
				createImageDir: true,
				checkResponse: func(t *testing.T, body []byte) {
					var frames []frame.Frame
					err := json.Unmarshal(body, &frames)
					require.NoError(t, err)
					assert.LessOrEqual(t, len(frames), 3)
				},
			},
		},
		{
			name: "fuzzy frame bad count type",
			args: args{
				endpoint:       "/frame/fuzzy/asdf/hjkl",
				method:         http.MethodGet,
				wantStatus:     http.StatusBadRequest,
				createImageDir: true,
			},
		},
		{
			name: "fuzzy frame bad count value",
			args: args{
				endpoint:       "/frame/fuzzy/asdf/-1",
				method:         http.MethodGet,
				wantStatus:     http.StatusBadRequest,
				createImageDir: true,
			},
		},
		{
			name: "fuzzy bad imageDir",
			args: args{
				endpoint:       "/frame/fuzzy/some/3",
				method:         http.MethodGet,
				wantStatus:     http.StatusInternalServerError,
				createImageDir: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageDir := t.TempDir()
			if !tt.args.createImageDir {
				imageDir = filepath.Join(basepath, "NotExist")
			} else {
				for i := 0; i < 10; i++ {
					_, err := os.Create(filepath.Join(imageDir, strconv.Itoa(i)+".jpg"))
					require.NoError(t, err)
				}
			}

			server := NewServer(imageDir)

			req, err := http.NewRequest(tt.args.method, tt.args.endpoint, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tt.args.wantStatus, res.StatusCode)
			if tt.args.checkResponse != nil {
				tt.args.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestRestPostUploadEndpoint(t *testing.T) {
	err := gofakeit.Seed(0)
	require.NoError(t, err)
	type args struct {
		fileContent []byte
		fieldname   string
		filename    string
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		fileExists bool
	}{
		{
			name: "valid jpeg image",
			args: args{
				fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
				fieldname:   "image",
				filename:    "test.jpg",
			},
			wantStatus: http.StatusCreated,
			fileExists: true,
		},
		{
			name: "valid jpeg image",
			args: args{
				fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
				fieldname:   "image",
				filename:    "test.jpeg",
			},
			wantStatus: http.StatusCreated,
			fileExists: true,
		},
		{
			name: "valid png image",
			args: args{
				fileContent: gofakeit.ImagePng(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
				fieldname:   "image",
				filename:    "test.png",
			},
			wantStatus: http.StatusCreated,
			fileExists: true,
		},
		{
			name: "file larger than 10MB",
			args: args{
				fileContent: make([]byte, 10*1024*1024+1),
				fieldname:   "image",
				filename:    "test.jpg",
			},
			wantStatus: http.StatusBadRequest,
			fileExists: false,
		},
		{
			name: "fieldname not image",
			args: args{
				fileContent: gofakeit.ImageJpeg(gofakeit.IntRange(1, 10), gofakeit.IntRange(1, 10)),
				fieldname:   "file",
				filename:    "test.jpg",
			},
			wantStatus: http.StatusBadRequest,
			fileExists: false,
		},
		{
			name: "file not an image",
			args: args{
				fileContent: []byte("not an image"),
				fieldname:   "image",
				filename:    "test.jpg",
			},
			wantStatus: http.StatusBadRequest,
			fileExists: false,
		},
		{
			name: "change permission of imageDir",
			args: args{
				fileContent: []byte("\xFF\xD8\xFF"),
				fieldname:   "image",
				filename:    "test.jpg",
			},
			wantStatus: http.StatusInternalServerError,
			fileExists: true,
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

			fw, err := bw.CreateFormFile(tt.args.fieldname, tt.args.filename)
			require.NoError(t, err)
			_, err = fw.Write(tt.args.fileContent)
			require.NoError(t, err)
			bw.Close()

			req, err := http.NewRequest(http.MethodPost, "/frame", &b)
			require.NoError(t, err)
			req.Header.Set("Content-Type", bw.FormDataContentType())
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			savedPath := filepath.Join(imageDir, tt.args.filename)
			_, err = os.Stat(savedPath)
			assert.Equal(t, tt.fileExists, !os.IsNotExist(err))
		})
	}
}
