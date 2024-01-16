package app

import (
	"regexp"
	"strings"
)

func Check(result string, regex string) (bool, error) {

	if strings.HasPrefix(result, "Merge branch") {
		return true, nil
	}

	re := regexp.MustCompile(regex)
	for _, commitMessage := range strings.Split(result, "\n") {
		if re.MatchString(commitMessage) {
			return true, nil
		}
	}
	return false, nil
}
