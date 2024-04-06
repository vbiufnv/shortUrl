package main

import (
	"shorturl/api/url"
	"shorturl/api/user"
	"shorturl/middleware"
	"shorturl/model"

	"github.com/gin-gonic/gin"
)



func main() {

	model.ConnDB()
	defer model.DB.Close()

	r:=gin.Default()
	
	r.POST("/login",user.Login)
	r.POST("/register",user.Register)
	r.GET("/:shortCode", url.RedirectHandler)

	UserRouter := r.Group("/url")  
 
   
    {
        UserRouter.Use(middleware.JWTAuthMiddleware())  //用户验证中间件
		UserRouter.GET("/update", url.Update)    
		UserRouter.GET("/add", url.Add)    
		UserRouter.GET("/search", url.Search)    
		UserRouter.GET("/delete",url.Delete )
		UserRouter.GET("/rank",url.Rank)   

    }

	r.Run(":8080")

   
	
}