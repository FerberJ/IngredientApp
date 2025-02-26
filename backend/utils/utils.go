package utils

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
