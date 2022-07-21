package songio

import "unicode"

func isEmptyOrWhitespace(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}
