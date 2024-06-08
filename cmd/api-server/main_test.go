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

	"AnimeFrameBot/internal/frame"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Join(filepath.Dir(b), "../..")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(filepath.Join(basepath, "images"))
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
	imageDir := t.TempDir()
	server := NewServer(imageDir)

	var b bytes.Buffer
	bw := multipart.NewWriter(&b)

	fw, err := bw.CreateFormFile("image", "test.jpg")
	require.NoError(t, err)
	_, err = fw.Write([]byte("\xFF\xD8\xFF"))
	require.NoError(t, err)
	bw.Close()

	req, err := http.NewRequest(http.MethodPost, "/frame/", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", bw.FormDataContentType())
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
	res := w.Result()
	assert.Equalf(t, http.StatusCreated, res.StatusCode, "%s", w.Body.String())

	savedPath := filepath.Join(imageDir, "test.jpg")
	_, err = os.Stat(savedPath)
	assert.Equal(t, false, os.IsNotExist(err))
}
