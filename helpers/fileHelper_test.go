package helpers

import "testing"

func TestProgramRootPathGetter(t *testing.T) {
	var path, err = getProgramRootPath()
	if err != nil {
		t.Error("Didn't expect error but received: ", err.Error())
	}

	if path == "" {
		t.Error("Expected root path but received empty path")
	}
}

func TestRandomFileNameGenerating(t *testing.T) {
	var rootPath, err = getProgramRootPath()
	if err != nil {
		t.Error("Didn't expect error but received: ", err.Error())
	}

	var filename = generateRandomFileName(rootPath, "test")
	if filename == "" {
		t.Error("Expected filename but received empty string")
	}
}

func TestRandomFileNameGeneratingExistence(t *testing.T) {
	var rootPath, err = getProgramRootPath()
	if err != nil {
		t.Error("Didn't expect error but received: ", err.Error())
	}

	var filename = generateRandomFileName(rootPath, "test")
	if fileExists(filename) {
		t.Error("Expected generated file not to exist, but it does")
	}
}

func TestFileTypeGetter(t *testing.T) {
	var rootPath, err = getProgramRootPath()
	if err != nil {
		t.Error("Didn't expect error but received: ", err.Error())
	}

	var filename = generateRandomFileName(rootPath, "pdf")
	var fileType = getFileType(filename)
	if fileType != "pdf" {
		t.Error("Expected 'pdf' file extension but received: ", fileType)
	}
}

func TestFileExtensionFromURLGetter(t *testing.T) {
	var testURL = "testurl/testimage.png"
	var fileExtension, found = getFileExtensionFromURL(testURL)
	if !found {
		t.Error("Expected to find the file type but it didn't")
	}

	if fileExtension != "png" {
		t.Error("Expected 'png' file type but received: ", fileExtension)
	}
}
