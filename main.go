package main

import (
	"github.com/labstack/echo"
	"github.com/sevenNt/echo-pprof"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"strconv"
	"time"
)

type DummyAllocation struct {
	I int
}

func main() {
	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	backgroundGoroutine := func(createAlloc bool) {
		localAllocations := make([]*DummyAllocation, 0)
		for i := 0; ; i++ {
			time.Sleep(1 * time.Second)
			if createAlloc {
				alloc := new(DummyAllocation)
				alloc.I = i
				localAllocations = append(localAllocations, alloc)
				log.Printf("new allocation ")
			}
		}
	}

	e.GET("/background/goroutines/:count", func(c echo.Context) error {
		countRaw := c.Param("count")
		count, err := strconv.Atoi(countRaw)
		if err != nil {
			return err
		}
		for i := 0; i < count; i++ {
			go backgroundGoroutine(false)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status": "ok",
		})
	})

	e.GET("/allocations/:count", func(c echo.Context) error {
		countRaw := c.Param("count")
		count, err := strconv.Atoi(countRaw)
		if err != nil {
			return err
		}
		for i := 0; i < count; i++ {
			go backgroundGoroutine(true)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status": "ok",
		})
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	echopprof.Wrap(e)

	// echopprof also plays well with *echo.Group
	// prefix := "/debug/pprof"
	// group := e.Group(prefix)
	// echopprof.WrapGroup(prefix, group)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlIndex)
	})

	e.GET("/debug/pprof/allocs", func(c echo.Context) error {
		pprof.Handler("allocs").ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.Start(":8080")
}

const htmlIndex = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>HTML 5 Boilerplate</title>
    <link rel="stylesheet" href="style.css">
  </head>
  <body>
	CONTROL: <br/>
	<a href="/allocations/5">create some allocations in background</a><br/>
	<a href="/background/goroutines/5">create some background goroutines</a><br/>
	<br/>
	DEBUG: <br/>
	<a href="/debug/pprof/">SEE PPROF</a>
  </body>
</html>
`
