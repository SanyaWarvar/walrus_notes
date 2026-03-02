package container

import (
	"wn/internal/infrastructure/cache/permissions"
	smtpCache "wn/internal/infrastructure/cache/smtp"
)

func (c *Container) getCaches() *cache {
	if c.caches == nil {
		c.caches = &cache{c: c}
	}
	return c.caches
}

type cache struct {
	c *Container

	smtp        *smtpCache.Cache
	permissions *permissions.Cache
}

func (s *cache) getSmtpCache() *smtpCache.Cache {
	if s.smtp == nil {
		s.smtp = smtpCache.NewCache(
			s.c.getLogger(),
			s.c.getCacheClient(),
		)
	}
	return s.smtp
}

func (s *cache) getPermissionsCache() *permissions.Cache {
	if s.permissions == nil {
		s.permissions = permissions.NewCache(
			s.c.getLogger(),
			s.c.getCacheClient(),
		)
	}
	return s.permissions
}
