package helpers

import (
	"errors"
	"image"
	"io/ioutil"
	"math/rand"
	"regexp"
	"runtime"

	"path/filepath"

	"os"

	"strings"

	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fontmaker/core"
)

var supportedFileTypes = []string{"txt", "pdf"}

// io.TempFile
func getProgramRootPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))
	abspath, err := filepath.Abs(basepath)
	return abspath, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomFileName(path string, extension string) string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	var fullpath = path + "\\" + string(b) + "." + extension
	if !fileExists(fullpath) {
		return fullpath
	}

	return generateRandomFileName(path, extension)
}

func fileExists(filePath string) bool {
	var _, err = os.Stat(filePath)

	if err == nil {
		return true
	}

	return false
}

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

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

func getFileExtensionFromURL(url string) (string, bool) {
	var splitURL = strings.Split(url, ".")
	var lastElement = splitURL[len(splitURL)-1]
	if strings.Contains(lastElement, "jpg") {
		return "jpg", true
	}
	if strings.Contains(lastElement, "png") {
		return "png", true
	}
	if strings.Contains(lastElement, "jpeg") {
		return "jpeg", true
	}

	return "", false
}

//StoreTextToFile stores text to the given file path. Creates the file if it's missing or appends to it
func StoreTextToFile(path string, title string, text string, images []string) (bool, error) {
	var fileType = getFileType(path)
	var titleAsBytes = []byte(title)
	var textAsBytes = []byte(text)
	var result = false
	var err error

	if fileType == "txt" {
		result, err = saveToTextFile(path, titleAsBytes, textAsBytes)
	} else if fileType == "pdf" {
		result, err = saveToPDFFile(path, titleAsBytes, textAsBytes, images)
	} else {
		err = errors.New("Invalid file type")
	}

	return result, err
}

func saveToTextFile(path string, title []byte, text []byte) (bool, error) {
	err := ioutil.WriteFile(path, title, 0644)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(path, text, 0644)
	var result = (err == nil)
	return result, err
}

func saveToPDFFile(path string, title []byte, text []byte, imageURLs []string) (bool, error) {
	pdf := gopdf.GoPdf{}
	var pageSizeHeight = 841.89
	var pageSizeWidth = 595.28
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageSizeWidth, H: pageSizeHeight}}) //595.28, 841.89 = A4
	pdf.AddPage()

	abspath, err := getProgramRootPath()
	if err != nil {
		return false, err
	}

	var fontPath = abspath + "/fonts/OpenSans-Regular.ttf"
	var fontSize = 12
	err = pdf.AddTTFFont("OpenSans-Regular", fontPath)

	if err != nil {
		return false, err
	}

	var imagePaths = []string{}
	var heightUsed = 0.0

	fontSize = 16
	heightUsed, err = writeTitleToPDF(&pdf, string(title), heightUsed, fontSize, fontPath, pageSizeWidth, pageSizeHeight)
	if err != nil {
		return false, err
	}

	for _, imageURL := range imageURLs {
		var imageExtension, isImage = getFileExtensionFromURL(imageURL)
		if !isImage {
			continue
		}

		var imagePath = generateRandomFileName(abspath, imageExtension)
		err := saveImageFromURL(imageURL, imagePath)
		if err != nil {
			continue
		}

		// Check if image will leave the page.
		// If thats the case, add new page and start from 0 there
		imageWidth, imageHeight, err := getImageDimension(imagePath)
		if err != nil {
			continue
		}

		// Print the image only if its not some kind of icon below 50 px from both sides
		if imageWidth+imageHeight > 50 {
			var floatImageHeight = float64(imageHeight)
			var floatImageWidth = float64(imageWidth)
			var imageProp = imageWidth / int(pageSizeWidth)
			floatImageWidth = floatImageWidth / float64(imageProp+1)
			floatImageHeight = floatImageHeight / float64(imageProp+1)

			if heightUsed+floatImageHeight > pageSizeHeight {
				pdf.AddPage()
				heightUsed = 0
			}

			if len(imagePaths) < 2 {
				pdf.Image(imagePath, 0, heightUsed, &gopdf.Rect{H: floatImageHeight, W: floatImageWidth}) //print image
				heightUsed += floatImageHeight
			}
		}

		// save image file paths in order to delete them later
		imagePaths = append(imagePaths, imagePath)
	}
	heightUsed += 5

	fontSize = 12
	_, err = writeBodyToPDF(&pdf, string(text), heightUsed, fontSize, fontPath, pageSizeWidth, pageSizeHeight)
	if err != nil {
		return false, err
	}

	pdf.WritePdf(path)
	deleteFiles(imagePaths) // delete temporary created images

	return true, nil
}

func writeTitleToPDF(pdf *gopdf.GoPdf, title string, heightUsed float64, fontSize int, fontPath string, pageSizeWidth float64, pageSizeHeight float64) (float64, error) {
	if title == "" {
		return heightUsed, nil
	}

	var err = pdf.SetFont("OpenSans-Regular", "", fontSize)
	if err != nil {
		return heightUsed, err
	}

	titleHeight, err := calculateTextHeight(fontPath, fontSize)
	if err != nil {
		return heightUsed, err
	}

	heightUsed, _ = writeTextToPDF(pdf, 5, 0, title, pageSizeWidth, pageSizeHeight, titleHeight)
	heightUsed += 10 // add padding between header and rest of body

	return heightUsed, nil
}

func writeBodyToPDF(pdf *gopdf.GoPdf, bodyText string, heightUsed float64, fontSize int, fontPath string, pageSizeWidth float64, pageSizeHeight float64) (float64, error) {
	var err = pdf.SetFont("OpenSans-Regular", "", fontSize)
	if err != nil {
		return heightUsed, err
	}

	textHeight, err := calculateTextHeight(fontPath, fontSize)
	if err != nil {
		return heightUsed, err
	}
	textHeight += 4

	heightUsed, err = writeTextToPDF(pdf, 5, heightUsed, string(bodyText), pageSizeWidth, pageSizeHeight, textHeight)
	if err != nil {
		return heightUsed, err
	}

	return heightUsed, nil
}

func calculateTextHeight(fontPath string, fontSize int) (float64, error) {
	var parser core.TTFParser
	var err = parser.Parse(fontPath)
	if err != nil {
		return 0, err
	}

	//Measure Height
	//get  CapHeight (https://en.wikipedia.org/wiki/Cap_height)
	cap := float64(float64(parser.CapHeight()) * 1000.00 / float64(parser.UnitsPerEm()))
	//convert
	realHeight := cap * (float64(fontSize) / 1000.0)

	return realHeight, nil
}

func writeTextToPDF(pdf *gopdf.GoPdf, startX float64, startY float64, text string, pageSizeWidth float64, pageSizeHeight float64, textHeight float64) (float64, error) {
	var textWords = strings.Split(text, " ")
	var currentLineText = ""
	var currentLine = 1.0

	for _, word := range textWords {
		var wordToUse = word
		var newLineWord = ""
		// If there is new line in the current word, then we have something like word1\nword2
		// Then we split the line after the first word and write the second on the new line
		if strings.Contains(word, "\n") {
			var wordsByLines = strings.Split(word, "\n")
			wordToUse = wordsByLines[0]
			if len(wordsByLines) > 1 {
				newLineWord = wordsByLines[1]
			}
		}

		var tempLineText = currentLineText + " " + wordToUse
		var textWidth, err = pdf.MeasureTextWidth(tempLineText)
		if err != nil {
			return 0, err
		}

		// This means that if we add the current word,
		// the text on the current line will move out of the page,
		// so we add the current word to a new line and write the current line on the pdf
		if textWidth > pageSizeWidth-15 {
			currentLine, startY = writeTextLineToPDF(pdf, startX, startY, textHeight, currentLine, currentLineText, pageSizeHeight, false)
			currentLineText = wordToUse
		} else {
			currentLineText = tempLineText
		}

		// This means that we split the current word and after the first, we must move to new line
		if newLineWord != "" {
			currentLine, startY = writeTextLineToPDF(pdf, startX, startY, textHeight, currentLine, currentLineText, pageSizeHeight, true)
			currentLineText = newLineWord
		}
	}

	var heightUsed = startY + (textHeight * currentLine)
	pdf.SetX(startX)
	pdf.SetY(heightUsed)
	pdf.Cell(nil, currentLineText)

	heightUsed += textHeight
	pdf.Br(textHeight)

	return heightUsed, nil
}

func writeTextLineToPDF(pdf *gopdf.GoPdf, startX float64, startY float64, textHeight float64, currentLine float64, text string, pageSizeHeight float64, newParagraph bool) (newCurrentLine float64, newStartY float64) {
	pdf.SetX(startX)
	var currentHeight = startY + (textHeight * currentLine)
	if currentHeight+textHeight > pageSizeHeight {
		pdf.AddPage()
		currentHeight = 5
		currentLine = 1
		startY = 5
	} else {
		currentLine++
		if newParagraph {
			currentLine++
		}
	}
	pdf.SetY(currentHeight)
	pdf.Cell(nil, text)

	if newParagraph {
		pdf.Br(textHeight)
	}

	return currentLine, startY
}

func deleteFiles(filePaths []string) {
	for _, filePath := range filePaths {
		if fileExists(filePath) {
			os.Remove(filePath)
		}
	}
}
