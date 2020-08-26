package fcache

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber"
	gc "github.com/patrickmn/go-cache"
)

// createNewCacheForTest is needed as the cache is global to the import in the package. If this is run before
// each new test, it prevents data from other tests leaking over.
func createNewCacheForTest() {
	Cache = gc.New(Config.DefaultTTL, Config.CleanupInterval)
}

func Test_valueStored(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	handlerKey := "sampleKey"

	app.Get("/", New(handlerKey), func(c *fiber.Ctx) {
		c.Send(responseText)
	})

	app.Test(httptest.NewRequest("GET", "/", nil))

	value, _ := Cache.Get(handlerKey)

	if string(value.([]byte)) != responseText {
		t.Fatal("Value in cache does not match expected")
	}
}

func Test_cacheValueReturned(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	modResponse := "thisIsAModifiedResponse"
	handlerKey := "sampleKey"

	app.Get("/", New(handlerKey), func(c *fiber.Ctx) {
		c.Send(responseText)
	})

	Cache.Set(handlerKey, []byte(modResponse), Config.DefaultTTL)

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	returnedResponse := buf.String()

	if returnedResponse != modResponse {
		t.Fatal("Value in cache does not match expected")
	}
}

func Test_customTTL(t *testing.T) {
	createNewCacheForTest()
	app := fiber.New()

	responseText := "thisIsAResponse"
	handlerKey := "sampleKey"

	app.Get("/", NewWithTTL(handlerKey, time.Second*2), func(c *fiber.Ctx) {
		c.Send(responseText)
	})

	app.Test(httptest.NewRequest("GET", "/", nil))

	time.Sleep(time.Second * 3) // by this point, the cache should have expired

	responseText = "thisIsDifferent"

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	returnedResponse := buf.String()

	if returnedResponse != responseText {
		t.Fatal("Cache was not refreshed after TTL expired")
	}
}
