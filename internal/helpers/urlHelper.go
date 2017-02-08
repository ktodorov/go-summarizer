package helpers

import "net/http"
import "regexp"

func getHTMLFromURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	htmlBytes, err := readFromReader(response.Body)
	if err != nil {
		return "", err
	}

	var htmlString = string(htmlBytes)
	return htmlString, nil
}

// ExtractMainInfoFromURL searches the main content from the given url and returns the text and images
func ExtractMainInfoFromURL(url string) (string, string, []string, error) {
	var htmlString, err = getHTMLFromURL(url)
	if err != nil {
		logError(err)
		return "", "", nil, err
	}

	titleFromHTML, textFromHTML, imagesFromHTML, err := getMainInfoFromHTML(htmlString)
	if err != nil {
		logError(err)
		return "", "", nil, err
	}

	return titleFromHTML, textFromHTML, imagesFromHTML, nil
}

// IsURL checks if the given text is a website url address
func IsURL(text string) bool {
	var urlRegex, err = regexp.Compile("[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)")
	if err != nil {
		return false
	}

	var isURL = urlRegex.MatchString(text)
	return isURL
}
