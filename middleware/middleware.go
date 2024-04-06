package middleware

import (
	
	"net/http"
	"shorturl/model"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var UserName string

var Secret = []byte("JwtSecret")



//解析token  
func JWTAuthMiddleware() func(c *gin.Context) {
    return func(c *gin.Context) {
        
        authHeader := c.Request.Header.Get("Authorization")     

        if authHeader == "" {
            c.JSON(http.StatusOK, gin.H{
                "msg":  "请求头中auth为空",
            })
            c.Abort()       //登陆失败
            return
        }
        
        
        //分割请求头
        parts := strings.SplitN(authHeader, " ", 2) 
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusOK, gin.H{
                "msg":  "请求头中auth格式有误",
            })
            c.Abort()
            return
        }
        
       //解析
       token, err := jwt.ParseWithClaims(parts[1], &model.MyClaims{}, func(token *jwt.Token) ( interface{},  error) {
        return Secret, nil
        })

       if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "msg":  "无效的Token",
        })
        c.Abort()
        return 
       }

        //获取当前用户名
       UserName=token.Claims.(*model.MyClaims).Username
       
    }
}



