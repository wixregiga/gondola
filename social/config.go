package social

import (
	"gondola/social/facebook"
	"gondola/social/pinterest"
	"gondola/social/twitter"
)

type FacebookConfig struct {
	App         *facebook.App
	AccessToken string
}

type TwitterConfig struct {
	App   *twitter.App
	Token *twitter.Token
}

type PinterestConfig struct {
	Account *pinterest.Account
	// Might be *pinterest.Board, string (board name) or nil (first board found)
	Board interface{}
}
