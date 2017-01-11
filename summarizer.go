package goSummarizer

import (
	"fmt"
	"goSummarizer/helpers"
)

func Summarize(text string) string {
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

	return summary
}

func SummarizeFromURL(url string) string {
	extractedText, err := helpers.ExtractMainTextFromURL(url)
	if err != nil {
		return ""
	}

	var result = Summarize(extractedText)
	return result
}

func StartListening() {
	var text = helpers.ReadInputFromUser("Enter url or text for summarizing: ")
	var isURL = helpers.IsURL(text)
	var result = ""
	if isURL {
		result = SummarizeFromURL(text)
	} else {
		result = Summarize(text)
	}

	var answer = helpers.ReadInputFromUser("Do you want to store the summary to a file? (y/n)")
	var isPositiveAnswer = helpers.IsPositiveAnswer(answer)
	if !isPositiveAnswer {
		fmt.Println("Goodbye.")
		return
	}

	var path = helpers.ReadInputFromUser("Enter folder for storing the file: ")
	var stored = helpers.StoreTextToFile(path, result)
	if stored {
		fmt.Println("File stored!")
	} else {
		fmt.Println("Something went wrong! Please try again")
	}
}
