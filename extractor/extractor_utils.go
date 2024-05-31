package extractor

import (
	"regexp"
	"strings"
)

func ExtractDate(fileName string) {
	// date_pattern, _ := regexp.Compile(fromRegex)
	date_pattern, _ := regexp.Compile(`(\d{4})(\d{2})(\d{2})\.\w+`)
	date := strings.Split(fileName, "_")
	date_match := date_pattern.FindStringSubmatch(date[1])

	ParsedData.SetYearAndMonth(date_match[1], date_match[2])
}