package container

import (
	"wn/internal/application/auth"
	"wn/internal/application/file"
	"wn/internal/application/layout"
	"wn/internal/application/note"
	"wn/internal/application/permissions"
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

	user        *userApp.Service
	auth        *auth.Service
	note        *note.Service
	layout      *layout.Service
	file        *file.Service
	permissions *permissions.Application
}

func (s *applications) getUserApplicationService() *userApp.Service {
	if s.user == nil {
		s.user = userApp.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getUserService(),
			s.c.getServices().getFileService(),
			s.c.getServices().getLayoutService(),
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
			s.c.getServices().getLayoutService(),
		)
	}
	return s.auth
}

func (s *applications) getNoteApplicationService() *note.Service {
	if s.note == nil {
		s.note = note.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getNoteService(),
			s.c.getServices().getPermissionsService(),
			s.c.getRepositories().getLayoutRepository(),
		)
	}
	return s.note
}

func (s *applications) getLayoutApplicationService() *layout.Service {
	if s.layout == nil {
		s.layout = layout.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getLayoutService(),
			s.c.getServices().getPermissionsService(),
		)
	}
	return s.layout
}

func (s *applications) getFileApplciationService() *file.Service {
	if s.file == nil {
		s.file = file.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getFileService(),
		)
	}
	return s.file
}

func (s *applications) getPermissionsApplicationService() *permissions.Application {
	if s.permissions == nil {
		s.permissions = permissions.NewApplication(
			s.c.getTransactionManager(),
			s.c.getLogger(),

			s.c.getServices().getPermissionsService(),
			s.c.getRepositories().getPermissionsRepository(),
			s.c.getCaches().getPermissionsCache(),
			s.c.getRepositories().getLayoutRepository(),
			s.c.getRepositories().getNoteRepository(),
		)
	}
	return s.permissions
}
