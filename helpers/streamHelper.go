package helpers

import (
	"fmt"
	"io"
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
