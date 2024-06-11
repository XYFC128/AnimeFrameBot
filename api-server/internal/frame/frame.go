package frame

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
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

func extractSubtitle(fileName string) string {
	parts := strings.Split(fileName, "_")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "_")
	}
	return fileName
}

func isValidFileName(filename string) bool {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return false
	}

	hashPart := parts[len(parts)-1]
	hashAndExt := strings.Split(hashPart, ".")
	if len(hashAndExt) != 2 {
		return false
	}

	hash := hashAndExt[0]
	ext := strings.ToLower(hashAndExt[1])

	if len(hash) != 64 {
		return false
	}

	validExt := ext == "jpg" || ext == "jpeg" || ext == "png"
	return validExt
}

func renameFileWithHash(imageDir string, fileName string) (string, error) {
	filePath := filepath.Join(imageDir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileBytes)
	hashString := hex.EncodeToString(hash[:])

	ext := filepath.Ext(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), ext)
	newFileName := baseName + "_" + hashString + ext

	err = os.Rename(filePath, filepath.Join(imageDir, newFileName))
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func initFrames(imageDir string) ([]Frame, error) {
	var frames []Frame
	files, err := os.ReadDir(imageDir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()

		if !isValidFileName(fileName) {
			newFileName, err := renameFileWithHash(imageDir, fileName)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			fileName = newFileName
		}

		subtitle := extractSubtitle(fileName)
		frames = append(frames, Frame{Filename: fileName, Subtitle: subtitle})
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
		if strings.EqualFold(input, frame.Subtitle) {
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
