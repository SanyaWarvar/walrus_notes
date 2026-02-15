package container

import (
	v1 "wn/internal/endpoint/controller/http/api/v1"
	"wn/internal/endpoint/controller/http/api/v1/auth"
	"wn/internal/endpoint/controller/http/api/v1/file"
	"wn/internal/endpoint/controller/http/api/v1/layout"
	"wn/internal/endpoint/controller/http/api/v1/note"
	"wn/internal/endpoint/controller/http/api/v1/socket"
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

			layout.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getLayoutApplicationService(),
			),

			socket.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getServices().getSocketService(),
				c.getServices().getMultyplayerService(),
			),

			file.NewController(
				c.getLogger(),
				c.getResponseBuilder(),
				c.getApplication().getFileApplciationService(),
			),
		)
	}
	return c.httpDispatcher
}
