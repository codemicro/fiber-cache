// package fcache provides caching middleware for the Fiber web framework. The caching engine can be accessed through the Cache variable.
package fcache

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	gc "github.com/patrickmn/go-cache"
)

var (
	Cache           *gc.Cache
	currentKeyIndex = 0
	statusCodes     = make(map[string]int)
	codesMutex      sync.Mutex
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

type internalConfig struct {
	CleanupInterval time.Duration
	DefaultTTL      time.Duration
}

func saveStatusCode(key string, code int) {
	codesMutex.Lock()
	statusCodes[key] = code
	codesMutex.Unlock()
}

func getStatusCode(key string) int {
	codesMutex.Lock()
	defer codesMutex.Unlock()
	return statusCodes[key]
}

func createMiddleware(key string, ttl time.Duration) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		val, found := Cache.Get(key)
		if found {
			c.Response().SetBody(val.([]byte))
			c.Response().SetStatusCode(getStatusCode(key))
			return nil
		}

		c.Locals("cacheKey", key)

		c.Next()

		Cache.Set(key, c.Response().Body(), ttl)

		saveStatusCode(key, c.Response().StatusCode())

		return nil

	}
}

func generateKey() string {
	key := "cacheKey-" + strconv.Itoa(currentKeyIndex)
	currentKeyIndex += 1
	return key
}

// New returns a new instance of the caching middleware, with an automatically generated key and the default TTL.
func New() func(*fiber.Ctx) error {
	return createMiddleware(generateKey(), Config.DefaultTTL)
}

// NewWithKey returns a new instance of the caching middleware with the default TTL and the option to set your own cache key. If this is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithKey(key string) func(*fiber.Ctx) error {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, Config.DefaultTTL)
}

// NewWithTTL returns a new instance of the caching middleware with the option to define your own cache key and your own TTL. If the cache key you set is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithTTL(key string, ttl time.Duration) func(*fiber.Ctx) error {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, ttl)
}
