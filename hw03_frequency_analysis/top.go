package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var /* const */ splitter = regexp.MustCompile(`\s+(-\s+)*`)

//Regex to remove all special characters from text and avoid apostrophe problems
// var /* const */ replacer = regexp.MustCompile(`(?<![A-Za-z])[']|['](?![A-Za-z])|[,.!:"]`)
var /* const */ replacer = strings.NewReplacer(",", "", ".", "", "!", "", ":", "", "'", "", "\"", "")

func Top10(input string) []string {
	//Return nil if string is empty
	if input == "" {
		return nil
	}

	// Generate base structs
	words := make(map[string]int)
	// Collect all words
	for _, el := range splitter.Split(replacer.Replace(input), -1) {
		if el != "" {
			words[strings.ToLower(el)]++
		}
	}

	fmt.Println(words)

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
	if len(buffer) > 10 {
		return buffer[0:10]
	} else {
		return buffer
	}

}
