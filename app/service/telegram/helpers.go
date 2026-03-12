package telegram

import (
	"math/rand/v2"
	"regexp"
	"strings"
)

const botUsername = "maznevich_bot"
const replyChance = 0.02

var phrases = []string{
	"надо отыграться",
	"стрим",
	"бот",
}

var travlyaRegex = regexp.MustCompile("травл(?:я|и|ю|ей|е|ям|ями|ях|явш|емый|емую|емого)?")

func needReply(text string) bool {
	lowerText := strings.ToLower(text)

	for _, phrase := range phrases {
		if strings.Contains(lowerText, phrase) {
			return true
		}
	}

	return travlyaRegex.MatchString(lowerText) ||
		strings.Contains(lowerText, botUsername) ||
		rand.Float32() < replyChance
}
