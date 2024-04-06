package user

import (
	"context"
	"log"
	"net/http"
	"shorturl/kitex_gen/short/user"
	"shorturl/kitex_gen/short/user/userservice"
	"shorturl/middleware"
	"shorturl/model"
	"strings"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

//登录实现
func Login(c *gin.Context)  {
    req:=user.NewUserRequest()
	

	req.Username=c.PostForm("username")
	req.Password=c.PostForm("password")

	if strings.ContainsAny(req.Username, " \t\n")||strings.ContainsAny(req.Password," \t\n"){
		c.JSON(http.StatusOK,gin.H{
			"message":"用户名or密码格式错误",
		})
		return
	}

	cli,err:=userservice.NewClient("short.user",client.WithHostPorts("0.0.0.0:8888"))
	if err!=nil{
		log.Fatal(err)
	}
	resp,_:=cli.Login(context.Background(),req)
	
	//token
	if resp.Message=="登录成功"{
		claim := model.MyClaims{
			Username: req.Username, 	
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),		
				ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), 
				Issuer:    "my-shorturl",              //签发人
			},
		}
		//Header&Payload
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

		//+Signature
		tokenString, err := token.SignedString(middleware.Secret)
		if err!=nil{
			log.Fatal(err)
		}  


		c.JSON(http.StatusOK,gin.H{
			"token":tokenString,
		})

	}
	c.JSON(http.StatusOK,gin.H{
		"message":resp.Message,
	})
}


//注册实现
func Register(c *gin.Context)  {
    req:=user.NewUserRequest()
	
	req.Username=c.PostForm("username")
	req.Password=c.PostForm("password")

	if strings.ContainsAny(req.Username, " \t\n")||strings.ContainsAny(req.Password," \t\n"){
		c.JSON(http.StatusOK,gin.H{
			"message":"用户名or密码包含空白字符",
		})
		return
	}

	if len(req.Password)<6{
		c.JSON(http.StatusOK,gin.H{
			"message":"密码至少包含6位",
		})
		return
	}


	cli,err:=userservice.NewClient("short.user",	client.WithHostPorts("0.0.0.0:8888"))
	if err!=nil{
		log.Fatal(err)
	}


	resp,_:=cli.Register(context.Background(),req)
	c.JSON(http.StatusOK,gin.H{
		"message":resp.Message,
	})
	
}