package helpers

import (
	"fmt"
)

func logError(err error) {
	fmt.Printf("Error occured: {%s}\n", err.Error())
}
