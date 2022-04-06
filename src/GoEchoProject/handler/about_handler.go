package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func AboutHandler(c echo.Context) error {
	// Please note the the second parameter "about.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"name": "About",
		"msg":  "All about hsohans!",
	})
}
