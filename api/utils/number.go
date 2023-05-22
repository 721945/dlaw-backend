package utils

import "strings"

func ConvertThaiNumToEngNum(text string) string {
	newText := text
	mapping := map[string]string{
		"๐": "0",
		"๑": "1",
		"๒": "2",
		"๓": "3",
		"๔": "4",
		"๕": "5",
		"๖": "6",
		"๗": "7",
		"๘": "8",
		"๙": "9",
	}

	for thai, normal := range mapping {
		newText = strings.Replace(newText, thai, normal, -1)
	}

	return newText
}
