package validations

import (
	"bufio"
	"log"
	"os"
)

func ValidateState(state string) bool {
	file, err := os.Open("web/data/states.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, line := range txtlines {
		if line == state {
			return true
		}
	}
	return false
}
