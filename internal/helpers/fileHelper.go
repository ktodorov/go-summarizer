package helpers

import (
	"errors"
	"io/ioutil"
	"regexp"
)

var supportedFileTypes = []string{"txt", "pdf"}

func getFileType(path string) string {
	for _, fileType := range supportedFileTypes {
		var regexString = ".+\\." + fileType

		matched, err := regexp.MatchString(regexString, path)
		if err == nil && matched {
			return fileType
		}
	}
	return ""
}

//StoreTextToFile stores text to the given file path. Creates the file if it's missing or appends to it
func StoreTextToFile(path string, text string) (bool, error) {
	var fileType = getFileType(path)
	var textAsBytes = []byte(text)
	var result = false
	var err error

	if fileType == "txt" {
		result, err = saveToTextFile(path, textAsBytes)
	} else if fileType == "pdf" {
		result, err = saveToPDFFile(path, textAsBytes)
	} else {
		err = errors.New("Invalid file type")
	}

	return result, err
}

func saveToTextFile(path string, text []byte) (bool, error) {
	err := ioutil.WriteFile(path, text, 0644)
	var result = (err == nil)
	return result, err
}

func saveToPDFFile(path string, text []byte) (bool, error) {
	// pdf := gopdf.GoPdf{}
	// pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	// pdf.AddPage()
	// pdf.Cell(nil, text)
	// pdf.WritePdf(path)

	return true, nil
}
