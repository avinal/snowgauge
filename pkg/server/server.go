package server

import (
	sf "github.com/avinal/snowgauge/pkg/lib"
	"github.com/avinal/snowgauge/pkg/server/routes"
	"github.com/labstack/echo/v4"
	"sync"
)

func Init(proxy sf.SnowflakeProxy) {
	e := echo.New()
	routes.Init(e)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		e.Logger.Fatal(proxy.Start())
		wg.Done()
	}()
	go func() {
		e.Logger.Fatal(e.Start(":9030"))
		wg.Done()
	}()
	wg.Wait()
}
