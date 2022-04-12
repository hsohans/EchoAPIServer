package v1

import (
	"GoEchoProject/connections"
	"GoEchoProject/models"
	v1service "GoEchoProject/services/api/v1"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthenticationController struct {
	DbInfo *gorm.DB
}

func GetAuthenticationController(c connections.Connections) *AuthenticationController {
	return &AuthenticationController{
		DbInfo: c.DbInfo,
	}
}

func (a *AuthenticationController) CreateToken(c echo.Context) (err error) {
	/* Request Body Data Mapping */
	var apiRequest models.UserInfo // -> &추가
	//apiRequest := new(models.UserInfo) // &제거
	if err = c.Bind(&apiRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return nil
	}

	// Authentication의 CreateToken을 호출한다.
	token, err := v1service.GetAuthenticationService(a.DbInfo).CreateToken(apiRequest, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return nil
	}

	// 토근을 발급한다.
	c.JSON(http.StatusOK, token)
	return nil
}
