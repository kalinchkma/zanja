package server

import (
	"geko/internal/cache"
	"geko/internal/ratelimiter"
)

type Config struct {
	Addr           string
	Env            string
	MailerCfg      MailerConfig
	AuthCfg        AuthConfig
	RedisCfg       cache.RedisConfig
	RateLimiterCfg ratelimiter.RateLimiterConfig
}

type AuthConfig struct {
}

type MailerConfig struct {
}
