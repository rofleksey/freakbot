package telegram

import (
	"regexp"
	"strings"
)

const botUsername = "maznevich_bot"

var travlyaRegex = regexp.MustCompile("травл(?:я|и|ю|ей|е|ям|ями|ях|явш|емый|емую|емого)?")

func needReply(text string) bool {
	return travlyaRegex.MatchString(strings.ToLower(text)) || strings.Contains(strings.ToLower(text), botUsername)
}
