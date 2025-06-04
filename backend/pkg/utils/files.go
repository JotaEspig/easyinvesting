package utils

import (
	"fmt"
	"io"
	"os"
)

func ReadFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		pwd, _ := os.Getwd()
		errorStr := fmt.Sprintf(
			"ReadFile: unable to find file: '%s' at '%s'\n",
			filename,
			pwd,
		)
		HandleErr(err, errorStr)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(content)
}
