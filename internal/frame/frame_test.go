package frame

import (
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

func TestMatchSubtitles(t *testing.T) {
	var frames []Frame
	for i := 'a'; i <= 'e'; i++ {
		frames = append(frames, Frame{Filename: string(i) + ".png", Subtitle: string(i)})
	}
	tests := []struct {
		input       string
		numFrames   int
		expectFrame []Frame
		expectError string
	}{
		{
			input:       "a",
			numFrames:   1,
			expectFrame: []Frame{{Filename: "a.png", Subtitle: "a"}},
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
