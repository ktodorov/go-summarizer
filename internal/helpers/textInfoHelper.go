package helpers

import (
	"bytes"
	"strconv"
)

// GetSummaryInfo Returns summary information statistics for the original text and summarized text
func GetSummaryInfo(originalText string, summarizedText string) string {
	// Print the ratio between the summary length and the original length
	var summaryInfo bytes.Buffer

	appendLine(&summaryInfo, "Summary info:")

	var originalTextLength = len(originalText)
	var summarizedTextLength = len(summarizedText)
	var ratio = (100 - (100 * (summarizedTextLength / originalTextLength)))

	appendLine(&summaryInfo, " - Original length: ", strconv.Itoa(originalTextLength))
	appendLine(&summaryInfo, " - Summary length:  ", strconv.Itoa(summarizedTextLength))
	appendLine(&summaryInfo, " - Summary ratio:   ", strconv.Itoa(ratio))

	var summaryInfoString = summaryInfo.String()
	return summaryInfoString
}

func appendLine(mainString *bytes.Buffer, stringsToAppend ...string) {
	for _, stringToAppend := range stringsToAppend {
		mainString.WriteString(stringToAppend)
	}

	mainString.WriteString("\n")
}
