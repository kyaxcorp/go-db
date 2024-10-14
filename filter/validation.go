package filter

import "strings"

//func secureFieldName(fieldName string) string {
//	return "`" + fieldName + "`"
//}

func CountWords(s string) int {
	return len(strings.Fields(s))
}

func validateFieldName(fieldName string) bool {
	if CountWords(fieldName) != 1 {
		return false
	}

	// Remove symbols, check symbols like:
	/*
		; ( ) / * \ | ' " @ # $ % & - = +
	*/

	if strings.ContainsAny(fieldName, "?!^{}[]`,:<>;()/*\\|'\"@#$%&-=+ ") {
		return false
	}

	return true
}

func validateFieldNameAndPanic(fieldName string) {
	if !validateFieldName(fieldName) {
		panic("invalid field format -> " + fieldName)
	}
}
