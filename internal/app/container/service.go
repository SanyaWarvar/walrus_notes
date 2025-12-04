package container

import (
	"wn/internal/domain/services/file"
	"wn/internal/domain/services/layout"
	"wn/internal/domain/services/note"
	smtpSrv "wn/internal/domain/services/smtp"
	"wn/internal/domain/services/socket"
	tokenSrv "wn/internal/domain/services/token"
	userSrv "wn/internal/domain/services/user"
)

func (c *Container) getServices() *services {
	if c.services == nil {
		c.services = &services{c: c}
	}
	return c.services
}

type services struct {
	c *Container

	user          *userSrv.Service
	smtp          *smtpSrv.Service
	token         *tokenSrv.Service
	file          *file.Service
	note          *note.Service
	layout        *layout.Service
	socketManager *socket.Service
}

func (s *services) getUserService() *userSrv.Service {
	if s.user == nil {
		s.user = userSrv.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getRepositories().getUserRepository(),
		)
	}
	return s.user
}

func (s *services) getSMTPService() *smtpSrv.Service {
	if s.smtp == nil {
		s.smtp = smtpSrv.NewService(
			s.c.getLogger(),
			smtpSrv.NewConfig(
				s.c.getConfig().Email.OwnerEmail,
				s.c.getConfig().Email.OwnerPassword,
				s.c.getConfig().Email.Address,
				s.c.getConfig().Email.CodeLenght,
				s.c.getConfig().Email.CodeExp,
				s.c.getConfig().Email.MinTTL,
			),
			s.c.getCaches().getSmtpCache(),
		)
	}
	return s.smtp
}

func (s *services) getTokenService() *tokenSrv.Service {
	if s.token == nil {
		s.token = tokenSrv.NewService(
			s.c.getConfig().Jwt.RefreshTTL,
			s.c.getConfig().Jwt.AccessTTL,
			s.c.getConfig().Jwt.JwtSecret,
			s.c.getRepositories().getTokenRepository(),
		)

	}
	return s.token
}

func (s *services) getFileService() *file.Service {
	if s.file == nil {
		s.file = file.NewService(
			s.c.getLogger(),
			s.c.getRepositories().getFileRepository(),
		)

	}
	return s.file
}

func (s *services) getNoteService() *note.Service {
	if s.note == nil {
		s.note = note.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getRepositories().getNoteRepository(),
			s.c.getRepositories().getLayoutRepository(),
			s.c.getRepositories().getLinksRepository(),
			s.c.getRepositories().getPositionsRepository(),
		)

	}
	return s.note
}

func (s *services) getLayoutService() *layout.Service {
	if s.layout == nil {
		s.layout = layout.NewService(
			s.c.getTransactionManager(),
			s.c.getLogger(),
			s.c.getRepositories().getLayoutRepository(),
			s.c.getRepositories().getLinksRepository(),
			s.c.getRepositories().getNoteRepository(),
			s.c.getRepositories().getPositionsRepository(),
		)

	}
	return s.layout
}

func (s *services) getSocketService() *socket.Service {
	if s.socketManager == nil {
		s.socketManager = socket.NewService(s.c.getLogger())
	}
	return s.socketManager
}
