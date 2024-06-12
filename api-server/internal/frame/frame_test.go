package frame

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtf8len(t *testing.T) {
	tests := []struct {
		input     string
		expectLen int
	}{
		{input: "Hello, World!", expectLen: 13},
		{input: "Hello, 世界!", expectLen: 10},
		{input: "こんにちは, 世界!", expectLen: 10},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expectLen, utf8len(tt.input))
	}
}

func TestExtractSubtitle(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{input: "a.png", expect: "a.png"},
		{input: "a_b.png", expect: "a"},
		{input: "a_b_c.png", expect: "a_b"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expect, extractSubtitle(tt.input))
	}
}

func TestIsValidFileName(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{input: "a.png", expect: false},
		{input: "a_b.png", expect: false},
		{input: "a_b", expect: false},
		{input: "a_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.png", expect: true},
		{input: "a_b_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.png", expect: true},
		{input: "a_b_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.jpg", expect: true},
		{input: "a_b_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.jpeg", expect: true},
		{input: "a_b_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA.txt", expect: false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expect, isValidFileName(tt.input))
	}
}

func TestRenameFileWithHash(t *testing.T) {
	tests := []struct {
		imageDir    string
		fileName    string
		expectError string
	}{
		{
			imageDir: "test_images",
			fileName: "a.png",
		},
		{
			imageDir:    "nonexistent",
			fileName:    "a.png",
			expectError: "open nonexistent/a.png: no such file or directory",
		},
	}

	for _, tt := range tests {
		if tt.imageDir != "nonexistent" {
			err := os.Mkdir(tt.imageDir, 0o755)
			assert.NoError(t, err)
			f, err := os.Create(tt.imageDir + "/" + tt.fileName)
			assert.NoError(t, err)
			f.Close()
		}

		_, err := renameFileWithHash(tt.imageDir, tt.fileName)
		if tt.expectError != "" {
			assert.EqualError(t, err, tt.expectError)
		} else {
			assert.NoError(t, err)
			e := os.RemoveAll(tt.imageDir)
			assert.NoError(t, e)
		}
	}
}

func TestInitFrames(t *testing.T) {
	tests := []struct {
		imageDir     string
		expectFrames []Frame
		expectError  string
	}{
		{
			imageDir: "test_images",
			expectFrames: []Frame{
				{Filename: "a.png", Subtitle: "a"},
				{Filename: "b.png", Subtitle: "b"},
				{Filename: "c.png", Subtitle: "c"},
				{Filename: "d.png", Subtitle: "d"},
				{Filename: "e.png", Subtitle: "e"},
			},
		},
		{
			imageDir:    "nonexistent",
			expectError: "open nonexistent: no such file or directory",
		},
	}

	for _, tt := range tests {
		if tt.imageDir != "nonexistent" {
			err := os.Mkdir(tt.imageDir, 0o755)
			assert.NoError(t, err)
			for i := 'a'; i <= 'e'; i++ {
				f, err := os.Create(tt.imageDir + "/" + string(i) + ".png")
				assert.NoError(t, err)
				f.Close()
			}
		}

		f, err := initFrames(tt.imageDir)
		if tt.expectError != "" {
			assert.EqualError(t, err, tt.expectError)
		} else {
			e := os.RemoveAll(tt.imageDir)
			assert.NoError(t, e)
			assert.NoError(t, err)

			for index, frame := range f {
				fileName := frame.Filename
				ext := filepath.Ext(fileName)
				baseName := strings.TrimSuffix(fileName, ext)
				realBaseName := strings.Split(baseName, "_")[0]
				realFileName := realBaseName + ext
				assert.Equal(t, realFileName, tt.expectFrames[index].Filename)
			}
		}
	}
}

func TestMatchSubtitles(t *testing.T) {
	frames := []Frame{
		{Filename: "a.png", Subtitle: "apple"},
		{Filename: "b.png", Subtitle: "banana"},
		{Filename: "c.png", Subtitle: "grape"},
		{Filename: "d.png", Subtitle: "apple"},
	}

	tests := []struct {
		input       string
		numFrames   int
		expectFrame []Frame
		expectError string
	}{
		{
			input:     "appl",
			numFrames: 2,
			expectFrame: []Frame{
				{Filename: "a.png", Subtitle: "apple"},
				{Filename: "d.png", Subtitle: "apple"},
			},
		},
		{
			input:     "appl",
			numFrames: 3,
			expectFrame: []Frame{
				{Filename: "a.png", Subtitle: "apple"},
				{Filename: "d.png", Subtitle: "apple"},
				{Filename: "c.png", Subtitle: "grape"},
			},
		},
		{
			input:     "appl",
			numFrames: 4,
			expectFrame: []Frame{
				{Filename: "a.png", Subtitle: "apple"},
				{Filename: "d.png", Subtitle: "apple"},
				{Filename: "c.png", Subtitle: "grape"},
				{Filename: "b.png", Subtitle: "banana"},
			},
		},
		{
			input:       "b",
			numFrames:   1,
			expectFrame: []Frame{{Filename: "b.png", Subtitle: "banana"}},
		},
		{
			input:       "",
			numFrames:   1,
			expectFrame: []Frame{},
		},
		{
			input:       "b",
			numFrames:   0,
			expectFrame: []Frame{},
		},
		{
			input:       "",
			numFrames:   -1,
			expectError: "invalid number of frames: -1",
		},
		{
			input:       "",
			numFrames:   5,
			expectError: "invalid number of frames: 5",
		},
	}

	for _, tt := range tests {
		f, err := matchSubtitles(frames, tt.input, tt.numFrames)
		if tt.expectError != "" {
			assert.EqualError(t, err, tt.expectError)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectFrame, f)
		}
	}
}

func TestMatchSubtitlesExact(t *testing.T) {
	frames := []Frame{
		{Filename: "a.png", Subtitle: "apple"},
		{Filename: "b.png", Subtitle: "banana"},
		{Filename: "c.png", Subtitle: "grape"},
		{Filename: "d.png", Subtitle: "peach"},
	}

	tests := []struct {
		input       string
		numFrames   int
		expectFrame []Frame
		expectError string
	}{
		{
			input:       "apple",
			numFrames:   1,
			expectFrame: []Frame{{Filename: "a.png", Subtitle: "apple"}},
		},
		{
			input:       "banana",
			numFrames:   1,
			expectFrame: []Frame{{Filename: "b.png", Subtitle: "banana"}},
		},
		{
			input:       "unmatch",
			numFrames:   1,
			expectFrame: []Frame{},
		},
		{
			input:       "apple",
			numFrames:   0,
			expectFrame: []Frame{},
		},
		{
			input:       "apple",
			numFrames:   4,
			expectFrame: []Frame{{Filename: "a.png", Subtitle: "apple"}},
		},
		{
			input:       "apple",
			numFrames:   -1,
			expectError: "invalid number of frames: -1",
		},
		{
			input:       "apple",
			numFrames:   5,
			expectError: "invalid number of frames: 5",
		},
	}

	for _, tt := range tests {
		f, err := matchSubtitlesExact(frames, tt.input, tt.numFrames)
		if tt.expectError != "" {
			assert.EqualError(t, err, tt.expectError)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectFrame, f)
		}
	}
}

func TestGetRandomFrames(t *testing.T) {
	var frames []Frame
	for i := 'a'; i <= 'e'; i++ {
		frames = append(frames, Frame{Filename: string(i) + ".png", Subtitle: string(i)})
	}
	tests := []struct {
		numFrames   int
		expectError string
	}{
		{numFrames: 0},
		{numFrames: 3},
		{numFrames: 5},
		{
			numFrames:   -1,
			expectError: "invalid number of frames: -1",
		},
		{
			numFrames:   6,
			expectError: "invalid number of frames: 6",
		},
	}

	for _, tt := range tests {
		f, err := getRandomFrames(frames, tt.numFrames)
		if tt.expectError != "" {
			assert.EqualError(t, err, tt.expectError)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.numFrames, len(f))
		}
	}
}

func FuzzGetRandomFrames(f *testing.F) {
	var frames []Frame
	for i := 'a'; i <= 'e'; i++ {
		frames = append(frames, Frame{Filename: string(i) + ".png", Subtitle: string(i)})
	}
	f.Fuzz(func(t *testing.T, numFrames int) {
		frames, err := getRandomFrames(frames, numFrames)
		if numFrames < 0 || numFrames > len(frames) {
			assert.EqualError(t, err, "invalid number of frames: "+fmt.Sprint(numFrames))
		} else {
			assert.NoError(t, err)
			assert.Equal(t, numFrames, len(frames))
		}
	})
}
