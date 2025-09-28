package container

import smtpCache "wn/internal/infrastructure/cache/smtp"

func (c *Container) getCaches() *cache {
	if c.caches == nil {
		c.caches = &cache{c: c}
	}
	return c.caches
}

type cache struct {
	c *Container

	smtp *smtpCache.Cache
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
