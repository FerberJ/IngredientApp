package utils

import (
	"fmt"
	"time"
)

func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func ContainsAllKeywords(keywords, selectedFilter []string) bool {
	keywordMap := make(map[string]bool)

	for _, tag := range keywords {
		keywordMap[tag] = true
	}

	for _, val := range selectedFilter {
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

func GetTotalTime(t1Str, t2Str string) string {
	t1, _ := time.Parse("15:04", t1Str)
	t2, _ := time.Parse("15:04", t2Str)

	// Convert to duration since midnight
	d1 := time.Duration(t1.Hour())*time.Hour + time.Duration(t1.Minute())*time.Minute
	d2 := time.Duration(t2.Hour())*time.Hour + time.Duration(t2.Minute())*time.Minute

	// Sum the durations
	sum := d1 + d2

	// Format the result
	hours := int(sum.Hours())
	minutes := int(sum.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
