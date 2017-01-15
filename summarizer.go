package goSummarizer

import (
	"errors"
	"goSummarizer/helpers"
)

type Summarizer struct {
	url            string
	fullText       string
	summarizedText string
	summarized     bool
}

func CreateFromUrl(url string) *Summarizer {
	var summarizer = new(Summarizer)
	summarizer.url = url
	return summarizer
}

func CreateFromText(text string) *Summarizer {
	var summarizer = new(Summarizer)
	summarizer.fullText = text
	return summarizer
}

func (s *Summarizer) Summarize() (string, error) {
	if s.IsSummarized() {
		return s.summarizedText, nil
	}

	if s.fullText == "" && s.url == "" {
		return "", errors.New("You must submit text or url for summarizing")
	}

	if s.fullText != "" {
		s.summarizedText = s.summarizeFromText()
	} else if s.url != "" {
		extractedText, err := helpers.ExtractMainTextFromURL(s.url)
		if err != nil {
			return "", err
		}

		s.fullText = extractedText
		s.summarizedText = s.summarizeFromText()
	}

	s.summarized = true
	return s.summarizedText, nil
}

func (s *Summarizer) summarizeFromText() string {
	// Build the sentences dictionary
	var sentencesDictionary = helpers.GetSentencesRanks(s.fullText)
	// Build the summary with the sentences dictionary
	var summary = helpers.GetSummary(s.fullText, sentencesDictionary)

	return summary
}

func (s *Summarizer) GetSummaryInfo() (string, error) {
	if !s.IsSummarized() {
		return "", errors.New("You must first summarize the text in order to get information for it")
	}

	var summaryInfo = helpers.GetSummaryInfo(s.fullText, s.summarizedText)
	return summaryInfo, nil
}

func (s *Summarizer) IsSummarized() bool {
	return s.summarized
}

// func StartListening() {
// 	var text = helpers.ReadInputFromUser("Enter url or text for summarizing: ")
// 	var isURL = helpers.IsURL(text)
// 	var result = ""
// 	if isURL {
// 		result = SummarizeFromURL(text)
// 	} else {
// 		result = Summarize(text)
// 	}

// 	var answer = helpers.ReadInputFromUser("Do you want to store the summary to a file? (y/n)")
// 	var isPositiveAnswer = helpers.IsPositiveAnswer(answer)
// 	if !isPositiveAnswer {
// 		fmt.Println("Goodbye.")
// 		return
// 	}

// 	var path = helpers.ReadInputFromUser("Enter folder for storing the file: ")
// 	var stored = helpers.StoreTextToFile(path, result)
// 	if stored {
// 		fmt.Println("File stored!")
// 	} else {
// 		fmt.Println("Something went wrong! Please try again")
// 	}
// }
