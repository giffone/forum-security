package sqlitedb

func minimalIDToShow(max, sep int) int {
	return max - sep
}

func maximumIDToShow(max, max2 int) int {
	if max2 == 0 {
		return max
	}
	return max2
}