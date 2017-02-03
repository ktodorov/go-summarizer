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
)

var supportedFileTypes = []string{"txt", "pdf"}

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
	var fullpath = path + "/" + string(b) + "." + extension
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

func saveToPDFFile(path string, text []byte, imageURLs []string) (bool, error) {
	pdf := gopdf.GoPdf{}
	var pageSizeHeight = 841.89
	var pageSizeWidth = 595.28
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageSizeWidth, H: pageSizeHeight}}) //595.28, 841.89 = A4
	pdf.AddPage()

	abspath, err := getProgramRootPath()
	if err != nil {
		return false, err
	}

	var fontsPath = abspath + "/fonts/OpenSans-Regular.ttf"

	err = pdf.AddTTFFont("OpenSans-Regular", fontsPath)

	if err != nil {
		return false, err
	}

	err = pdf.SetFont("OpenSans-Regular", "", 12)
	if err != nil {
		return false, err
	}

	var imagePaths = []string{}
	var heightUsed = 0.0
	for _, imageURL := range imageURLs {
		var imageExtension, isImage = getFileExtensionFromURL(imageURL)
		if !isImage {
			continue
		}

		var imagePath = generateRandomFileName(abspath, imageExtension)
		err := saveImageFromURL(imageURL, imagePath)
		if err != nil {
			return false, err
		}

		// Check if image will leave the page.
		// If thats the case, add new page and start from 0 there
		imageWidth, imageHeight, err := getImageDimension(imagePath)
		if err != nil {
			return false, err
		}

		var floatImageHeight = float64(imageHeight)
		var floatImageWidth = float64(imageWidth)
		var imageProp = imageWidth / int(pageSizeWidth)
		floatImageWidth = floatImageWidth / float64(imageProp+1)
		floatImageHeight = floatImageHeight / float64(imageProp+1)

		if heightUsed+floatImageHeight > pageSizeHeight {
			pdf.AddPage()
			heightUsed = 0
		}

		pdf.Image(imagePath, 0, heightUsed, &gopdf.Rect{H: floatImageHeight, W: floatImageWidth}) //print image
		heightUsed += floatImageHeight

		// save image file paths in order to delete them later
		imagePaths = append(imagePaths, imagePath)
	}

	pdf.SetX(5) //move current location
	pdf.SetY(heightUsed + 5)
	pdf.Cell(nil, string(text)) //print text
	pdf.WritePdf(path)

	deleteFiles(imagePaths) // delete temporary created images

	return true, nil
}

func deleteFiles(filePaths []string) {

	for _, filePath := range filePaths {
		if fileExists(filePath) {
			os.Remove(filePath)
		}
	}
}
