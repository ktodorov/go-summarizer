package helpers

import (
	"regexp"
	"strings"
)

// Naive method for splitting a text into sentences
func getContentSentences(content string) []string {
	content = strings.Replace(content, "\n", ". ", -1)
	var sentences = strings.Split(content, ". ")
	return sentences
}

// Naive method for splitting a text into paragraphs
func getContentParagraphs(content string) []string {
	var paragraphs = strings.Split(content, "\n\n")
	return paragraphs
}

// Caculate the intersection between 2 sentences
func sentencesIntersectedWordsCount(sent1 string, sent2 string) float32 {
	// split the sentence into words/tokens
	var words1 = strings.Split(sent1, " ")
	var words2 = strings.Split(sent2, " ")

	// If there is not intersection, just return 0
	if len(words1) == 0 && len(words2) == 0 {
		return 0
	}

	var intersectionCount = 0

	var matchedWords = make(map[string]bool)
	for _, word1 := range words1 {
		if _, exists := matchedWords[word1]; exists {
			continue
		}

		for _, word2 := range words2 {
			if word1 != word2 {
				continue
			}

			if _, exists := matchedWords[word2]; !exists {
				intersectionCount++
				matchedWords[word2] = true
			}
		}
	}

	for _, word2 := range words2 {
		if _, exists := matchedWords[word2]; exists {
			continue
		}

		for _, word1 := range words1 {
			if word1 != word2 {
				continue
			}

			if _, exists := matchedWords[word2]; !exists {
				intersectionCount++
				matchedWords[word2] = true
			}
		}
	}

	// We normalize the result by the average number of words
	var numerator = float32(len(matchedWords))
	var denominator = float32((len(words1) + len(words2)) / 2)
	var result = numerator / denominator
	return result
}

// Format a sentence - remove all non-alphbetic chars from the sentence
// We'll use the formatted sentence as a key in our sentences dictionary
func formatSentence(sentence string) string {
	var regex, err = regexp.Compile("\\W+")
	if err != nil {
		return ""
	}
	var replacedSentence = regex.ReplaceAllString(sentence, "")

	// sentence = re.sub(r'\W+', '', sentence)
	return replacedSentence
}

// Convert the content into a dictionary <K, V>
// k = The formatted sentence
// V = The rank of the sentence
func GetSentencesRanks(content string) map[string]float32 {

	// Split the content into sentences
	var sentences = getContentSentences(content)

	// Calculate the intersection of every two sentences
	var sentencesCount = len(sentences)
	var values = [][]float32{}

	for i := 0; i < sentencesCount; i++ {
		values = append(values, []float32{})
		for j := 0; j < sentencesCount; j++ {
			values[i] = append(values[i], sentencesIntersectedWordsCount(sentences[i], sentences[j]))
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
	var maxValue float32
	for _, s := range sentences {
		var trimmedSentence = formatSentence(s)
		if trimmedSentence != "" && sentencesDictionary[trimmedSentence] > maxValue {
			maxValue = sentencesDictionary[trimmedSentence]
			bestSentence = s
		}
	}

	return bestSentence
}

// Build the summary
func GetSummary(content string, sentencesDictionary map[string]float32) string {

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

	var result = strings.Join(summary, "\n")
	return result
}
