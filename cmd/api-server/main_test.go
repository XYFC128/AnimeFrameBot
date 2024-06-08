package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"strconv"

	"AnimeFrameBot/internal/frame"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestRestGetEndpoints(t *testing.T) {
	type args struct {
		endpoint      string
		method        string
		wantStatus    int
		checkResponse func(t *testing.T, body []byte)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "random frame normal",
			args: args{
				endpoint:   "/frame/random/3",
				method:     http.MethodGet,
				wantStatus: http.StatusOK,
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
				endpoint:   "/frame/random/asdf",
				method:     http.MethodGet,
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name: "random frame bad count value",
			args: args{
				endpoint:   "/frame/random/-1",
				method:     http.MethodGet,
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name: "random bad imageDir",
			args: args{
				endpoint:   "/frame/random/3",
				method:     http.MethodGet,
				wantStatus: http.StatusInternalServerError,
			},
		},
		{
			name: "fuzzy frame normal",
			args: args{
				endpoint:   "/frame/fuzzy/some/3",
				method:     http.MethodGet,
				wantStatus: http.StatusOK,
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
				endpoint:   "/frame/fuzzy/asdf/hjkl",
				method:     http.MethodGet,
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name: "fuzzy frame bad count value",
			args: args{
				endpoint:   "/frame/fuzzy/asdf/-1",
				method:     http.MethodGet,
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name: "fuzzy bad imageDir",
			args: args{
				endpoint:   "/frame/fuzzy/some/3",
				method:     http.MethodGet,
				wantStatus: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server http.Handler
			if tt.name == "random bad imageDir" || tt.name == "fuzzy bad imageDir" {
				server = NewServer("invalid")
			} else {
				server = NewServer(filepath.Join(basepath, "images"))
			}
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
	tests := []struct {
		name        string
		fileContent []byte
		fieldname string
		filename    string
		wantStatus  int
		fileExists  bool
	}{
		{
			name:        "valid jpeg image",
			fileContent: []byte("\xFF\xD8\xFF"),
			fieldname:   "image",
			filename:    "test.jpg",
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
			fileContent: []byte("\xFF\xD8\xFF"),
			fieldname:   "file",
			filename:    "test.jpg",
			wantStatus:  http.StatusBadRequest,
			fileExists:  false,
		},
		{
			name:        "file not an image",
			fileContent: []byte("invalid"),
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
				err := os.Chmod(imageDir, 0000)
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

			savedPath := filepath.Join(imageDir, tt.filename)
			_, err = os.Stat(savedPath)
			assert.Equal(t, tt.fileExists, !os.IsNotExist(err))
		})
	}
}

func FuzzRandomFrame(f *testing.F) {
	server := NewServer(filepath.Join(basepath, "images"))
	files, err := os.ReadDir(filepath.Join(basepath, "images"))
	require.NoError(f, err)
	imageCount := len(files)

	f.Fuzz(func(t *testing.T, count int) {
		req, err := http.NewRequest(http.MethodGet, "/frame/random/" + strconv.Itoa(count), nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		res := w.Result()

		if count >= 0 && count <= imageCount {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			var frames []frame.Frame
			e := json.Unmarshal(w.Body.Bytes(), &frames)
			require.NoError(t, e)
			assert.LessOrEqual(t, len(frames), count)
		} else {
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		}
	})
}