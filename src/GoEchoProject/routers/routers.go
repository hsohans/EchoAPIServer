package routers

import (
	"GoEchoProject/connections"
	apiControllerV1 "GoEchoProject/controllers/api/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

//SetupRouter function will perform all route operations
func SetupRouter(c connections.Connections) *echo.Echo {
	e := echo.New()

	// Logger 설정 (HTTP requests)
	e.Use(middleware.Logger())
	// Recover 설정 (recovers panics, prints stack trace)
	e.Use(middleware.Recover())

	// CORS 설정 (control domain access)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		MaxAge:       86400,
		//AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// swagger 2.0 설정
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Controller 설정
	apiAuthentication := apiControllerV1.GetAuthenticationController(c)
	apiUser := apiControllerV1.GetUserController(c)

	// Router 설정
	v1 := e.Group("/api/v1")
	v1.POST("/token", apiAuthentication.CreateToken)
	v1.GET("/users", apiUser.GetUsers)

	// 아이디 및 비밀번호 확인(BasicAuth)시 JWT 토큰 발급 및 Redis 저장

	// userMiddleware 설정 (입력된 JWT 토큰 검증 및 검증된 요청자 API 접근 허용)

	return e
}
