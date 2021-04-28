// +build !appengine

package app

import "gondola/log"

func (c *Context) logger() log.Interface {
	if c == nil || c.app.Logger == nil {
		return nullLogger{}
	}
	return c.app.Logger
}
