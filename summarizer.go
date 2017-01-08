package goSummarizer

import "fmt"
import "goSummarizer/helpers"

func Summarize(text string) {
	fmt.Println(text)
}

func SummarizeFromURL(url string) {
	extractedText, err := helpers.ExtractMainTextFromURL(url)
	if err != nil {
		return
	}

	Summarize(extractedText)
}
