package fcache

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber"
	gc "github.com/patrickmn/go-cache"
)

var (
	Cache           *gc.Cache
	Config          internalConfig
	statusCodes     = make(map[string]int)
	currentKeyIndex = 0
	AutoGenerateKey = ""
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

func generateKey() string {
	key := "cacheKey-" + strconv.Itoa(currentKeyIndex)
	currentKeyIndex += 1
	return key
}

func New() func(*fiber.Ctx) {
	return createMiddleware(generateKey(), Config.DefaultTTL)
}

func NewWithKey(key string) func(*fiber.Ctx) {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, Config.DefaultTTL)
}

func NewWithTTL(key string, ttl time.Duration) func(*fiber.Ctx) {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, ttl)
}
