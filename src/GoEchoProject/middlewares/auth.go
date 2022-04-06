package middlewares

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

/*
UserMiddlewares function to add auth
*/
func UserMiddlewares() echo.HandlerFunc {
	return nil
	/*return func(c echo.Context) error {
		cc := &CustomContext{c}
		return next(cc)
	}*/
}
