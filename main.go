package main

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/templates"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = make(map[*websocket.Conn]bool)
	mu       sync.Mutex
)

type SubmitRequest struct {
	Name      string   `json:"name"`
	ArrayData []string `json:"items[]"`
}
type SubmitRequestForm struct {
	Name      string   `form:"name"`
	ArrayData []string `form:"items[]"`
}

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		component := templates.Homepage()
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return component.Render(c.Request().Context(), c.Response().Writer)
	})
	e.GET("/ws/submit", func(c echo.Context) error {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		mu.Lock()
		clients[ws] = true
		mu.Unlock()

		for {
			// Read
			_, msg, err := ws.ReadMessage()

			if err != nil {
				mu.Lock()
				delete(clients, ws)
				mu.Unlock()
				log.Println(err)
				c.Logger().Error(err)
				break
			}

			var data SubmitRequest
			if err := json.Unmarshal(msg, &data); err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("not decoded"))
				fmt.Printf("JSON unmarshal error: %v", err)
				continue
			}

			ws.WriteMessage(websocket.TextMessage, []byte("decoded successfully"))
			fmt.Println(data)
		}
		return nil
	})
	e.POST("/submit", func(c echo.Context) error {
		var data SubmitRequestForm
		if err := c.Bind(&data); err != nil {
			c.String(http.StatusOK, "failed to bind data")
		}
		fmt.Println(data)
		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
