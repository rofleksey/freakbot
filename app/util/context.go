package util

type ContextKey string

func (c ContextKey) String() string {
	return "freakbot_" + string(c)
}

var ChatIDContextKey ContextKey = "chat_id"
var ChatNameContextKey ContextKey = "chat_name"
var UserIDContextKey ContextKey = "user_id"
var UsernameContextKey ContextKey = "username"
