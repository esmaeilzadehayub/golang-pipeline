package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("error in getting hostname", err)
	}
	e.GET("/api/v1", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s - %s", hostname, time.Now().Format(time.RFC3339)))
	})
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
