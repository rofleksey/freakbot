package telegram

import (
	"regexp"
	"strings"
)

const thePhrase = "Спасибо за травлю в интернете!"

var travlyaRegex = regexp.MustCompile("травл(?:я|и|ю|ей|е|ям|ями|ях|явш|емый|емую|емого)?")

func containsBullying(text string) bool {
	return travlyaRegex.MatchString(strings.ToLower(text))
}
