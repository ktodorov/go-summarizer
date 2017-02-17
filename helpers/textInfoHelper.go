package helpers

import (
	"bytes"
	"strconv"
)

// GetSummaryInfo Returns summary information statistics for the original text and summarized text
func GetSummaryInfo(originalText string, summarizedText string, imagesCount int) string {
	// Print the ratio between the summary length and the original length
	var summaryInfo bytes.Buffer

	appendLine(&summaryInfo, "Summary info:")

	var originalTextLength = float64(len(originalText))
	var summarizedTextLength = float64(len(summarizedText))
	var ratio = (100 - (100 * (summarizedTextLength / originalTextLength)))

	var originalLengthString = strconv.FormatFloat(originalTextLength, 'f', -1, 64)
	var summarizedLengthString = strconv.FormatFloat(summarizedTextLength, 'f', -1, 64)
	var ratioString = strconv.FormatFloat(ratio, 'f', 2, 64)

	appendLine(&summaryInfo, " - Original length: ", originalLengthString, " symbols")
	appendLine(&summaryInfo, " - Summary length:  ", summarizedLengthString, " symbols")
	appendLine(&summaryInfo, " - Summary ratio:   ", ratioString, "%")
	if imagesCount > 0 {
		appendLine(&summaryInfo, " - Images found:    ", strconv.Itoa(imagesCount))
	}

	var summaryInfoString = summaryInfo.String()
	return summaryInfoString
}

func appendLine(mainString *bytes.Buffer, stringsToAppend ...string) {
	for _, stringToAppend := range stringsToAppend {
		mainString.WriteString(stringToAppend)
	}

	mainString.WriteString("\n")
}
