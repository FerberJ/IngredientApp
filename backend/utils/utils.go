package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

func FormatDuration(input string) (string, error) {
	// Check if the input is already in HH:mm format
	timeRegex := regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$`)
	if timeRegex.MatchString(input) {
		// If it's already in correct format, ensure it has leading zeros
		parts := strings.Split(input, ":")
		hour, _ := strconv.Atoi(parts[0])
		return fmt.Sprintf("%02d:%s", hour, parts[1]), nil
	}

	// Check if it's in ISO 8601 duration format (PT...)
	if !strings.HasPrefix(input, "PT") {
		return "", fmt.Errorf("invalid format: input must be in HH:mm or PT... format")
	}

	// Remove the PT prefix
	duration := input[2:]

	// Initialize hours and minutes
	hours := 0
	minutes := 0

	// Extract hours if present
	hourIndex := strings.Index(duration, "H")
	if hourIndex != -1 {
		hourStr := duration[:hourIndex]
		h, err := strconv.Atoi(hourStr)
		if err != nil {
			return "", fmt.Errorf("invalid hour format: %v", err)
		}
		hours = h
		duration = duration[hourIndex+1:]
	}

	// Extract minutes if present
	minuteIndex := strings.Index(duration, "M")
	if minuteIndex != -1 {
		minuteStr := duration[:minuteIndex]
		m, err := strconv.Atoi(minuteStr)
		if err != nil {
			return "", fmt.Errorf("invalid minute format: %v", err)
		}
		minutes = m
	}

	// Convert total minutes to hours and minutes
	totalMinutes := hours*60 + minutes
	hours = totalMinutes / 60
	minutes = totalMinutes % 60

	// Format as HH:mm
	return fmt.Sprintf("%02d:%02d", hours, minutes), nil
}

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
