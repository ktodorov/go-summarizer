package helpers

import (
	"fmt"
)

func logError(err error) {
	fmt.Printf("Error occurred: {%s}\n", err.Error())
}
