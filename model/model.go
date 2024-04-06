package model

import (
	"context"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

type Url struct{
	gorm.Model
	LongUrl string
	ShortCode string
	Visited int64 
	CreatedUser string
}

type MyClaims struct {
    Username string 	`json:"username"`
    jwt.StandardClaims
}

var  DB *gorm.DB
var RedisClient *redis.Client



func ConnDB()  {
	var err error
	dsn := "root:147258aa@tcp(172.17.0.3:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"

	DB,err=gorm.Open("mysql",dsn)

	if err!=nil{
		log.Fatal(err)
	}

	DB.AutoMigrate(&Url{},&User{})

	//连接redis
	RedisClient=redis.NewClient(&redis.Options{
		Addr: "172.17.0.2:6379",
		Password: "",
		DB: 0,
	})

	_,err=RedisClient.Ping(context.Background()).Result()
	if err!=nil{
		log.Fatal(err)
	}
}