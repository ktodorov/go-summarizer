package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func IsPositiveAnswer(answer string) bool {
	var isPositive = (answer == "y" || answer == "Y")
	return isPositive
}

func ReadInputFromUser(message string) string {
	reader := bufio.NewReader(os.Stdin)
	if message != "" {
		fmt.Println(message)
	}

	text, _ := reader.ReadString('\n')
	var trimmedText = strings.TrimSpace(text)

	// Listen while the user enter something instead of pressing enter or space
	for trimmedText == "" {
		trimmedText = ReadInputFromUser("")
	}

	return trimmedText
}
