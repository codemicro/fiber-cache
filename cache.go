// package fcache provides caching middleware for the Fiber web framework. The caching engine can be accessed through the Cache variable.
package fcache

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber"
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

func createMiddleware(key string, ttl time.Duration) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		val, found := Cache.Get(key)
		if found {
			c.Fasthttp.Response.SetBody(val.([]byte))
			c.Fasthttp.Response.SetStatusCode(getStatusCode(key))
			return
		}

		c.Locals("cacheKey", key)

		c.Next()

		Cache.Set(key, c.Fasthttp.Response.Body(), ttl)

		saveStatusCode(key, c.Fasthttp.Response.StatusCode())

	}
}

func generateKey() string {
	key := "cacheKey-" + strconv.Itoa(currentKeyIndex)
	currentKeyIndex += 1
	return key
}

// New returns a new instance of the caching middleware, with an automatically generated key and the default TTL.
func New() func(*fiber.Ctx) {
	return createMiddleware(generateKey(), Config.DefaultTTL)
}

// NewWithKey returns a new instance of the caching middleware with the default TTL and the option to set your own cache key. If this is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithKey(key string) func(*fiber.Ctx) {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, Config.DefaultTTL)
}

// NewWithTTL returns a neew instance of the caching middleware with the option to define your own cache key and your own TTL. If the cache key you set is an empty string or AutoGenerateKey, a key will be automatically generated.
func NewWithTTL(key string, ttl time.Duration) func(*fiber.Ctx) {
	if key == AutoGenerateKey {
		key = generateKey()
	}
	return createMiddleware(key, ttl)
}
