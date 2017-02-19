package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func readFromReader(reader io.Reader) ([]byte, error) {
	var resultBytes []byte
	b := make([]byte, 1024)
	var err error
	var bytesRead int

	for err != io.EOF {
		bytesRead, err = reader.Read(b)
		if err != nil && bytesRead == 0 {
			fmt.Println("error occurred: ", err.Error(), "\nbytes read: ", bytesRead)
			return nil, err
		}

		resultBytes = append(resultBytes, b...)
	}

	return resultBytes, nil
}

func saveImageFromURL(imageURL string, imagePath string) error {
	// don't worry about errors
	response, e := http.Get(imageURL)

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	if e != nil {
		return e
	}

	//open a file for writing
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
