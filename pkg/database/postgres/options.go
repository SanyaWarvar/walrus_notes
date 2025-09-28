package postgres

import "time"

// Option -.
type Option func(*Pool)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Pool) {
		c.maxPoolSize = size
	}
}

// MaxPoolSize -.
func MinPoolSize(size int) Option {
	return func(c *Pool) {
		c.minPoolSize = size
	}
}

// ConnMaxLifetime -.
func ConnMaxLifetime(maxLifetime time.Duration) Option {
	return func(c *Pool) {
		c.connMaxLifetime = maxLifetime
	}
}

// ConnMaxIdletime -.
func ConnMaxIdletime(maxIdletime time.Duration) Option {
	return func(c *Pool) {
		c.connMaxIdletime = maxIdletime
	}
}

// HealthCheckPeriod -.
func HealthCheckPeriod(healthCheckPeriod time.Duration) Option {
	return func(c *Pool) {
		c.healthCheckPeriod = healthCheckPeriod
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Pool) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Pool) {
		c.connTimeout = timeout
	}
}
