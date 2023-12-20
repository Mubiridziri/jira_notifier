package helpers

import "bytes"

func InStringArray(array []string, need string) int {
	for index, el := range array {
		if el == need {
			return index
		}
	}

	return -1
}

func GetLastFoundSymbolIndex(str, need string) int {
	runes := bytes.Runes([]byte(str))
	foundIndex := -1

	for index, symbolRune := range runes {
		symbol := string(symbolRune)
		if symbol == need {
			foundIndex = index
		}
	}
	return foundIndex
}
