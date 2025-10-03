package v1

import (
	"wn/internal/endpoint/controller/http/api/v1/auth"
	"wn/internal/endpoint/controller/http/api/v1/layout"
	"wn/internal/endpoint/controller/http/api/v1/note"
	"wn/internal/endpoint/controller/http/api/v1/user"

	"github.com/gin-gonic/gin"
)

type Dispatcher struct {
	apiPath string

	auth   *auth.Controller
	user   *user.Controller
	note   *note.Controller
	layout *layout.Controller
}

func NewDispatcher(
	apiPath string,

	auth *auth.Controller,
	user *user.Controller,
	note *note.Controller,
	layout *layout.Controller,
) *Dispatcher {
	return &Dispatcher{
		apiPath: apiPath,
		auth:    auth,
		user:    user,
		note:    note,
		layout:  layout,
	}
}

func (d *Dispatcher) Init(router *gin.RouterGroup, authorization gin.HandlerFunc) {
	api := router.Group("/v1")
	{
		d.auth.Init(api)
		authorizedGroup := api.Group("", authorization)
		{
			d.user.Init(api, authorizedGroup)
			d.note.Init(api, authorizedGroup)
			d.layout.Init(api, authorizedGroup)
		}
	}
}
