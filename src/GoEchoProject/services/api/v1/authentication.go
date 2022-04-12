package v1service

import (
	dao "GoEchoProject/dao/api/v1"
	"GoEchoProject/models"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"os"
	"time"
)

//Gorm Object Struct
type AuthenticationService struct {
	txn *gorm.DB
}

func GetAuthenticationService(txn *gorm.DB) *AuthenticationService {
	return &AuthenticationService{
		txn: txn,
	}
}

func (h *AuthenticationService) CreateToken(apiRequest models.UserInfo, c echo.Context) (string, error) {

	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	results, err := dao.GetUserDao(h.txn).GetUser(apiRequest, c)
	if err != nil {
		fmt.Println("DB Error =========+>", err)
		return "", err
	}
	if len(results) == 0 {
		fmt.Println("DB Error =========+>", "Not Found User")
		return "Not Found User", err
	}

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = apiRequest.Username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
