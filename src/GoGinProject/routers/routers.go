package routers

import (
	apiControllerV1 "GoGinProject/controllers/api/v1"
	apiControllerV2 "GoGinProject/controllers/api/v2"
	"GoGinProject/middlewares"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//SetupRouter function will perform all route operations
func SetupRouter() *gin.Engine {

	r := gin.Default()

	//Giving access to storage folder
	r.Static("/storage", "storage")

	//Giving access to template folder
	r.Static("/templates", "templates")
	r.LoadHTMLGlob("templates/*")

	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	//API route for version 1
	v1 := r.Group("/api/v1")

	//If you want to pass your route through specific middlewares
	v1.Use(middlewares.UserMiddlewares())
	{
		v1.GET("test-user-list", apiControllerV1.TestUserList)
		v1.POST("user-list", apiControllerV1.UserList)
		r.GET("/user/:name/*action", func(c *gin.Context) {
			name := c.Param("name")
			action := c.Param("action")
			message := name + " is " + action
			c.String(http.StatusOK, message)
		})
		r.GET("/log", func(c *gin.Context) {
			c.File("gin.log")
		})
		r.GET("/long_async", func(c *gin.Context) {
			// create copy to be used inside the goroutine
			cCp := c.Copy()
			go func() {
				// simulate a long task with time.Sleep(). 5 seconds
				time.Sleep(5 * time.Second)

				// note that you are using the copied context "cCp", IMPORTANT
				log.Println("Done! in path " + cCp.Request.URL.Path)
			}()
		})

		r.GET("/long_sync", func(c *gin.Context) {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// since we are NOT using a goroutine, we do not have to copy the context
			log.Println("Done! in path " + c.Request.URL.Path)
		})
	}

	//API route for version 2
	v2 := r.Group("/api/v2")

	v2.POST("user-list", apiControllerV2.UserList)

	return r

}
