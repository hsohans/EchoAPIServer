package routers

import (
	"GoEchoProject/handler"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"html/template"
	"io"
)

// Define the template registry struct
type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

//SetupRouter function will perform all route operations
func SetupRouter() *echo.Echo {
	e := echo.New()

	// Logger 설정
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

	templates := make(map[string]*template.Template)
	templates["home.html"] = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
	templates["about.html"] = template.Must(template.ParseFiles("templates/about.html", "templates/base.html"))
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.Static("/static", "view/static")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// swagger setting
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	/*e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})*/
	/*e.GET("/json", func(c echo.Context) error {
		return c.JSONBlob(
			http.StatusOK,
			[]byte(`{ "id": "1", "msg": "Hello, Boatswain!" }`),
		)
	})
	e.GET("/html", func(c echo.Context) error {
		return c.HTML(
			http.StatusOK,
			"<h1>Hello, Boatswain!</h1>",
		)
	})*/
	// Route => handler
	e.GET("/", handler.HomeHandler)
	e.GET("/about", handler.AboutHandler)

	return e
}
