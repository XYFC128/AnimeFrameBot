package main

import (
	"fmt"

    "apiserver/utils"
)


func main() {
    frames := utils.InitFrames()
	for true {
		var input string
		fmt.Print("Enter a subtitle: ")
		fmt.Scanln(&input)

		matchedFrames := utils.MatchSubtitles(frames, input, 3)
		for _, frame := range matchedFrames {
			fmt.Println("Matched frame: ", frame.Filename)
		}

		randomFrames := utils.GetRandomFrames(frames, 3)
		for _, frame := range randomFrames {
			fmt.Println("Random frame: ", frame.Filename)
		}
	}
}