package container

import (
	"wn/internal/infrastructure/repository/file"
	"wn/internal/infrastructure/repository/layout"
	"wn/internal/infrastructure/repository/links"
	"wn/internal/infrastructure/repository/note"
	"wn/internal/infrastructure/repository/positions"
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

	user      *userRepo.Repository
	token     *tokensRepo.Repository
	file      *file.Repository
	note      *note.Repository
	layout    *layout.Repository
	links     *links.Repository
	positions *positions.Repository
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

func (r *repositories) getNoteRepository() *note.Repository {
	if r.note == nil {
		r.note = note.NewRepository(r.c.getDBPool())
	}
	return r.note
}

func (r *repositories) getLayoutRepository() *layout.Repository {
	if r.layout == nil {
		r.layout = layout.NewRepository(r.c.getDBPool())
	}
	return r.layout
}

func (r *repositories) getLinksRepository() *links.Repository {
	if r.links == nil {
		r.links = links.NewRepository(r.c.getDBPool())
	}
	return r.links
}

func (r *repositories) getPositionsRepository() *positions.Repository {
	if r.positions == nil {
		r.positions = positions.NewRepository(r.c.getDBPool())
	}
	return r.positions
}
