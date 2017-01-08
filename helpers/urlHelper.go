package helpers

import "net/http"

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
