package frame

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Frame struct {
	Filename string `json:"name"`
	Subtitle string `json:"subtitle"`
}

type FrameDistance struct {
	Frame
	Distance int
}

func initFrames(imageDir string) ([]Frame, error) {
	var frames []Frame
	files, err := os.ReadDir(imageDir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, file := range files {
		subtitle := strings.TrimSpace(strings.Split(file.Name(), ".")[0])
		frames = append(frames, Frame{Filename: file.Name(), Subtitle: subtitle})
	}

	return frames, nil
}

func utf8len(s string) int {
	return len([]rune(s))
}

func matchSubtitles(frames []Frame, input string, numFrames int) ([]Frame, error) {
	if numFrames > len(frames) || numFrames < 0 {
		return nil, fmt.Errorf("invalid number of frames: %d", numFrames)
	}

	var frameDistances []FrameDistance
	for _, frame := range frames {
		distance := fuzzy.LevenshteinDistance(input, frame.Subtitle)
		if distance < utf8len(frame.Subtitle) {
			frameDistances = append(frameDistances, FrameDistance{Frame: frame, Distance: distance})
		}
	}

	if len(frameDistances) == 0 {
		return []Frame{}, nil
	}

	sort.Slice(frameDistances, func(i, j int) bool {
		return frameDistances[i].Distance < frameDistances[j].Distance
	})

	var matchedFrames []Frame
	for i := 0; i < min(numFrames, len(frameDistances)); i++ {
		matchedFrames = append(matchedFrames, frameDistances[i].Frame)
	}

	return matchedFrames, nil
}

func matchSubtitlesExact(frames []Frame, input string, numFrames int) ([]Frame, error) {
	if numFrames > len(frames) || numFrames < 0 {
		return nil, fmt.Errorf("invalid number of frames: %d", numFrames)
	}

	var exactMatchedFrames []Frame
	for _, frame := range frames {
		if strings.Contains(strings.ToLower(frame.Subtitle), strings.ToLower(input)) {
			exactMatchedFrames = append(exactMatchedFrames, frame)
		}
	}

	if len(exactMatchedFrames) == 0 {
		return []Frame{}, nil
	}

	if len(exactMatchedFrames) < numFrames {
		return exactMatchedFrames, nil
	} else {
		return getRandomFrames(exactMatchedFrames, numFrames)
	}
}

func getRandomFrames(frames []Frame, numFrames int) ([]Frame, error) {
	if numFrames > len(frames) || numFrames < 0 {
		return nil, fmt.Errorf("invalid number of frames: %d", numFrames)
	}

	randomIndices := rand.Perm(len(frames))[:numFrames]
	var randomFrames []Frame
	for _, index := range randomIndices {
		randomFrames = append(randomFrames, frames[index])
	}

	return randomFrames, nil
}
