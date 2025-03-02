package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func ContainsAllKeywords(keywords, valSlice []string) bool {
	keywordMap := make(map[string]bool)

	for _, tag := range keywords {
		keywordMap[tag] = true
	}

	for _, val := range valSlice {
		if !keywordMap[val] {
			return false
		}
	}

	return true
}

func RemoveString(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func PtToString(pt string) string {
	ret := ""

	hourMatchString := "[0-9]+H"
	minuteMatchString := "[0-9]+M"

	hourRe := regexp.MustCompile(hourMatchString)
	minuteRe := regexp.MustCompile(minuteMatchString)

	hourStr := hourRe.FindString(pt)
	minuteStr := minuteRe.FindString(pt)

	if hourStr != "" {
		hourNum, _ := strconv.Atoi(strings.TrimSuffix(hourStr, "H"))

		hourUnit := "Hours"
		if hourNum == 1 {
			hourUnit = "Hour"
		}

		hourStr = strings.ReplaceAll(hourStr, "H", " "+hourUnit)
		ret += hourStr
	}

	if minuteStr != "" {
		minuteNum, _ := strconv.Atoi(strings.TrimSuffix(minuteStr, "M"))

		minuteUnit := "Minutes"
		if minuteNum == 1 {
			minuteUnit = "Minute"
		}

		minuteStr = strings.ReplaceAll(minuteStr, "M", " "+minuteUnit)
		ret += " " + minuteStr
	}

	return ret
}
