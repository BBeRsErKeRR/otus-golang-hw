package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	if input != "" {
		// Generate base structs
		words := make(map[string]int)
		re := regexp.MustCompile(`\s+(-\s+)*`)
		replace := regexp.MustCompile(`[,.!:"']+`)
		// Collect all words
		for _, el := range re.Split(replace.ReplaceAllString(input, ""), -1) {
			if el != "" {
				words[strings.ToLower(el)]++
			}
		}
		// Check if words more then 10
		if len(words) >= 10 {
			// Generate slice of words
			buffer := make([]string, 0, len(words))
			for k := range words {
				buffer = append(buffer, k)
			}
			// Sorting
			sort.Slice(buffer, func(i, j int) bool {
				return words[buffer[i]] > words[buffer[j]] ||
					(words[buffer[i]] == words[buffer[j]] && buffer[i] < buffer[j])
			})
			// return only first 10 words
			res := buffer[0:10]
			fmt.Println(res)
			return res
		}
	}
	return nil
}
