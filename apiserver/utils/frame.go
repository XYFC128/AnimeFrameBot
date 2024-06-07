package utils

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Frame struct {
	Filename string
	Subtitle string
}

type FrameDistance struct {
	Frame    Frame
	Distance int
}

func InitFrames() []Frame {
	var frames []Frame
	imageDir := "./images"
	files, err := os.ReadDir(imageDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, file := range files {
		subtitle := strings.TrimSpace(strings.Split(file.Name(), ".")[0])
		frames = append(frames, Frame{Filename: file.Name(), Subtitle: subtitle})
	}

	return frames
}

func utf8len(s string) int {
	return len([]rune(s))
}

func MatchSubtitles(frames []Frame, input string, numFrames int) []Frame {
	var frameDistances []FrameDistance
	for _, frame := range frames {
		distance := fuzzy.LevenshteinDistance(input, frame.Subtitle)
		if distance < utf8len(frame.Subtitle) {
			frameDistances = append(frameDistances, FrameDistance{Frame: frame, Distance: distance})
		}
	}

	if len(frameDistances) == 0 {
		return nil
	}

	sort.Slice(frameDistances, func(i, j int) bool {
		return frameDistances[i].Distance < frameDistances[j].Distance
	})

	var matchedFrames []Frame
	for i := 0; i < min(numFrames, len(frameDistances)); i++ {
		matchedFrames = append(matchedFrames, frameDistances[i].Frame)
	}

	return matchedFrames	
}

func GetRandomFrames(frames []Frame, numFrames int) []Frame {
	randomIndices := rand.Perm(len(frames))[:numFrames]
	var randomFrames []Frame
	for _, index := range randomIndices {
		randomFrames = append(randomFrames, frames[index])
	}

	return randomFrames
}
