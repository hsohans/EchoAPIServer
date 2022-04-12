package connections

import (
	"GoEchoProject/helpers"
	"GoEchoProject/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

type Connections struct {
	DbInfo *gorm.DB
}

func SetupConnection() (conn Connections, err error) {

	// DB 설정
	dbConnInfo := models.DBConfig{}
	dbConnInfo.DbType = os.Getenv("db_type")
	dbConnInfo.UserName = os.Getenv("db_user")
	dbConnInfo.UserPassword = os.Getenv("db_pass")
	dbConnInfo.DbName = os.Getenv("db_name")
	dbConnInfo.Host = os.Getenv("db_host")
	dbConnInfo.Port = os.Getenv("db_port")

	paasConnectionString := helpers.GetConnectionString(dbConnInfo.Host, dbConnInfo.Port, dbConnInfo.UserName, dbConnInfo.UserPassword, dbConnInfo.DbName)
	fmt.Println(paasConnectionString)
	//logger.Infof("DB Connection Info : %v\n", paasConnectionString)

	paasDbAccessObj, paasDbErr := gorm.Open(dbConnInfo.DbType, paasConnectionString+"?charset=utf8&parseTime=true")
	if paasDbErr != nil {
		fmt.Println("%v\n", paasDbErr)
		return conn, paasDbErr
	}
	//fmt.Println(paasDbAccessObj)

	c1 := Connections{
		DbInfo: paasDbAccessObj,
	}
	return c1, paasDbErr
}
