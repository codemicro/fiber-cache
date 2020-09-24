// package fcache provides caching middleware for the Fiber web framework. The caching engine can be accessed through the Cache variable.
package fcache

import (
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	gc "github.com/patrickmn/go-cache"
)

var (
	Cache           *gc.Cache
	currentKeyIndex = 0
	Config          internalConfig
)

const (
	AutoGenerateKey = ""
	NoExpiration    = gc.NoExpiration
)

func init() {
	Config.DefaultTTL = 5 * time.Minute
	Config.CleanupInterval = 10 * time.Minute
	Cache = gc.New(Config.DefaultTTL, Config.CleanupInterval)
}

type CacheEntry struct {
	Body []byte
	StatusCode int
	ContentType []byte
}

type internalConfig struct {
	CleanupInterval time.Duration
	DefaultTTL      time.Duration
}

func createMiddleware(key string, ttl time.Duration) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		fmt.Println()
		path := utils.ImmutableString(c.Path() + " ")

		if key == AutoGenerateKey {
			key = path
		}

		fmt.Println("Key", "'" + key + "'")

		val, found := Cache.Get(key)
		if found {
			fmt.Println("Found entry")
			entry := val.(CacheEntry)
			c.Response().SetBody(entry.Body)
			c.Response().SetStatusCode(entry.StatusCode)
			c.Response().Header.SetContentTypeBytes(entry.ContentType)
			return nil
		}
		fmt.Println("Not found entry")

		c.Locals("cacheKey", key)

		err := c.Next()

		if err == nil {
			Cache.Set(key, CacheEntry{
				Body:        c.Response().Body(),
				StatusCode:  c.Response().StatusCode(),
				ContentType: c.Response().Header.ContentType(),
			}, ttl)
		}

		return err

	}
}

// New returns a new instance of the caching middleware, with an automatically generated key and the default TTL.
func New() func(*fiber.Ctx) error {
	return createMiddleware(AutoGenerateKey, Config.DefaultTTL)
}

// NewWithKey returns a new instance of the caching middleware with the default TTL and the option to set your own cache key. If this is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithKey(key string) func(*fiber.Ctx) error {
	return createMiddleware(key, Config.DefaultTTL)
}

// NewWithTTL returns a new instance of the caching middleware with the option to define your own cache key and your own TTL. If the cache key you set is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithTTL(key string, ttl time.Duration) func(*fiber.Ctx) error {
	return createMiddleware(key, ttl)
}
