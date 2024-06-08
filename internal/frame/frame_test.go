package frame

import (
	"os"
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
			err := os.Mkdir(tt.imageDir, 0755)
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
			assert.Equal(t, tt.expectFrames, f)
		}
	}
}

func TestMatchSubtitles(t *testing.T) {
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
			input:       "appl",
			numFrames:   3,
			expectFrame: []Frame{
				{Filename: "a.png", Subtitle: "apple"},
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
