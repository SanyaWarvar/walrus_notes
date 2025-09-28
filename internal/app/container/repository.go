package container

import (
	"wn/internal/infrastructure/repository/file"
	tokensRepo "wn/internal/infrastructure/repository/tokens"
	userRepo "wn/internal/infrastructure/repository/user"
)

func (c *Container) getRepositories() *repositories {
	if c.repositories == nil {
		c.repositories = &repositories{c: c}
	}
	return c.repositories
}

type repositories struct {
	c *Container

	user  *userRepo.Repository
	token *tokensRepo.Repository
	file  *file.Repository
}

func (r *repositories) getUserRepository() *userRepo.Repository {
	if r.user == nil {
		r.user = userRepo.NewRepository(r.c.getDBPool())
	}
	return r.user
}

func (r *repositories) getTokenRepository() *tokensRepo.Repository {
	if r.token == nil {
		r.token = tokensRepo.NewRepository(r.c.getDBPool())
	}
	return r.token
}

func (r *repositories) getFileRepository() *file.Repository {
	if r.file == nil {
		r.file = file.NewRepository(r.c.getDBPool())
	}
	return r.file
}
