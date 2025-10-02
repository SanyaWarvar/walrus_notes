package container

import (
	v1 "wn/internal/endpoint/controller/http/api/v1"
	"wn/internal/endpoint/controller/http/api/v1/auth"
	"wn/internal/endpoint/controller/http/api/v1/note"
	"wn/internal/endpoint/controller/http/api/v1/user"
)

func (c *Container) getHTTPDispatcher() *v1.Dispatcher {
	if c.httpDispatcher == nil {
		c.httpDispatcher = v1.NewDispatcher(
			c.getConfig().Internal.Path,

			auth.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getUserApplicationService(),
				c.getApplication().getAuthApplicationService(),
			),

			user.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getUserApplicationService(),
			),

			note.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getNoteApplicationService(),
			),
		)
	}
	return c.httpDispatcher
}
