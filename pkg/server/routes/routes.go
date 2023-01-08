package routes

import (
	"github.com/avinal/snowgauge/pkg/server/controllers"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	controllers.Init()
	c := controllers.NewController()

	e.GET("/", c.HomeController)
	e.GET("/stream", c.StreamContoller)
}
