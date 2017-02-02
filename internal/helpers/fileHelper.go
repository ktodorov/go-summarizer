package helpers

import (
	"errors"
	"io/ioutil"
	"regexp"
	"runtime"

	"fmt"

	"path/filepath"

	"github.com/signintech/gopdf"
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
func StoreTextToFile(path string, text string, images []string) (bool, error) {
	var fileType = getFileType(path)
	var textAsBytes = []byte(text)
	var result = false
	var err error

	if fileType == "txt" {
		result, err = saveToTextFile(path, textAsBytes)
	} else if fileType == "pdf" {
		result, err = saveToPDFFile(path, textAsBytes, images)
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

func saveToPDFFile(path string, text []byte, images []string) (bool, error) {
	pdf := gopdf.GoPdf{}
	var goPdfRect = gopdf.Rect{W: 595.28, H: 841.89}
	pdf.Start(gopdf.Config{PageSize: goPdfRect}) //595.28, 841.89 = A4
	pdf.AddPage()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))
	abspath, err := filepath.Abs(basepath + "/fonts/OpenSans-Regular.ttf")
	if err != nil {
		return false, err
	}

	fmt.Println("abspath: ", abspath)

	err = pdf.AddTTFFont("OpenSans-Regular", abspath)

	if err != nil {
		return false, err
	}

	err = pdf.SetFont("OpenSans-Regular", "", 14)
	if err != nil {
		return false, err
	}

	for _, image := range images {
		imagePath, err := saveImageFromURL(image)
		if err != nil {
			return false, err
		}

		pdf.Image(imagePath, 0, 0, &goPdfRect) //print image
		pdf.AddPage()
	}

	pdf.SetX(50) //move current location
	pdf.SetY(50)
	pdf.Cell(&goPdfRect, string(text)) //print text
	pdf.WritePdf(path)

	return true, nil
}
