package util

type ContextKey string

func (c ContextKey) String() string {
	return "freakbot_" + string(c)
}

var ChatIDContextKey ContextKey = "chat_id"
var UserIDContextKey ContextKey = "user_id"
var UsernameContextKey ContextKey = "username"
