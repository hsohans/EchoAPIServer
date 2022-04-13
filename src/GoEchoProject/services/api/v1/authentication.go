package v1service

import (
	dao "GoEchoProject/dao/api/v1"
	"GoEchoProject/models"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
	"os"
	"strconv"
	"time"
)

//Gorm Object Struct
type AuthenticationService struct {
	DbInfo    *gorm.DB
	RedisInfo *redis.Client
}

func GetAuthenticationService(DbInfo *gorm.DB, RedisInfo *redis.Client) *AuthenticationService {
	return &AuthenticationService{
		DbInfo:    DbInfo,
		RedisInfo: RedisInfo,
	}
}

func (h *AuthenticationService) CreateToken(apiRequest models.UserInfo, c echo.Context) (models.TokenDetails, error) {
	// 아이디 및 비밀번호 확인 시 JWT 토큰 발급 및 Redis 저장

	// Token 모델을 선언한다.
	td := models.TokenDetails{}

	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다.
	results, err := dao.GetUserDao(h.DbInfo).GetUser(apiRequest, c)
	if err != nil {
		return td, err
	}
	if len(results) == 0 {
		return td, fmt.Errorf("Not Found User")
	}
	// 계정에 대한 비밀번호를 확인한다.
	if results[0].Password != apiRequest.Password {
		return td, fmt.Errorf("The Password is incorrect")
	}

	// Access Token 만료 시간 (현재시간 + 15분)
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()
	// Refresh Token 만료 시간 (현재시간 + 7일)
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	// Access Token을 생성한다.
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = apiRequest.Username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	td.AccessToken = token

	// Refresh Token을 생성한다.
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = apiRequest.Username
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return td, err
	}

	// at는 AccessToken의 접근 유효 시간
	// rt는 RefreshToken의 만료 시간
	redis_at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
	redis_rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	// Todo Redis 생성시 value 처리
	// JWT을 Redis에 저장한다.
	errAccess := h.RedisInfo.Set(td.AccessUuid, strconv.Itoa(int(1)), redis_at.Sub(now)).Err()
	if errAccess != nil {
		return td, errAccess
	}
	errRefresh := h.RedisInfo.Set(td.RefreshUuid, strconv.Itoa(int(1)), redis_rt.Sub(now)).Err()
	if errRefresh != nil {
		return td, errRefresh
	}
	return td, nil
}
