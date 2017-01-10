package goSummarizer

import "goSummarizer/helpers"
import "fmt"

func Summarize(text string) {
	// Build the sentences dictionary
	var sentencesDictionary = helpers.GetSentencesRanks(text)

	// Build the summary with the sentences dictionary
	var summary = helpers.GetSummary(text, sentencesDictionary)
	fmt.Println(summary)

	// // Print the summary
	// fmt.Println("summary is\n", summary)

	// // Print the ratio between the summary length and the original length
	// fmt.Println("")
	// fmt.Println("Original Length: ", len(text))
	// fmt.Println("Summary Length: ", len(summary))
	// fmt.Println("Summary Ratio: ", (100 - (100 * (len(summary) / len(text)))))
}

func SummarizeFromURL(url string) {
	extractedText, err := helpers.ExtractMainTextFromURL(url)
	if err != nil {
		return
	}

	Summarize(extractedText)
}
