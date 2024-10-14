package helper

import (
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"strings"
)

func RemoveOmit(tx *gorm.DB, columns ...string) {
	var removeColumnsSlice []string
	removeColumnsMap := make(map[string]bool)
	if len(columns) == 1 && strings.ContainsRune(columns[0], ',') {
		removeColumnsSlice = strings.FieldsFunc(columns[0], utils.IsValidDBNameChar)
	} else {
		removeColumnsSlice = columns
	}

	// Transform to map for faster indexing
	for _, columnName := range removeColumnsSlice {
		removeColumnsMap[columnName] = false
	}

	var newOmits []string
	// create the new slice
	for _, omitField := range tx.Statement.Omits {
		// If it's not the same then add to the new slice

		if _, shouldBeRemoved := removeColumnsMap[omitField]; !shouldBeRemoved {
			newOmits = append(newOmits, omitField)
		}
	}
	// set back the new created slice
	tx.Statement.Omits = newOmits
}
