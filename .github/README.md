# fiber-cache

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
fcache.NewWithTTL(fcache.AutoGenerateKey, time.Second*20)
```

If you want to set your own custom key, you can do the following:

```go
fcache.NewWithKey("yourKeyHere")
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

