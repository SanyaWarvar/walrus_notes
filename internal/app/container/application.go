package container

import (
	"wn/internal/application/auth"
	userApp "wn/internal/application/user"
)

func (c *Container) getApplication() *applications {
	if c.applications == nil {
		c.applications = &applications{c: c}
	}
	return c.applications
}

type applications struct {
	c *Container

	user *userApp.Service
	auth *auth.Service
}

func (s *applications) getUserApplicationService() *userApp.Service {
	if s.user == nil {
		s.user = userApp.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getUserService(),
			s.c.getServices().getFileService(),
		)
	}
	return s.user
}

func (s *applications) getAuthApplicationService() *auth.Service {
	if s.auth == nil {
		s.auth = auth.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getUserService(),
			s.c.getServices().getSMTPService(),
			s.c.getServices().getTokenService(),
		)
	}
	return s.auth
}
