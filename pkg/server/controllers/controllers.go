package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/avinal/snowgauge/pkg/server/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
	"log"
)

type Controller interface {
	HomeController(e echo.Context) error
	StreamContoller(e echo.Context) error
}

type controller struct {
}

var model models.Model

func (c *controller) HomeController(e echo.Context) error {
	return e.File("views/index.html")
}

func (c controller) StreamContoller(e echo.Context) error {
	websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		status, err := model.GetLiveNetworkStat()
		if err != nil {
			log.Printf("error: %v", err)
		}
		for {
			newVal := <-status
			jsonResponse, _ := json.Marshal(newVal)
			err := websocket.Message.Send(conn, fmt.Sprintln(string(jsonResponse)))
			if err != nil {
				log.Printf("error: %v", err)
			}
		}
	}).ServeHTTP(e.Response(), e.Request())
	return nil
}

func NewController() Controller {
	return &controller{}
}

func Init() {
	model = models.NewModel()
}
