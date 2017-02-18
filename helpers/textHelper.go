package helpers

import (
	"regexp"
	"strings"
)

// Naive method for splitting a text into sentences
func getContentSentences(content string) []string {
	content = strings.Replace(content, "\n", ". ", -1)
	var sentences = strings.Split(content, ". ")

	// filter invalid sentences
	var validSentences = []string{}
	for _, sentence := range sentences {
		if strings.TrimSpace(sentence) != "" {
			validSentences = append(validSentences, sentence)
		}
	}

	return validSentences
}

// Naive method for splitting a text into paragraphs
func getContentParagraphs(content string) []string {
	var paragraphs = strings.Split(content, "\n\n")

	// filter invalid paragraphs
	var validParagraphs = []string{}
	for _, paragraph := range paragraphs {
		if strings.TrimSpace(paragraph) != "" {
			validParagraphs = append(validParagraphs, paragraph)
		}
	}

	return validParagraphs
}

// Caculate the intersection between 2 sentences
func sentencesIntersectedWordsCount(sent1 string, sent2 string) float32 {
	// split the sentence into words/tokens
	var words1 = splitWordsToMap(sent1)
	var words2 = splitWordsToMap(sent2)

	// If there is not intersection, just return 0
	if len(words1) == 0 && len(words2) == 0 {
		return 0
	}

	var intersectionCount = 0
	for word1 := range words1 {
		var _, exists = words2[word1]
		if exists {
			intersectionCount++
		}
	}

	// We normalize the result by the average number of words
	var numerator = float32(intersectionCount)
	var denominator = float32((len(words1) + len(words2)) / 2)
	var result = numerator / denominator
	return result
}

// Split words from string into map object with words as keys and true as value to all
func splitWordsToMap(text string) map[string]bool {
	var words = strings.Split(text, " ")
	var wordsMap = make(map[string]bool)
	for _, word := range words {
		if _, exists := wordsMap[word]; !exists {
			wordsMap[word] = true
		}
	}

	return wordsMap
}

// Format a sentence - remove all non-alphbetic chars from the sentence
// We'll use the formatted sentence as a key in our sentences dictionary
func formatSentence(sentence string) string {
	var regex, err = regexp.Compile("[^a-zA-Zа-яА-я]")
	if err != nil {
		return ""
	}

	var replacedSentence = regex.ReplaceAllString(sentence, "")
	return replacedSentence
}

func getSentencesRanks(content string) map[string]float32 {
	// Split the content into sentences
	var sentences = getContentSentences(content)

	// Calculate the intersection of every two sentences
	var sentencesCount = len(sentences)
	var values = [][]float32{}

	for i := 0; i < sentencesCount; i++ {
		values = append(values, []float32{})
		for j := 0; j < sentencesCount; j++ {
			if i == j {
				values[i] = append(values[i], 0)
			} else {
				values[i] = append(values[i], sentencesIntersectedWordsCount(sentences[i], sentences[j]))
			}
		}
	}

	// Build the sentences dictionary
	// The score of a sentences is the sum of all its intersection
	var sentencesDictionary = make(map[string]float32)
	for i := 0; i < sentencesCount; i++ {
		var score float32

		for j := 0; j < sentencesCount; j++ {
			if i == j {
				continue
			}
			score += values[i][j]
		}

		sentencesDictionary[formatSentence(sentences[i])] = score
	}

	return sentencesDictionary
}

// Return the best sentence in a paragraph
func getBestSentence(paragraph string, sentencesDictionary map[string]float32) string {

	// Split the paragraph into sentences
	var sentences = getContentSentences(paragraph)

	// Ignore short paragraphs
	if len(sentences) < 2 {
		return ""
	}

	// Get the best sentence according to the sentences dictionary
	var bestSentence = ""
	var maxValue float32 = -1
	for _, s := range sentences {
		var trimmedSentence = formatSentence(s)
		if trimmedSentence != "" && sentencesDictionary[trimmedSentence] > maxValue {
			maxValue = sentencesDictionary[trimmedSentence]
			bestSentence = s
		}
	}

	return bestSentence
}

// GetSummary builds the summary from the given content text
func GetSummary(content string) string {
	// Build the sentences dictionary
	var sentencesDictionary = getSentencesRanks(content)

	// Split the content into paragraphs
	var paragraphs = getContentParagraphs(content)

	// Add the title
	var summary = []string{}

	// Add the best sentence from each paragraph
	for _, paragraph := range paragraphs {
		var currentBestSentence = getBestSentence(paragraph, sentencesDictionary)
		var sentence = strings.TrimSpace(currentBestSentence)
		if sentence != "" {
			summary = append(summary, sentence)
		}
	}

	if len(summary) == 0 && len(sentencesDictionary) == len(paragraphs) && len(sentencesDictionary) > 1 {
		// Then we have one sentence per paragraph
		// This way we combine all sentences in one paragraph
		var newContent = strings.Replace(content, "\n\n", " ", -1)
		var result = GetSummary(newContent)
		return result
	}

	var result = strings.Join(summary, "\n")
	return result
}
