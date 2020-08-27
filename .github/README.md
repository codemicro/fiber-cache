# fiber-cache
![Run tests](https://github.com/codemicro/fiber-cache/workflows/Run%20tests/badge.svg) ![Gosec audit](https://github.com/codemicro/fiber-cache/workflows/Gosec%20audit/badge.svg) [![Godoc](https://godoc.org/github.com/codemicro/fiber-cache?status.svg)](https://pkg.go.dev/github.com/codemicro/fiber-cache@v1.0.0)

fiber-cache is middleware that provides caching for specific routes in a Fiber application.

### Examples

The most basic caching you can do is like this:

```go
app.Get("/your/route", fcache.New(), yourHandler)
```

This will then cache that endpoint for the default TTL (Time-To-Live) specified. If you have not changed this value, it is 5 minutes.

You can change the default TTL for all endpoints by doing the following:

```go
fcache.Config.DefaultTTL = time.Minute * 2
```

If you want to override the default TTL for a specific endpoint, you can do the following:

```go
app.Get("/your/route", fcache.NewWithTTL(fcache.AutoGenerateKey, time.Second*20), yourHandler)
```

If you want to set your own custom key, you can do the following:

```go
app.Get("/your/route", fcache.NewWithKey("yourKeyHere"), yourHandler)
```

The cache key for the current route is stored in `c.Locals` as `cacheKey`. You can also access the underlying cache engine through the `fcache.Cache` variable. [Please click here for information about that.](https://github.com/patrickmn/go-cache)

```go
app.Get("/your/route", fcache.New(), func(c *fiber.Ctx) {
    cacheKey := c.Locals("cacheKey").(string) // -> cacheKey-0
    
    data, found := fcache.Cache.Get(cacheKey)

    // ...

    c.Send("My cache key is: " + cacheKey) // -> My cache key is: cacheKey-0
})
```

### Reference

[http://godoc.org/github.com/patrickmn/go-cache](http://godoc.org/github.com/patrickmn/go-cache)

### Installation
You must have Go 1.11 or higher installed before attempting installation.

Installation is done using the `go get` command:

```
go get -u github.com/codemicro/fiber-cache
```

You can then import the package as follows:

```go
import (
    fcache "github.com/codemicro/fiber-cache"
)
```

### Benchmarks

Consider the following example:

```go
app.Get("/longtime", fcache.New(), func(c *fiber.Ctx) {
    time.Sleep(time.Second * 10)
    c.Send("Hello world")
})
```

When this endpoint is first requested, it takes 10.035 seconds for a response to be returned. After that, across the next 20 requests, it takes an average of 0.01094 seconds for a response to be returned.

### Licence
fiber-cache is free and open source software covered by the Mozilla Public Licence v2.

#### Third party library licences
* [gofiber/fiber](https://github.com/gofiber/fiber/blob/master/LICENSE)
* [patrickmn/go-cache](https://github.com/patrickmn/go-cache/blob/master/LICENSE)

