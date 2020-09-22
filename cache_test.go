package fcache

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	gc "github.com/patrickmn/go-cache"
)

// createNewCacheForTest is needed as the cache is global to the import in the package. If this is run before
// each new test, it prevents data from other tests leaking over.
func createNewCacheForTest() {
	Cache = gc.New(Config.DefaultTTL, Config.CleanupInterval)
}

func getResponseBody(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}

func Test_valueStored(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	handlerKey := "sampleKey"

	app.Get("/", NewWithKey(handlerKey), func(c *fiber.Ctx) error {
		return c.SendString(responseText)
	})

	app.Test(httptest.NewRequest("GET", "/", nil))

	value, _ := Cache.Get(handlerKey)

	if string(value.(CacheEntry).Body) != responseText {
		t.Fatal("Value in cache does not match expected")
	}
}

func Test_cacheValueReturned(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	modResponse := "thisIsAModifiedResponse"
	handlerKey := "sampleKey"

	app.Get("/", NewWithKey(handlerKey), func(c *fiber.Ctx) error {
		c.SendString(responseText)
		return nil
	})

	Cache.Set(handlerKey, CacheEntry{
		Body:        []byte(modResponse),
		StatusCode:  200,
		ContentType: []byte("text/plain"),
	}, Config.DefaultTTL)

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	returnedResponse := getResponseBody(resp)

	if returnedResponse != modResponse {
		t.Fatal("Value in cache does not match expected")
	}
}

func Test_customTTL(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	handlerKey := "sampleKey"

	app.Get("/", NewWithTTL(handlerKey, time.Second*2), func(c *fiber.Ctx) error {
		c.SendString(responseText)
		return nil
	})

	app.Test(httptest.NewRequest("GET", "/", nil))

	time.Sleep(time.Second * 3) // by this point, the cache should have expired

	responseText = "thisIsDifferent"

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	returnedResponse := getResponseBody(resp)

	if returnedResponse != responseText {
		t.Fatal("Cache was not refreshed after TTL expired")
	}
}

func Test_automaticKeyGeneration(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText1 := "thisIsAResponse"
	responseText2 := "thisIsADifferentResponse"

	app.Get("/", New(), func(c *fiber.Ctx) error {
		c.SendString(responseText1)
		return nil
	})

	app.Get("/other", New(), func(c *fiber.Ctx) error {
		c.SendString(responseText2)
		return nil
	})

	resp1, _ := app.Test(httptest.NewRequest("GET", "/", nil))
	resp2, _ := app.Test(httptest.NewRequest("GET", "/other", nil))

	returnedResponse1 := getResponseBody(resp1)
	returnedResponse2 := getResponseBody(resp2)

	if returnedResponse1 == returnedResponse2 {
		t.Fatal("A collision has occured between automatically generated keys")
	} else if returnedResponse1 != responseText1 || returnedResponse2 != responseText2 {
		t.Fatal("Incorrect values are being stored in cache for automatically generated keys")
	}
}

func Test_automaticKeyGenerationWithFlag(t *testing.T) {
	createNewCacheForTest()

	app := fiber.New()
	responseText := "hello world this is a response"
	var cacheKey string

	app.Get("/", NewWithKey(AutoGenerateKey), func(c *fiber.Ctx) error {
		cacheKey = c.Locals("cacheKey").(string)
		c.SendString(responseText)
		return nil
	})

	app.Test(httptest.NewRequest("GET", "/", nil))

	if cacheKey == "" {
		t.Fatal("Automatic key generation is failing or the cache key is not being set")
	}
}

func Test_correctContentTypeSet(t *testing.T) {
	createNewCacheForTest()

	app := fiber.New()

	app.Get("/", New(), func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"hello":"world"})
	})

	app.Test(httptest.NewRequest("GET", "/", nil)) // load cache response

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	if resp.Header["Content-Type"][0] != "application/json" {
		t.Fatal("The content type header is not being set correctly")
	}
}