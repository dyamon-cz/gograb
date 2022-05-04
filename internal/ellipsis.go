/*
Package internal
Copyright © 2022 Pavel Sidlo <github.com/dyamon-cz>
*/
package internal

import "unicode"

// ellipsis https://stackoverflow.com/a/59955803
func ellipsis(str string, max int) string {
	lastSpaceIx := -1
	length := 0
	for i, r := range str {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		length++
		if length >= max {
			if lastSpaceIx != -1 {
				return str[:lastSpaceIx] + " ..."
			}
			// If here, string is longer than max, but has no spaces
		}
	}
	// If here, string is shorter than max

	return str
}
