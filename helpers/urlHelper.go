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

func ExtractMainTextFromURL(url string) (string, error) {
	var htmlString, err = getHTMLFromURL(url)
	if err != nil {
		logError(err)
		return "", err
	}

	textFromHTML, err := getTextFromHTML(htmlString)
	if err != nil {
		logError(err)
		return "", err
	}

	return textFromHTML, nil
}

func IsURL(text string) bool {
	var urlRegex, err = regexp.Compile("[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)")
	if err != nil {
		return false
	}

	var isURL = urlRegex.MatchString(text)
	return isURL
}
