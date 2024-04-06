package middleware

import (
	
	"net/http"
	"shorturl/model"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var UserName string


//用于签名的字符串
var Secret = []byte("JwtSecret")



//获取解析token  
func JWTAuthMiddleware() func(c *gin.Context) {
    return func(c *gin.Context) {
        
        authHeader := c.Request.Header.Get("Authorization")         //取得请求头中的json
        if authHeader == "" {
            c.JSON(http.StatusOK, gin.H{
                "msg":  "请求头中auth为空",
            })
            c.Abort()       //未登录 阻止调用后续的函数   
            return
        }
        
        
        // (取到token之后) 按空格分割
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



