package telegram

import (
	"math/rand/v2"
	"regexp"
	"strings"
)

const botUsername = "maznevich_bot"
const replyChance = 0.015

var travlyaRegex = regexp.MustCompile("травл(?:я|и|ю|ей|е|ям|ями|ях|явш|емый|емую|емого)?")

func needReply(text string) bool {
	return travlyaRegex.MatchString(strings.ToLower(text)) ||
		strings.Contains(strings.ToLower(text), botUsername) ||
		rand.Float32() < replyChance
}
