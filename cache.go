package fcache

import (
	"time"

	"github.com/gofiber/fiber"
	gc "github.com/patrickmn/go-cache"
)

var (
	Config      internalConfig
	statusCodes = make(map[string]int)
	Cache       *gc.Cache
)

func init() {
	Config.DefaultTTL = time.Duration(5 * time.Minute)
	Config.CleanupInterval = time.Duration(10 * time.Minute)
	Cache = gc.New(Config.DefaultTTL, Config.CleanupInterval)
}

type internalConfig struct {
	CleanupInterval time.Duration
	DefaultTTL      time.Duration
}

func createMiddleware(key string, ttl time.Duration) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		val, found := Cache.Get(key)
		if found {
			statusCode := statusCodes[key]
			c.Fasthttp.Response.SetBody(val.([]byte))
			c.Fasthttp.Response.SetStatusCode(statusCode)
			return
		}

		c.Next()

		Cache.Set(key, c.Fasthttp.Response.Body(), ttl)
		statusCodes[key] = c.Fasthttp.Response.StatusCode()

	}
}

func New(key string) func(*fiber.Ctx) {
	return createMiddleware(key, Config.DefaultTTL)
}

func NewWithTTL(key string, ttl time.Duration) func(*fiber.Ctx) {
	return createMiddleware(key, ttl)
}
