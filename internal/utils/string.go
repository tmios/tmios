package utils

import "strings"

func Overlay(str string, overlay string, start int, end int) (overlayed string) {
	r := []rune(str)
	l := len([]rune(r))

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	overlayed += overlay
	overlayed += string(r[end:])
	return overlayed
}

func ReplaceDup(str string) string {
	odsMessage := strings.ReplaceAll(str, "Duplicati", "NDP")
	adsMessage := strings.ReplaceAll(odsMessage, "duplicati", "ndp")
	return adsMessage
}
