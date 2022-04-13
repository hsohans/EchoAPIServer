package middlewares

import (
	"GoEchoProject/connections"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type CustomContext struct {
	echo.Context
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

// 1. 토큰 추출
func ExtractToken(c echo.Context) (string, error) {
	var err error
	req := c.Request()
	bearToken := req.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		// Bearer 제외
		return strArr[1], err
	} else if len(bearToken) == 0 || len(strArr) == 0 {
		// Token 정보가 없으면
		return "", fmt.Errorf("Please enter token")
	}
	return bearToken, err
}

// 2. 토큰 검증 (signing method 검증, 서명 검증)
func VerifyToken(bearToken string, c echo.Context) (*jwt.Token, error) {
	token, err := jwt.Parse(bearToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// 3. 토큰 만료 검증
func TokenValid(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("Invaild token")
	}
	return nil
}

// 4. 메타 데이터 추출 (메타데이터 이용한 Redis 확인)
func ExtractTokenMetadata(token *jwt.Token) (map[string]interface{}, error) {
	// metadata 초기화 선언
	metadata := make(map[string]interface{})
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return metadata, fmt.Errorf("No access_uuid in metadata")
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return metadata, fmt.Errorf("No user_id in metadata")
		}
		metadata["access_uuid"] = accessUuid
		metadata["user_id"] = userId
		return metadata, nil
	}
	return metadata, fmt.Errorf("Token is invalid")
}

// 5. Redis 추출 (UUID를 이용한 userId 추출)
func FetchAuth(metadata map[string]interface{}, RedisInfo *redis.Client) (string, error) {
	userid, err := RedisInfo.Get(metadata["access_uuid"].(string)).Result()
	if err != nil {
		return "", fmt.Errorf("User is not exist")
	}
	return userid, nil
}

func CheckToken(conn connections.Connections) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// 1. 토큰 추출
			bearToken, err := ExtractToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			// 2. 토큰 검증 (signing method 검증, 서명 검증)
			token, err := VerifyToken(bearToken, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// 3. 토큰 만료 검증 (동작 방식 확인 필요)
			err = TokenValid(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// 4. 메타 데이터 추출 (메타데이터 이용한 Redis 확인)
			metadata, err := ExtractTokenMetadata(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			fmt.Println(metadata)

			// 5. Redis 추출 (UUID를 이용한 userId 추출)
			userId, err := FetchAuth(metadata, conn.RedisInfo)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			fmt.Println(userId + "is authorized")

			// Continue
			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}
